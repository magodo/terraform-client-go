// This is derived from github.com/hashicorp/terraform/internal/plugin5/grpc_provider.go (15ecdb66c84cd8202b0ae3d34c44cb4bbece5444)

package tf5client

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/hashicorp/go-plugin"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/client"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/convert"
	"github.com/magodo/tfstate/terraform/jsonschema"
	"github.com/zclconf/go-cty/cty"
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
	client tfprotov5.ProviderServer

	// schema stores the schema for this provider. This is used to properly
	// serialize the state for requests.
	schemas client.GetProviderSchemaResponse

	configured   bool
	configuredMu sync.Mutex
}

func New(pluginClient *plugin.Client, grpcClient tfprotov5.ProviderServer) (client.Interface, error) {
	c := &Client{
		pluginClient: pluginClient,
		client:       grpcClient,
	}

	resp, err := grpcClient.GetProviderSchema(context.Background(), &tfprotov5.GetProviderSchemaRequest{})
	if err != nil {
		return nil, err
	}
	if diags := convert.DecodeDiagnostics(resp.Diagnostics); diags.HasErrors() {
		return nil, diags.Err()
	}

	schemas := client.GetProviderSchemaResponse{
		Provider:      convert.ProtoToProviderSchema(resp.Provider),
		ProviderMeta:  convert.ProtoToProviderSchema(resp.ProviderMeta),
		ResourceTypes: map[string]tfjson.Schema{},
		DataSources:   map[string]tfjson.Schema{},
		ServerCapabilities: client.ServerCapabilities{
			PlanDestroy: false,
		},
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

func (c *Client) GetProviderSchema() (*client.GetProviderSchemaResponse, client.Diagnostics) {
	return &c.schemas, nil
}

func (c *Client) ValidateProviderConfig(ctx context.Context, request client.ValidateProviderConfigRequest) (*client.ValidateProviderConfigResponse, client.Diagnostics) {
	var diags client.Diagnostics

	ty := jsonschema.SchemaBlockImpliedType(c.schemas.Provider.Block)

	mp, err := msgpack.Marshal(request.Config, ty)
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	resp, err := c.client.PrepareProviderConfig(ctx, &tfprotov5.PrepareProviderConfigRequest{
		Config: &tfprotov5.DynamicValue{
			MsgPack: mp,
		},
	})
	if err != nil {
		diags = append(diags, client.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	config, err := decodeDynamicValue(resp.PreparedConfig, ty)
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}

	return &client.ValidateProviderConfigResponse{
		PreparedConfig: config,
	}, diags
}

func (c *Client) ValidateResourceConfig(ctx context.Context, request client.ValidateResourceConfigRequest) (*client.ValidateResourceConfigResponse, client.Diagnostics) {
	var diags client.Diagnostics

	schema := c.schemas
	resourceSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, client.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	mp, err := msgpack.Marshal(request.Config, jsonschema.SchemaBlockImpliedType(resourceSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	resp, err := c.client.ValidateResourceTypeConfig(ctx, &tfprotov5.ValidateResourceTypeConfigRequest{
		TypeName: request.TypeName,
		Config:   &tfprotov5.DynamicValue{MsgPack: mp},
	})
	if err != nil {
		diags = append(diags, client.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	return &client.ValidateResourceConfigResponse{}, diags
}

func (c *Client) ValidateDataResourceConfig(ctx context.Context, request client.ValidateDataResourceConfigRequest) (*client.ValidateDataResourceConfigResponse, client.Diagnostics) {
	var diags client.Diagnostics

	schema := c.schemas
	datasourceSchema, ok := schema.DataSources[request.TypeName]
	if !ok {
		diags = append(diags, client.ErrorDiagnostics("no schema", fmt.Errorf("unknown data source type %q", request.TypeName))...)
		return nil, diags
	}

	mp, err := msgpack.Marshal(request.Config, jsonschema.SchemaBlockImpliedType(datasourceSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	resp, err := c.client.ValidateDataSourceConfig(ctx, &tfprotov5.ValidateDataSourceConfigRequest{
		TypeName: request.TypeName,
		Config:   &tfprotov5.DynamicValue{MsgPack: mp},
	})
	if err != nil {
		diags = append(diags, client.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	return &client.ValidateDataResourceConfigResponse{}, diags
}

func (c *Client) UpgradeResourceState(ctx context.Context, request client.UpgradeResourceStateRequest) (*client.UpgradeResourceStateResponse, client.Diagnostics) {
	var diags client.Diagnostics

	schema := c.schemas

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, client.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	protoReq := &tfprotov5.UpgradeResourceStateRequest{
		TypeName: request.TypeName,
		Version:  int64(request.Version),
		RawState: &tfprotov5.RawState{
			JSON:    request.RawStateJSON,
			Flatmap: request.RawStateFlatmap,
		},
	}

	resp, err := c.client.UpgradeResourceState(ctx, protoReq)
	if err != nil {
		diags = append(diags, client.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	ty := jsonschema.SchemaBlockImpliedType(resSchema.Block)
	state := cty.NullVal(ty)
	if resp.UpgradedState != nil {
		state, err = decodeDynamicValue(resp.UpgradedState, ty)
		if err != nil {
			diags = append(diags, client.ErrorDiagnostics("decode dynamic value", err)...)
			return nil, diags
		}
	}
	return &client.UpgradeResourceStateResponse{
		UpgradedState: state,
	}, diags
}

func (c *Client) ConfigureProvider(ctx context.Context, request client.ConfigureProviderRequest) (*client.ConfigureProviderResponse, client.Diagnostics) {
	c.configuredMu.Lock()
	defer c.configuredMu.Unlock()
	if c.configured {
		return nil, client.Diagnostics{
			{
				Severity: client.Error,
				Summary:  "Provider already configured",
				Detail:   "This operation requires an unconfigured provider, but this provider was already configured.",
			},
		}
	}

	var diags client.Diagnostics

	schema := c.schemas
	mp, err := msgpack.Marshal(
		request.Config,
		jsonschema.SchemaBlockImpliedType(schema.Provider.Block),
	)
	if err != nil {
		diags := client.ErrorDiagnostics("msgpack marshal", err)
		return nil, diags
	}
	resp, err := c.client.ConfigureProvider(ctx, &tfprotov5.ConfigureProviderRequest{
		TerraformVersion: request.TerraformVersion,
		Config: &tfprotov5.DynamicValue{
			MsgPack: mp,
		},
	})
	if err != nil {
		diags := client.RPCErrorDiagnostics(err)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	c.configured = true

	return &client.ConfigureProviderResponse{}, diags
}

func (c *Client) Stop(ctx context.Context) error {
	resp, err := c.client.StopProvider(ctx, &tfprotov5.StopProviderRequest{})
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error)
	}
	return nil
}

func (c *Client) ReadResource(ctx context.Context, request client.ReadResourceRequest) (*client.ReadResourceResponse, client.Diagnostics) {
	var diags client.Diagnostics
	schema := c.schemas

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, client.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaSchema := schema.ProviderMeta

	mp, err := msgpack.Marshal(request.PriorState, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov5.ReadResourceRequest{
		TypeName:     request.TypeName,
		CurrentState: &tfprotov5.DynamicValue{MsgPack: mp},
		Private:      request.Private,
	}

	// The second check here is not something from terraform's implementation, should be derived from the schema drift in tfjson module.
	if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
		metaMP, err := msgpack.Marshal(request.ProviderMeta, jsonschema.SchemaBlockImpliedType(metaSchema.Block))
		if err != nil {
			diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
			return nil, diags
		}
		protoReq.ProviderMeta = &tfprotov5.DynamicValue{MsgPack: metaMP}
	}

	resp, err := c.client.ReadResource(ctx, protoReq)
	if err != nil {
		diags = append(diags, client.RPCErrorDiagnostics(err)...)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(resp.NewState, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}

	return &client.ReadResourceResponse{
		NewState: state,
		Private:  resp.Private,
	}, diags
}

func (c *Client) PlanResourceChange(ctx context.Context, request client.PlanResourceChangeRequest) (*client.PlanResourceChangeResponse, client.Diagnostics) {
	var diags client.Diagnostics
	schema := c.schemas

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, client.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaSchema := schema.ProviderMeta
	capabilities := schema.ServerCapabilities

	var response client.PlanResourceChangeResponse

	// If the provider doesn't support planning a destroy operation, we can
	// return immediately.
	if request.ProposedNewState.IsNull() && !capabilities.PlanDestroy {
		response.PlannedState = request.ProposedNewState
		response.PlannedPrivate = request.PriorPrivate
		return &response, nil
	}

	priorMP, err := msgpack.Marshal(request.PriorState, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	configMP, err := msgpack.Marshal(request.Config, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	propMP, err := msgpack.Marshal(request.ProposedNewState, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov5.PlanResourceChangeRequest{
		TypeName:         request.TypeName,
		PriorState:       &tfprotov5.DynamicValue{MsgPack: priorMP},
		Config:           &tfprotov5.DynamicValue{MsgPack: configMP},
		ProposedNewState: &tfprotov5.DynamicValue{MsgPack: propMP},
		PriorPrivate:     request.PriorPrivate,
	}

	// The second check here is not something from terraform's implementation, should be derived from the schema drift in tfjson module.
	if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
		metaMP, err := msgpack.Marshal(request.ProviderMeta, jsonschema.SchemaBlockImpliedType(resSchema.Block))
		if err != nil {
			diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
			return nil, diags
		}
		protoReq.ProviderMeta = &tfprotov5.DynamicValue{MsgPack: metaMP}
	}

	protoResp, err := c.client.PlanResourceChange(ctx, protoReq)
	if err != nil {
		diags = append(diags, client.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(protoResp.PlannedState, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}
	response.PlannedState = state

	for _, p := range protoResp.RequiresReplace {
		response.RequiresReplace = append(response.RequiresReplace, convert.DecodeAttributePath(p))
	}

	response.PlannedPrivate = protoResp.PlannedPrivate

	response.LegacyTypeSystem = protoResp.UnsafeToUseLegacyTypeSystem

	return &response, diags
}

func (c *Client) ApplyResourceChange(ctx context.Context, request client.ApplyResourceChangeRequest) (*client.ApplyResourceChangeResponse, client.Diagnostics) {
	var diags client.Diagnostics
	schema := c.schemas

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, client.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaSchema := schema.ProviderMeta

	priorMP, err := msgpack.Marshal(request.PriorState, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}
	plannedMP, err := msgpack.Marshal(request.PlannedState, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}
	configMP, err := msgpack.Marshal(request.Config, jsonschema.SchemaBlockImpliedType(resSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       request.TypeName,
		PriorState:     &tfprotov5.DynamicValue{MsgPack: priorMP},
		PlannedState:   &tfprotov5.DynamicValue{MsgPack: plannedMP},
		Config:         &tfprotov5.DynamicValue{MsgPack: configMP},
		PlannedPrivate: request.PlannedPrivate,
	}

	// The second check here is not something from terraform's implementation, should be derived from the schema drift in tfjson module.
	if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
		metaMP, err := msgpack.Marshal(request.ProviderMeta, jsonschema.SchemaBlockImpliedType(metaSchema.Block))
		if err != nil {
			diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
			return nil, diags
		}
		protoReq.ProviderMeta = &tfprotov5.DynamicValue{MsgPack: metaMP}
	}

	protoResp, err := c.client.ApplyResourceChange(ctx, protoReq)
	if err != nil {
		diags = append(diags, client.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(protoResp.NewState, jsonschema.SchemaBlockImpliedType(metaSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	return &client.ApplyResourceChangeResponse{
		NewState:         state,
		Private:          protoResp.Private,
		LegacyTypeSystem: protoResp.UnsafeToUseLegacyTypeSystem,
	}, diags
}

func (c *Client) ImportResourceState(ctx context.Context, request client.ImportResourceStateRequest) (*client.ImportResourceStateResponse, client.Diagnostics) {
	var diags client.Diagnostics

	schema := c.schemas
	resp, err := c.client.ImportResourceState(ctx, &tfprotov5.ImportResourceStateRequest{
		TypeName: request.TypeName,
		ID:       request.ID,
	})
	if err != nil {
		return nil, client.RPCErrorDiagnostics(err)
	}

	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	var response client.ImportResourceStateResponse
	for _, imported := range resp.ImportedResources {
		resource := client.ImportedResource{
			TypeName: imported.TypeName,
			Private:  imported.Private,
		}

		resSchema, ok := schema.ResourceTypes[imported.TypeName]
		if !ok {
			diags = append(diags, client.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", imported.TypeName))...)
			continue
		}

		state, err := decodeDynamicValue(imported.State, jsonschema.SchemaBlockImpliedType(resSchema.Block))
		if err != nil {
			diags = append(diags, client.ErrorDiagnostics("decode dynamic value", err)...)
			return nil, diags
		}
		resource.State = state
		response.ImportedResources = append(response.ImportedResources, resource)

	}

	return &response, diags
}

func (c *Client) ReadDataSource(ctx context.Context, request client.ReadDataSourceRequest) (*client.ReadDataSourceResponse, client.Diagnostics) {
	var diags client.Diagnostics
	schema := c.schemas

	dsSchema, ok := schema.DataSources[request.TypeName]
	if !ok {
		diags = append(diags, client.ErrorDiagnostics("no schema", fmt.Errorf("unknown data source type %q", request.TypeName))...)
		return nil, diags
	}

	metaSchema := schema.ProviderMeta

	mp, err := msgpack.Marshal(request.Config, jsonschema.SchemaBlockImpliedType(dsSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov5.ReadDataSourceRequest{
		TypeName: request.TypeName,
		Config:   &tfprotov5.DynamicValue{MsgPack: mp},
	}

	// The second check here is not something from terraform's implementation, should be derived from the schema drift in tfjson module.
	if metaSchema.Block != nil && len(metaSchema.Block.NestedBlocks)+len(metaSchema.Block.Attributes) != 0 {
		metaMP, err := msgpack.Marshal(request.ProviderMeta, jsonschema.SchemaBlockImpliedType(metaSchema.Block))
		if err != nil {
			diags = append(diags, client.ErrorDiagnostics("msgpack marshal", err)...)
			return nil, diags
		}
		protoReq.ProviderMeta = &tfprotov5.DynamicValue{MsgPack: metaMP}
	}

	resp, err := c.client.ReadDataSource(ctx, protoReq)
	if err != nil {
		diags = append(diags, client.RPCErrorDiagnostics(err)...)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(resp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(resp.State, jsonschema.SchemaBlockImpliedType(dsSchema.Block))
	if err != nil {
		diags = append(diags, client.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}

	return &client.ReadDataSourceResponse{
		State: state,
	}, diags
}

func (c *Client) Close() {
	c.pluginClient.Kill()
}

// Decode a DynamicValue from either the JSON or MsgPack encoding.
// Derived from github.com/hashicorp/terraform/internal/plugin/grpc_provider.go (15ecdb66c84cd8202b0ae3d34c44cb4bbece5444)
func decodeDynamicValue(v *tfprotov5.DynamicValue, ty cty.Type) (cty.Value, error) {
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
