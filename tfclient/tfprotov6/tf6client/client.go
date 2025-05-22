// This is derived from github.com/hashicorp/terraform/internal/plugin6/grpc_provider.go

package tf6client

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-plugin"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/configschema"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/convert"
	"github.com/magodo/terraform-client-go/tfclient/typ"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	ctyjson "github.com/zclconf/go-cty/cty/json"
	"github.com/zclconf/go-cty/cty/msgpack"
)

// Client handles the client, or core side of the plugin rpc connection.
// The Client methods are mostly a translation layer between the
// terraform providers types and the grpc proto types, directly converting
// between the two.
type Client struct {
	// PluginClient provides a reference to the plugin.Client which controls the plugin process.
	// This allows the Client a way to shutdown the plugin process.
	pluginClient *plugin.Client

	// Proto client use to make the grpc service calls.
	client tfprotov6.ProviderServer

	// schema stores the schema for this provider. This is used to properly
	// serialize the state for requests.
	schemas typ.GetProviderSchemaResponse

	configured   bool
	configuredMu sync.Mutex
}

func New(pluginClient *plugin.Client, grpcClient tfprotov6.ProviderServer, schema *typ.GetProviderSchemaResponse) (*Client, error) {
	c := &Client{
		pluginClient: pluginClient,
		client:       grpcClient,
	}

	if schema != nil {
		c.schemas = *schema
		return c, nil
	}

	resp, err := grpcClient.GetProviderSchema(context.Background(), &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		return nil, err
	}
	if diags := convert.DecodeDiagnostics(resp.Diagnostics); diags.HasErrors() {
		return nil, diags.Err()
	}

	schemas := typ.GetProviderSchemaResponse{
		ResourceTypes:    map[string]tfjson.Schema{},
		ResourceTypesCty: map[string]cty.Type{},
		DataSources:      map[string]tfjson.Schema{},
		DataSourcesCty:   map[string]cty.Type{},
		Functions:        map[string]typ.FunctionDecl{},
		ServerCapabilities: typ.ServerCapabilities{
			PlanDestroy: false,
		},
	}
	if resp.Provider != nil {
		providerSchema := convert.ProtoToProviderSchema(resp.Provider)
		schemas.Provider = providerSchema
		schemas.ProviderCty = configschema.SchemaBlockImpliedType(providerSchema.Block)
	}
	if resp.ProviderMeta != nil {
		providerMetaSchema := convert.ProtoToProviderSchema(resp.ProviderMeta)
		schemas.ProviderMeta = providerMetaSchema
		schemas.ProviderMetaCty = configschema.SchemaBlockImpliedType(providerMetaSchema.Block)
	}
	if resp.ServerCapabilities != nil {
		schemas.ServerCapabilities.PlanDestroy = resp.ServerCapabilities.PlanDestroy
	}
	for name, schema := range resp.ResourceSchemas {
		resourceSchema := convert.ProtoToProviderSchema(schema)
		schemas.ResourceTypes[name] = resourceSchema
		schemas.ResourceTypesCty[name] = configschema.SchemaBlockImpliedType(resourceSchema.Block)
	}
	for name, schema := range resp.DataSourceSchemas {
		dataSourceSchema := convert.ProtoToProviderSchema(schema)
		schemas.DataSources[name] = dataSourceSchema
		schemas.DataSourcesCty[name] = configschema.SchemaBlockImpliedType(dataSourceSchema.Block)
	}
	for name, fun := range resp.Functions {
		schemas.Functions[name], err = convert.FunctionDeclFromProto(fun)
		if err != nil {
			return nil, err
		}
	}

	c.schemas = schemas

	return c, nil
}

func (c *Client) GetProviderSchema() (*typ.GetProviderSchemaResponse, typ.Diagnostics) {
	return &c.schemas, nil
}

func (c *Client) ValidateProviderConfig(ctx context.Context, request typ.ValidateProviderConfigRequest) (*typ.ValidateProviderConfigResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	ty := c.schemas.ProviderCty

	mp, err := msgpack.Marshal(request.Config, ty)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	resp, err := c.client.ValidateProviderConfig(ctx, &tfprotov6.ValidateProviderConfigRequest{
		Config: &tfprotov6.DynamicValue{
			MsgPack: mp,
		},
	})
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	config, err := decodeDynamicValue(resp.PreparedConfig, ty)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}

	return &typ.ValidateProviderConfigResponse{
		PreparedConfig: config,
	}, diags
}

func (c *Client) ValidateResourceConfig(ctx context.Context, request typ.ValidateResourceConfigRequest) (*typ.ValidateResourceConfigResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	schema := c.schemas
	resourceTyp, ok := schema.ResourceTypesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	mp, err := msgpack.Marshal(request.Config, resourceTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	resp, err := c.client.ValidateResourceConfig(ctx, &tfprotov6.ValidateResourceConfigRequest{
		TypeName: request.TypeName,
		Config:   &tfprotov6.DynamicValue{MsgPack: mp},
	})
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	return &typ.ValidateResourceConfigResponse{}, diags
}

func (c *Client) ValidateDataResourceConfig(ctx context.Context, request typ.ValidateDataResourceConfigRequest) (*typ.ValidateDataResourceConfigResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	schema := c.schemas
	datasourceTyp, ok := schema.DataSourcesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown data source type %q", request.TypeName))...)
		return nil, diags
	}

	mp, err := msgpack.Marshal(request.Config, datasourceTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	resp, err := c.client.ValidateDataResourceConfig(ctx, &tfprotov6.ValidateDataResourceConfigRequest{
		TypeName: request.TypeName,
		Config:   &tfprotov6.DynamicValue{MsgPack: mp},
	})
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	return &typ.ValidateDataResourceConfigResponse{}, diags
}

func (c *Client) UpgradeResourceState(ctx context.Context, request typ.UpgradeResourceStateRequest) (*typ.UpgradeResourceStateResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	schema := c.schemas

	resTyp, ok := schema.ResourceTypesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	protoReq := &tfprotov6.UpgradeResourceStateRequest{
		TypeName: request.TypeName,
		Version:  int64(request.Version),
		RawState: &tfprotov6.RawState{
			JSON:    request.RawStateJSON,
			Flatmap: request.RawStateFlatmap,
		},
	}

	resp, err := c.client.UpgradeResourceState(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state := cty.NullVal(resTyp)
	if resp.UpgradedState != nil {
		state, err = decodeDynamicValue(resp.UpgradedState, resTyp)
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("decode dynamic value", err)...)
			return nil, diags
		}
	}
	return &typ.UpgradeResourceStateResponse{
		UpgradedState: state,
	}, diags
}

func (c *Client) ConfigureProvider(ctx context.Context, request typ.ConfigureProviderRequest) (*typ.ConfigureProviderResponse, typ.Diagnostics) {
	c.configuredMu.Lock()
	defer c.configuredMu.Unlock()
	if c.configured {
		return nil, typ.Diagnostics{
			{
				Severity: typ.Error,
				Summary:  "Provider already configured",
				Detail:   "This operation requires an unconfigured provider, but this provider was already configured.",
			},
		}
	}

	var diags typ.Diagnostics

	schema := c.schemas
	mp, err := msgpack.Marshal(
		request.Config,
		schema.ProviderCty,
	)
	if err != nil {
		diags := typ.ErrorDiagnostics("msgpack marshal", err)
		return nil, diags
	}
	resp, err := c.client.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		TerraformVersion: request.TerraformVersion,
		Config: &tfprotov6.DynamicValue{
			MsgPack: mp,
		},
		ClientCapabilities: &tfprotov6.ConfigureProviderClientCapabilities{
			DeferralAllowed: request.ClientCapabilities.DeferralAllowed,
		},
	})
	if err != nil {
		diags := typ.RPCErrorDiagnostics(err)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}
	c.configured = true
	return &typ.ConfigureProviderResponse{}, nil
}

func (c *Client) Stop(ctx context.Context) error {
	resp, err := c.client.StopProvider(ctx, &tfprotov6.StopProviderRequest{})
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return nil
}

func (c *Client) ReadResource(ctx context.Context, request typ.ReadResourceRequest) (*typ.ReadResourceResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	schema := c.schemas

	resTyp, ok := schema.ResourceTypesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaTyp := schema.ProviderMetaCty

	mp, err := msgpack.Marshal(request.PriorState, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov6.ReadResourceRequest{
		TypeName:     request.TypeName,
		CurrentState: &tfprotov6.DynamicValue{MsgPack: mp},
		Private:      request.Private,
		ClientCapabilities: &tfprotov6.ReadResourceClientCapabilities{
			DeferralAllowed: request.ClientCapabilities.DeferralAllowed,
		},
	}

	// The second check here is not something from terraform's implementation, should be derived from the schema drift in tfjson module.
	//if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
	if !metaTyp.Equals(cty.EmptyObject) {
		metaMP, err := msgpack.Marshal(request.ProviderMeta, metaTyp)
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
			return nil, diags
		}
		protoReq.ProviderMeta = &tfprotov6.DynamicValue{MsgPack: metaMP}
	}

	resp, err := c.client.ReadResource(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(resp.NewState, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}

	return &typ.ReadResourceResponse{
		NewState: state,
		Private:  resp.Private,
		Deferred: convert.ProtoToDeferred(resp.Deferred),
	}, diags
}

func (c *Client) PlanResourceChange(ctx context.Context, request typ.PlanResourceChangeRequest) (*typ.PlanResourceChangeResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	schema := c.schemas

	resTyp, ok := schema.ResourceTypesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaTyp := schema.ProviderMetaCty
	capabilities := schema.ServerCapabilities

	var response typ.PlanResourceChangeResponse

	// If the provider doesn't support planning a destroy operation, we can
	// return immediately.
	if request.ProposedNewState.IsNull() && !capabilities.PlanDestroy {
		response.PlannedState = request.ProposedNewState
		response.PlannedPrivate = request.PriorPrivate
		return &response, nil
	}

	priorMP, err := msgpack.Marshal(request.PriorState, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	configMP, err := msgpack.Marshal(request.Config, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	propMP, err := msgpack.Marshal(request.ProposedNewState, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov6.PlanResourceChangeRequest{
		TypeName:         request.TypeName,
		PriorState:       &tfprotov6.DynamicValue{MsgPack: priorMP},
		Config:           &tfprotov6.DynamicValue{MsgPack: configMP},
		ProposedNewState: &tfprotov6.DynamicValue{MsgPack: propMP},
		PriorPrivate:     request.PriorPrivate,
		ClientCapabilities: &tfprotov6.PlanResourceChangeClientCapabilities{
			DeferralAllowed: request.ClientCapabilities.DeferralAllowed,
		},
	}

	// The second check here is not something from terraform's implementation, should be derived from the schema drift in tfjson module.
	//if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
	if !metaTyp.Equals(cty.EmptyObject) {
		metaVal := request.ProviderMeta
		if metaVal == cty.NilVal {
			metaVal = cty.NullVal(metaTyp)
		}
		metaMP, err := msgpack.Marshal(metaVal, metaTyp)
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
			return nil, diags
		}
		protoReq.ProviderMeta = &tfprotov6.DynamicValue{MsgPack: metaMP}
	}

	protoResp, err := c.client.PlanResourceChange(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(protoResp.PlannedState, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}
	response.PlannedState = state

	for _, p := range protoResp.RequiresReplace {
		response.RequiresReplace = append(response.RequiresReplace, convert.DecodeAttributePath(p))
	}

	response.PlannedPrivate = protoResp.PlannedPrivate

	response.LegacyTypeSystem = protoResp.UnsafeToUseLegacyTypeSystem

	response.Deferred = convert.ProtoToDeferred(protoResp.Deferred)

	return &response, diags
}

func (c *Client) ApplyResourceChange(ctx context.Context, request typ.ApplyResourceChangeRequest) (*typ.ApplyResourceChangeResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	schema := c.schemas

	resTyp, ok := schema.ResourceTypesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaTyp := schema.ProviderMetaCty

	priorMP, err := msgpack.Marshal(request.PriorState, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}
	plannedMP, err := msgpack.Marshal(request.PlannedState, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}
	configMP, err := msgpack.Marshal(request.Config, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov6.ApplyResourceChangeRequest{
		TypeName:       request.TypeName,
		PriorState:     &tfprotov6.DynamicValue{MsgPack: priorMP},
		PlannedState:   &tfprotov6.DynamicValue{MsgPack: plannedMP},
		Config:         &tfprotov6.DynamicValue{MsgPack: configMP},
		PlannedPrivate: request.PlannedPrivate,
	}

	// The second check here is not something from terraform's implementation, should be derived from the schema drift in tfjson module.
	//if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
	if !metaTyp.Equals(cty.EmptyObject) {
		metaVal := request.ProviderMeta
		if metaVal == cty.NilVal {
			metaVal = cty.NullVal(metaTyp)
		}
		metaMP, err := msgpack.Marshal(metaVal, metaTyp)
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
			return nil, diags
		}
		protoReq.ProviderMeta = &tfprotov6.DynamicValue{MsgPack: metaMP}
	}

	protoResp, err := c.client.ApplyResourceChange(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(protoResp.NewState, resTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	return &typ.ApplyResourceChangeResponse{
		NewState:         state,
		Private:          protoResp.Private,
		LegacyTypeSystem: protoResp.UnsafeToUseLegacyTypeSystem,
	}, diags
}

func (c *Client) ImportResourceState(ctx context.Context, request typ.ImportResourceStateRequest) (*typ.ImportResourceStateResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	schema := c.schemas
	resp, err := c.client.ImportResourceState(ctx, &tfprotov6.ImportResourceStateRequest{
		TypeName: request.TypeName,
		ID:       request.ID,
		ClientCapabilities: &tfprotov6.ImportResourceStateClientCapabilities{
			DeferralAllowed: request.ClientCapabilities.DeferralAllowed,
		},
	})
	if err != nil {
		return nil, typ.RPCErrorDiagnostics(err)
	}

	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	var response typ.ImportResourceStateResponse
	response.Deferred = convert.ProtoToDeferred(resp.Deferred)
	for _, imported := range resp.ImportedResources {
		resource := typ.ImportedResource{
			TypeName: imported.TypeName,
			Private:  imported.Private,
		}

		resTyp, ok := schema.ResourceTypesCty[imported.TypeName]
		if !ok {
			diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", imported.TypeName))...)
			continue
		}

		state, err := decodeDynamicValue(imported.State, resTyp)
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("decode dynamic value", err)...)
			return nil, diags
		}
		resource.State = state
		response.ImportedResources = append(response.ImportedResources, resource)

	}

	return &response, diags
}

func (c *Client) MoveResourceState(ctx context.Context, request typ.MoveResourceStateRequest) (*typ.MoveResourceStateResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	schema := c.schemas

	protoReq := &tfprotov6.MoveResourceStateRequest{
		SourceProviderAddress: request.SourceProviderAddress,
		SourceTypeName:        request.SourceTypeName,
		SourceSchemaVersion:   request.SourceSchemaVersion,
		SourceState: &tfprotov6.RawState{
			JSON: request.SourceStateJSON,
		},
		SourcePrivate:  request.SourcePrivate,
		TargetTypeName: request.TargetTypeName,
	}

	protoResp, err := c.client.MoveResourceState(ctx, protoReq)
	if err != nil {
		return nil, typ.RPCErrorDiagnostics(err)
	}

	targetTyp, ok := schema.ResourceTypesCty[request.TargetTypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TargetTypeName))...)
		return nil, diags
	}
	state, err := decodeDynamicValue(protoResp.TargetState, targetTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}

	return &typ.MoveResourceStateResponse{
		TargetState:   state,
		TargetPrivate: protoResp.TargetPrivate,
	}, nil
}

func (c *Client) ReadDataSource(ctx context.Context, request typ.ReadDataSourceRequest) (*typ.ReadDataSourceResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	schema := c.schemas

	dstTyp, ok := schema.DataSourcesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown data source type %q", request.TypeName))...)
		return nil, diags
	}

	metaTyp := schema.ProviderMetaCty

	mp, err := msgpack.Marshal(request.Config, dstTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov6.ReadDataSourceRequest{
		TypeName: request.TypeName,
		Config:   &tfprotov6.DynamicValue{MsgPack: mp},
		ClientCapabilities: &tfprotov6.ReadDataSourceClientCapabilities{
			DeferralAllowed: request.ClientCapabilities.DeferralAllowed,
		},
	}

	// The second check here is not something from terraform's implementation, should be derived from the schema drift in tfjson module.
	//if metaTyp.Block != nil && len(metaTyp.Block.NestedBlocks)+len(metaTyp.Block.Attributes) != 0 {
	if !metaTyp.Equals(cty.EmptyObject) {
		metaMP, err := msgpack.Marshal(request.ProviderMeta, metaTyp)
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
			return nil, diags
		}
		protoReq.ProviderMeta = &tfprotov6.DynamicValue{MsgPack: metaMP}
	}

	resp, err := c.client.ReadDataSource(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(resp.State, dstTyp)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}

	return &typ.ReadDataSourceResponse{
		State:    state,
		Deferred: convert.ProtoToDeferred(resp.Deferred),
	}, diags
}

func (c *Client) CallFunction(ctx context.Context, request typ.CallFunctionRequest) (*typ.CallFunctionResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	schema := c.schemas

	funcDecl, ok := schema.Functions[request.FunctionName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown function name type %q", request.FunctionName))...)
		return nil, diags
	}
	if len(request.Arguments) < len(funcDecl.Parameters) {
		diags = append(diags, typ.ErrorDiagnostics("call function error", fmt.Errorf("not enough arguments for function %q", request.FunctionName))...)
		return nil, diags
	}
	if funcDecl.VariadicParameter == nil && len(request.Arguments) > len(funcDecl.Parameters) {
		diags = append(diags, typ.ErrorDiagnostics("call function error", fmt.Errorf("too many arguments for function %q", request.FunctionName))...)
		return nil, diags
	}
	args := make([]*tfprotov6.DynamicValue, len(request.Arguments))
	for i, argVal := range request.Arguments {
		var paramDecl typ.FunctionParam
		if i < len(funcDecl.Parameters) {
			paramDecl = funcDecl.Parameters[i]
		} else {
			paramDecl = *funcDecl.VariadicParameter
		}

		argValRaw, err := msgpack.Marshal(argVal, paramDecl.Type)
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("call function error", fmt.Errorf("marshal argument: %v", err))...)
			return nil, diags
		}
		args[i] = &tfprotov6.DynamicValue{
			MsgPack: argValRaw,
		}
	}

	protoResp, err := c.client.CallFunction(ctx, &tfprotov6.CallFunctionRequest{
		Name:      request.FunctionName,
		Arguments: args,
	})
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}

	resp := &typ.CallFunctionResponse{}

	if protoResp.Error != nil {
		resp.Err = errors.New(protoResp.Error.Text)

		// If this is a problem with a specific argument, we can wrap the error
		// in a function.ArgError
		if protoResp.Error.FunctionArgument != nil {
			resp.Err = function.NewArgError(int(*protoResp.Error.FunctionArgument), resp.Err)
		}

		return resp, nil
	}

	resultVal, err := decodeDynamicValue(protoResp.Result, funcDecl.ReturnType)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("call function error", fmt.Errorf("decoding return value: %v", err))...)
		return nil, diags
	}

	resp.Result = resultVal
	return resp, nil
}

// GetResourceIdentitySchemas implements tfclient.Client.
func (c *Client) GetResourceIdentitySchemas() *typ.GetResourceIdentitySchemasResponse {
	panic("unimplemented")
}

// ValidateListResourceConfig implements tfclient.Client.
func (c *Client) ValidateListResourceConfig(context.Context, typ.ValidateListResourceConfigRequest) typ.Diagnostics {
	panic("unimplemented")
}

// UpgradeResourceIdentity implements tfclient.Client.
func (c *Client) UpgradeResourceIdentity(context.Context, typ.UpgradeResourceIdentityRequest) (*typ.UpgradeResourceIdentityResponse, typ.Diagnostics) {
	panic("unimplemented")
}

// ValidateEphemeralResourceConfig implements tfclient.Client.
func (c *Client) ValidateEphemeralResourceConfig(context.Context, typ.ValidateEphemeralResourceConfigRequest) typ.Diagnostics {
	panic("unimplemented")
}

// OpenEphemeralResource implements tfclient.Client.
func (c *Client) OpenEphemeralResource(context.Context, typ.OpenEphemeralResourceRequest) (*typ.OpenEphemeralResourceResponse, typ.Diagnostics) {
	panic("unimplemented")
}

// CloseEphemeralResource implements tfclient.Client.
func (c *Client) CloseEphemeralResource(context.Context, typ.CloseEphemeralResourceRequest) typ.Diagnostics {
	panic("unimplemented")
}

// RenewEphemeralResource implements tfclient.Client.
func (c *Client) RenewEphemeralResource(context.Context, typ.RenewEphemeralResourceRequest) (*typ.RenewEphemeralResourceResponse, typ.Diagnostics) {
	panic("unimplemented")
}

func (c *Client) Close() {
	c.pluginClient.Kill()
}

// Decode a DynamicValue from either the JSON or MsgPack encoding.
// Derived from github.com/hashicorp/terraform/internal/plugin6/grpc_provider.go (15ecdb66c84cd8202b0ae3d34c44cb4bbece5444)
func decodeDynamicValue(v *tfprotov6.DynamicValue, ty cty.Type) (cty.Value, error) {
	// always return a valid value
	var err error
	res := cty.NullVal(ty)
	if v == nil {
		return res, nil
	}

	switch {
	case len(v.MsgPack) > 0:
		res, err = msgpack.Unmarshal(v.MsgPack, ty)
	case len(v.JSON) > 0:
		res, err = ctyjson.Unmarshal(v.JSON, ty)
	}
	return res, err
}
