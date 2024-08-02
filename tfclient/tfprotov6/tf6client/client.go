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
		ResourceTypes: map[string]tfjson.Schema{},
		DataSources:   map[string]tfjson.Schema{},
		ServerCapabilities: typ.ServerCapabilities{
			PlanDestroy: false,
		},
	}
	if resp.Provider != nil {
		schemas.Provider = convert.ProtoToProviderSchema(resp.Provider)
	}
	if resp.ProviderMeta != nil {
		schemas.ProviderMeta = convert.ProtoToProviderSchema(resp.ProviderMeta)
	}
	if resp.ServerCapabilities != nil {
		schemas.ServerCapabilities.PlanDestroy = resp.ServerCapabilities.PlanDestroy
	}
	for name, schema := range resp.ResourceSchemas {
		schemas.ResourceTypes[name] = convert.ProtoToProviderSchema(schema)
	}
	for name, schema := range resp.DataSourceSchemas {
		schemas.DataSources[name] = convert.ProtoToProviderSchema(schema)
	}

	c.schemas = schemas

	return c, nil
}

func (c *Client) GetProviderSchema() (*typ.GetProviderSchemaResponse, typ.Diagnostics) {
	return &c.schemas, nil
}

func (c *Client) ValidateProviderConfig(ctx context.Context, request typ.ValidateProviderConfigRequest) (*typ.ValidateProviderConfigResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	ty := configschema.SchemaBlockImpliedType(c.schemas.Provider.Block)

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
	resourceSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	mp, err := msgpack.Marshal(request.Config, configschema.SchemaBlockImpliedType(resourceSchema.Block))
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
	datasourceSchema, ok := schema.DataSources[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown data source type %q", request.TypeName))...)
		return nil, diags
	}

	mp, err := msgpack.Marshal(request.Config, configschema.SchemaBlockImpliedType(datasourceSchema.Block))
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

	resSchema, ok := schema.ResourceTypes[request.TypeName]
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

	ty := configschema.SchemaBlockImpliedType(resSchema.Block)
	state := cty.NullVal(ty)
	if resp.UpgradedState != nil {
		state, err = decodeDynamicValue(resp.UpgradedState, ty)
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

	schema := c.schemas
	mp, err := msgpack.Marshal(
		request.Config,
		configschema.SchemaBlockImpliedType(schema.Provider.Block),
	)
	if err != nil {
		diags := typ.ErrorDiagnostics("msgpack marshal", err)
		return nil, diags
	}
	if _, err := c.client.ConfigureProvider(ctx, &tfprotov6.ConfigureProviderRequest{
		TerraformVersion: request.TerraformVersion,
		Config: &tfprotov6.DynamicValue{
			MsgPack: mp,
		},
		ClientCapabilities: &tfprotov6.ConfigureProviderClientCapabilities{
			DeferralAllowed: request.ClientCapabilities.DeferralAllowed,
		},
	}); err != nil {
		diags := typ.RPCErrorDiagnostics(err)
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

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaSchema := schema.ProviderMeta

	mp, err := msgpack.Marshal(request.PriorState, configschema.SchemaBlockImpliedType(resSchema.Block))
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
	if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
		metaMP, err := msgpack.Marshal(request.ProviderMeta, configschema.SchemaBlockImpliedType(metaSchema.Block))
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

	state, err := decodeDynamicValue(resp.NewState, configschema.SchemaBlockImpliedType(resSchema.Block))
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

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaSchema := schema.ProviderMeta
	capabilities := schema.ServerCapabilities

	var response typ.PlanResourceChangeResponse

	// If the provider doesn't support planning a destroy operation, we can
	// return immediately.
	if request.ProposedNewState.IsNull() && !capabilities.PlanDestroy {
		response.PlannedState = request.ProposedNewState
		response.PlannedPrivate = request.PriorPrivate
		return &response, nil
	}

	priorMP, err := msgpack.Marshal(request.PriorState, configschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	configMP, err := msgpack.Marshal(request.Config, configschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	propMP, err := msgpack.Marshal(request.ProposedNewState, configschema.SchemaBlockImpliedType(resSchema.Block))
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
	if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
		metaTy := configschema.SchemaBlockImpliedType(metaSchema.Block)
		metaVal := request.ProviderMeta
		if metaVal == cty.NilVal {
			metaVal = cty.NullVal(metaTy)
		}
		metaMP, err := msgpack.Marshal(metaVal, metaTy)
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

	state, err := decodeDynamicValue(protoResp.PlannedState, configschema.SchemaBlockImpliedType(resSchema.Block))
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

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaSchema := schema.ProviderMeta

	priorMP, err := msgpack.Marshal(request.PriorState, configschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}
	plannedMP, err := msgpack.Marshal(request.PlannedState, configschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}
	configMP, err := msgpack.Marshal(request.Config, configschema.SchemaBlockImpliedType(resSchema.Block))
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
	if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
		metaTy := configschema.SchemaBlockImpliedType(metaSchema.Block)
		metaVal := request.ProviderMeta
		if metaVal == cty.NilVal {
			metaVal = cty.NullVal(metaTy)
		}
		metaMP, err := msgpack.Marshal(metaVal, metaTy)
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

	state, err := decodeDynamicValue(protoResp.NewState, configschema.SchemaBlockImpliedType(metaSchema.Block))
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

		resSchema, ok := schema.ResourceTypes[imported.TypeName]
		if !ok {
			diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", imported.TypeName))...)
			continue
		}

		state, err := decodeDynamicValue(imported.State, configschema.SchemaBlockImpliedType(resSchema.Block))
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

	targetType, ok := schema.ResourceTypes[request.TargetTypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TargetTypeName))...)
		return nil, diags
	}
	state, err := decodeDynamicValue(protoResp.TargetState, configschema.SchemaBlockImpliedType(targetType.Block))
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

	dsSchema, ok := schema.DataSources[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown data source type %q", request.TypeName))...)
		return nil, diags
	}

	metaSchema := schema.ProviderMeta

	mp, err := msgpack.Marshal(request.Config, configschema.SchemaBlockImpliedType(dsSchema.Block))
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
	if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
		metaMP, err := msgpack.Marshal(request.ProviderMeta, configschema.SchemaBlockImpliedType(metaSchema.Block))
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

	state, err := decodeDynamicValue(resp.State, configschema.SchemaBlockImpliedType(dsSchema.Block))
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
