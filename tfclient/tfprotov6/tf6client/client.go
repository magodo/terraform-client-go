// This is derived from github.com/hashicorp/terraform/internal/plugin6/grpc_provider.go

package tf6client

import (
	"context"
	"errors"
	"fmt"
	"slices"
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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TFProtoV6Client interface {
	tfprotov6.ProviderServer
	tfprotov6.ListResourceServer
	tfprotov6.ActionServer
}

// Client handles the client, or core side of the plugin rpc connection.
// The Client methods are mostly a translation layer between the
// terraform providers types and the grpc proto types, directly converting
// between the two.
type Client struct {
	// PluginClient provides a reference to the plugin.Client which controls the plugin process.
	// This allows the Client a way to shutdown the plugin process.
	pluginClient *plugin.Client

	// Proto client use to make the grpc service calls.
	client TFProtoV6Client

	// schema stores the schema for this provider. This is used to properly
	// serialize the state for requests.
	schemas typ.GetProviderSchemaResponse

	configured   bool
	configuredMu sync.Mutex
}

func New(pluginClient *plugin.Client, grpcClient TFProtoV6Client, schema *typ.GetProviderSchemaResponse) (*Client, error) {
	ctx := context.Background()
	c := &Client{
		pluginClient: pluginClient,
		client:       grpcClient,
	}

	if schema != nil {
		c.schemas = *schema
		return c, nil
	}

	resp, err := grpcClient.GetProviderSchema(ctx, &tfprotov6.GetProviderSchemaRequest{})
	if err != nil {
		return nil, err
	}
	if diags := convert.DecodeDiagnostics(resp.Diagnostics); diags.HasErrors() {
		return nil, diags.Err()
	}

	schemas := typ.GetProviderSchemaResponse{
		ResourceTypes:             map[string]tfjson.Schema{},
		ResourceTypesCty:          map[string]cty.Type{},
		DataSources:               map[string]tfjson.Schema{},
		DataSourcesCty:            map[string]cty.Type{},
		Functions:                 map[string]typ.FunctionDecl{},
		EphemeralResourceTypes:    map[string]tfjson.Schema{},
		EphemeralResourceTypesCty: map[string]cty.Type{},
		ListResourceTypes:         map[string]tfjson.Schema{},
		ListResourceTypesCty:      map[string]cty.Type{},
		Actions:                   map[string]tfjson.Schema{},
		ActionsCty:                map[string]cty.Type{},
	}

	identResp, err := grpcClient.GetResourceIdentitySchemas(ctx, new(tfprotov6.GetResourceIdentitySchemasRequest))
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			// We don't treat this as an error if older providers don't implement this method,
			// so we create an empty map for identity schemas
			identResp = &tfprotov6.GetResourceIdentitySchemasResponse{
				IdentitySchemas: map[string]*tfprotov6.ResourceIdentitySchema{},
			}
		} else {
			return nil, err
		}
	}

	if diags := convert.DecodeDiagnostics(identResp.Diagnostics); diags.HasErrors() {
		return nil, diags.Err()
	}

	if resp.Provider != nil {
		providerSchema := convert.ProtoToProviderSchema(resp.Provider, nil)
		schemas.Provider = providerSchema
		schemas.ProviderCty = configschema.SchemaBlockImpliedType(providerSchema.Block)
	}
	if resp.ProviderMeta != nil {
		providerMetaSchema := convert.ProtoToProviderSchema(resp.ProviderMeta, nil)
		schemas.ProviderMeta = providerMetaSchema
		schemas.ProviderMetaCty = configschema.SchemaBlockImpliedType(providerMetaSchema.Block)
	}
	if resp.ServerCapabilities != nil {
		schemas.ServerCapabilities.PlanDestroy = resp.ServerCapabilities.PlanDestroy
	}
	for name, schema := range resp.ResourceSchemas {
		id := identResp.IdentitySchemas[name]
		resourceSchema := convert.ProtoToProviderSchema(schema, id)
		schemas.ResourceTypes[name] = resourceSchema
		schemas.ResourceTypesCty[name] = configschema.SchemaBlockImpliedType(resourceSchema.Block)
	}
	for name, schema := range resp.DataSourceSchemas {
		dataSourceSchema := convert.ProtoToProviderSchema(schema, nil)
		schemas.DataSources[name] = dataSourceSchema
		schemas.DataSourcesCty[name] = configschema.SchemaBlockImpliedType(dataSourceSchema.Block)
	}
	for name, ephem := range resp.EphemeralResourceSchemas {
		ephemSchema := convert.ProtoToProviderSchema(ephem, nil)
		schemas.EphemeralResourceTypes[name] = ephemSchema
		schemas.EphemeralResourceTypesCty[name] = configschema.SchemaBlockImpliedType(ephemSchema.Block)
	}
	for name, fun := range resp.Functions {
		schemas.Functions[name], err = convert.FunctionDeclFromProto(fun)
		if err != nil {
			return nil, err
		}
	}
	for name, list := range resp.ListResourceSchemas {
		listSchema := convert.ProtoToProviderSchema(list, nil)
		schemas.ListResourceTypes[name] = listSchema
		schemas.ListResourceTypesCty[name] = configschema.SchemaBlockImpliedType(listSchema.Block)
	}

	for name, action := range resp.ActionSchemas {
		actionSchema := convert.ProtoToProviderSchema(action.Schema, nil)
		schemas.Actions[name] = actionSchema
		schemas.ActionsCty[name] = configschema.SchemaBlockImpliedType(actionSchema.Block)
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
		ClientCapabilities: &tfprotov6.ValidateResourceConfigClientCapabilities{
			WriteOnlyAttributesAllowed: request.ClientCapabilities.WriteOnlyAttributesAllowed,
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

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

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

	if !request.CurrentIdentity.IsNull() {
		if resSchema.Identity == nil {
			diags = append(diags, typ.ErrorDiagnostics("identity type not found", fmt.Errorf("identity type not found for resoruce type %s", request.TypeName))...)
			return nil, diags
		}
		currentIdentityMP, err := msgpack.Marshal(request.CurrentIdentity, configschema.SchemaNestedAttributeTypeImpliedType(resSchema.Identity))
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("msgpach marshal", err)...)
			return nil, diags
		}
		protoReq.CurrentIdentity = &tfprotov6.ResourceIdentityData{
			IdentityData: &tfprotov6.DynamicValue{MsgPack: currentIdentityMP},
		}
	}

	protoResp, err := c.client.ReadResource(ctx, protoReq)
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
		diags = append(diags, typ.ErrorDiagnostics("decode dynamic value", err)...)
		return nil, diags
	}

	resp := &typ.ReadResourceResponse{
		NewState: state,
		Private:  protoResp.Private,
		Deferred: convert.ProtoToDeferred(protoResp.Deferred),
	}

	if protoResp.NewIdentity != nil && protoResp.NewIdentity.IdentityData != nil {
		if resSchema.Identity == nil {
			diags = append(diags, typ.ErrorDiagnostics("unknown identity type", fmt.Errorf("unknown identity type %s", request.TypeName))...)
			return nil, diags
		}

		resp.Identity, err = decodeDynamicValue(protoResp.NewIdentity.IdentityData, configschema.SchemaNestedAttributeTypeImpliedType(resSchema.Identity))
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("decode dynamic value for identity data", err)...)
		}
	}

	return resp, diags
}

func (c *Client) PlanResourceChange(ctx context.Context, request typ.PlanResourceChangeRequest) (*typ.PlanResourceChangeResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	schema := c.schemas

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	resTyp, ok := schema.ResourceTypesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	metaTyp := schema.ProviderMetaCty
	capabilities := schema.ServerCapabilities

	var resp typ.PlanResourceChangeResponse

	// If the provider doesn't support planning a destroy operation, we can
	// return immediately.
	if request.ProposedNewState.IsNull() && !capabilities.PlanDestroy {
		resp.PlannedState = request.ProposedNewState
		resp.PlannedPrivate = request.PriorPrivate
		return &resp, nil
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

	if !request.PriorIdentity.IsNull() {
		if resSchema.Identity == nil {
			diags = append(diags, typ.ErrorDiagnostics("identity type not found", fmt.Errorf("identity type not found for resoruce type %s", request.TypeName))...)
			return nil, diags
		}
		priorIdentityMP, err := msgpack.Marshal(request.PriorIdentity, configschema.SchemaNestedAttributeTypeImpliedType(resSchema.Identity))
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("msgpach marshal", err)...)
			return nil, diags
		}
		protoReq.PriorIdentity = &tfprotov6.ResourceIdentityData{
			IdentityData: &tfprotov6.DynamicValue{MsgPack: priorIdentityMP},
		}
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
	resp.PlannedState = state

	for _, p := range protoResp.RequiresReplace {
		resp.RequiresReplace = append(resp.RequiresReplace, convert.DecodeAttributePath(p))
	}

	resp.PlannedPrivate = protoResp.PlannedPrivate

	resp.LegacyTypeSystem = protoResp.UnsafeToUseLegacyTypeSystem

	resp.Deferred = convert.ProtoToDeferred(protoResp.Deferred)

	if protoResp.PlannedIdentity != nil && protoResp.PlannedIdentity.IdentityData != nil {
		if resSchema.Identity == nil {
			diags = append(diags, typ.ErrorDiagnostics("unknown identity type", fmt.Errorf("unknown identity type %s", request.TypeName))...)
			return nil, diags
		}

		resp.PlannedIdentity, err = decodeDynamicValue(protoResp.PlannedIdentity.IdentityData, configschema.SchemaNestedAttributeTypeImpliedType(resSchema.Identity))
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("decode dynamic value for identity data", err)...)
		}
	}

	return &resp, diags
}

func (c *Client) ApplyResourceChange(ctx context.Context, request typ.ApplyResourceChangeRequest) (*typ.ApplyResourceChangeResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	schema := c.schemas

	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

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

	if !request.PlannedIdentity.IsNull() {
		if resSchema.Identity == nil {
			diags = append(diags, typ.ErrorDiagnostics("identity type not found", fmt.Errorf("identity type not found for resoruce type %s", request.TypeName))...)
			return nil, diags
		}
		currentIdentityMP, err := msgpack.Marshal(request.PlannedIdentity, configschema.SchemaNestedAttributeTypeImpliedType(resSchema.Identity))
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("msgpach marshal", err)...)
			return nil, diags
		}
		protoReq.PlannedIdentity = &tfprotov6.ResourceIdentityData{
			IdentityData: &tfprotov6.DynamicValue{MsgPack: currentIdentityMP},
		}
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

	resp := &typ.ApplyResourceChangeResponse{
		NewState:         state,
		Private:          protoResp.Private,
		LegacyTypeSystem: protoResp.UnsafeToUseLegacyTypeSystem,
	}

	if protoResp.NewIdentity != nil && protoResp.NewIdentity.IdentityData != nil {
		if resSchema.Identity == nil {
			diags = append(diags, typ.ErrorDiagnostics("unknown identity type", fmt.Errorf("unknown identity type %s", request.TypeName))...)
			return nil, diags
		}

		resp.NewIdentity, err = decodeDynamicValue(protoResp.NewIdentity.IdentityData, configschema.SchemaNestedAttributeTypeImpliedType(resSchema.Identity))
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("decode dynamic value for identity data", err)...)
		}
	}

	return resp, diags
}

func (c *Client) ImportResourceState(ctx context.Context, request typ.ImportResourceStateRequest) (*typ.ImportResourceStateResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	schema := c.schemas

	protoReq := &tfprotov6.ImportResourceStateRequest{
		TypeName: request.TypeName,
		ID:       request.ID,
		ClientCapabilities: &tfprotov6.ImportResourceStateClientCapabilities{
			DeferralAllowed: request.ClientCapabilities.DeferralAllowed,
		},
	}

	if !request.Identity.IsNull() {
		resSchema, ok := schema.ResourceTypes[request.TypeName]
		if !ok {
			diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
			return nil, diags
		}

		if resSchema.Identity == nil {
			diags = append(diags, typ.ErrorDiagnostics("identity type not found", fmt.Errorf("identity type not found for resoruce type %s", request.TypeName))...)
			return nil, diags
		}
		mp, err := msgpack.Marshal(request.Identity, configschema.SchemaNestedAttributeTypeImpliedType(resSchema.Identity))
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("msgpach marshal", err)...)
			return nil, diags
		}
		protoReq.Identity = &tfprotov6.ResourceIdentityData{
			IdentityData: &tfprotov6.DynamicValue{MsgPack: mp},
		}
	}

	resp, err := c.client.ImportResourceState(ctx, protoReq)
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

		if imported.Identity != nil && imported.Identity.IdentityData != nil {
			importedIdentitySchema, ok := schema.ResourceTypes[imported.TypeName]
			if !ok {
				diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown imported resource type %q", imported.TypeName))...)
				continue
			}

			if importedIdentitySchema.Identity == nil {
				diags = append(diags, typ.ErrorDiagnostics("unknown identity type", fmt.Errorf("unknown identity type %s", imported.TypeName))...)
				continue
			}

			resource.Identity, err = decodeDynamicValue(imported.Identity.IdentityData, configschema.SchemaNestedAttributeTypeImpliedType(importedIdentitySchema.Identity))
			if err != nil {
				diags = append(diags, typ.ErrorDiagnostics("decode dynamic value for identity data", err)...)
				return &response, diags
			}
		}

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

	if len(request.SourceIdentity) > 0 {
		protoReq.SourceIdentity = &tfprotov6.RawState{JSON: request.SourceIdentity}
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

	resp := &typ.MoveResourceStateResponse{
		TargetState:   state,
		TargetPrivate: protoResp.TargetPrivate,
	}

	if protoResp.TargetIdentity != nil && protoResp.TargetIdentity.IdentityData != nil {
		targetResSchema := schema.ResourceTypes[request.TargetTypeName]

		if targetResSchema.Identity == nil {
			diags = append(diags, typ.ErrorDiagnostics("unknown identity type", fmt.Errorf("unknown identity type %s", request.TargetTypeName))...)
			return nil, diags
		}

		resp.TargetIdentity, err = decodeDynamicValue(protoResp.TargetIdentity.IdentityData, configschema.SchemaNestedAttributeTypeImpliedType(targetResSchema.Identity))
		if err != nil {
			diags = append(diags, typ.ErrorDiagnostics("decode dynamic value for identity data", err)...)
		}
	}

	return resp, diags
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

		return resp, diags
	}

	resultVal, err := decodeDynamicValue(protoResp.Result, funcDecl.ReturnType)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("call function error", fmt.Errorf("decoding return value: %v", err))...)
		return nil, diags
	}

	resp.Result = resultVal
	return resp, diags
}

func (c *Client) GetResourceIdentitySchemas(ctx context.Context) (*typ.GetResourceIdentitySchemasResponse, typ.Diagnostics) {
	var diags typ.Diagnostics
	resp := typ.GetResourceIdentitySchemasResponse{
		IdentityTypes: map[string]typ.IdentitySchema{},
	}

	protoResp, err := c.client.GetResourceIdentitySchemas(ctx, &tfprotov6.GetResourceIdentitySchemasRequest{})
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			return &resp, nil
		}

		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	for name, res := range protoResp.IdentitySchemas {
		resp.IdentityTypes[name] = typ.IdentitySchema{
			Version: res.Version,
			Body:    convert.ProtoToIdentitySchema(res.IdentityAttributes),
		}
	}

	return &resp, diags
}

func (c *Client) UpgradeResourceIdentity(ctx context.Context, request typ.UpgradeResourceIdentityRequest) (*typ.UpgradeResourceIdentityResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	schema, _ := c.GetProviderSchema()
	resSchema, ok := schema.ResourceTypes[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	protoReq := &tfprotov6.UpgradeResourceIdentityRequest{
		TypeName: request.TypeName,
		Version:  request.Version,
		RawIdentity: &tfprotov6.RawState{
			JSON: request.RawIdentityJSON,
		},
	}

	protoResp, err := c.client.UpgradeResourceIdentity(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}
	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	ty := configschema.SchemaNestedAttributeTypeImpliedType(resSchema.Identity)

	resp := &typ.UpgradeResourceIdentityResponse{
		UpgradedIdentity: cty.NullVal(ty),
	}
	if protoResp.UpgradedIdentity == nil {
		return resp, diags
	}

	identity, err := decodeDynamicValue(protoResp.UpgradedIdentity.IdentityData, ty)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("decode dynamic value for identity data", err)...)
		return nil, diags
	}

	resp.UpgradedIdentity = identity
	return resp, diags
}

func (c *Client) ValidateEphemeralResourceConfig(ctx context.Context, request typ.ValidateEphemeralResourceConfigRequest) typ.Diagnostics {
	var diags typ.Diagnostics

	schema := c.schemas

	ephemSchema, ok := schema.EphemeralResourceTypesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return diags
	}

	mp, err := msgpack.Marshal(request.Config, ephemSchema)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return diags
	}

	protoReq := &tfprotov6.ValidateEphemeralResourceConfigRequest{
		TypeName: request.TypeName,
		Config:   &tfprotov6.DynamicValue{MsgPack: mp},
	}

	protoResp, err := c.client.ValidateEphemeralResourceConfig(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return diags
	}

	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return diags
	}

	return diags
}

func (c *Client) OpenEphemeralResource(ctx context.Context, request typ.OpenEphemeralResourceRequest) (*typ.OpenEphemeralResourceResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	schema := c.schemas

	ephemSchema, ok := schema.EphemeralResourceTypesCty[request.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics("no schema", fmt.Errorf("unknown resource type %q", request.TypeName))...)
		return nil, diags
	}

	mp, err := msgpack.Marshal(request.Config, ephemSchema)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return nil, diags
	}

	protoReq := &tfprotov6.OpenEphemeralResourceRequest{
		TypeName: request.TypeName,
		Config:   &tfprotov6.DynamicValue{MsgPack: mp},
		ClientCapabilities: &tfprotov6.OpenEphemeralResourceClientCapabilities{
			DeferralAllowed: request.ClientCapabilities.DeferralAllowed,
		},
	}

	protoResp, err := c.client.OpenEphemeralResource(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	state, err := decodeDynamicValue(protoResp.Result, ephemSchema)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("decode dynamic value for result", err)...)
	}

	resp := &typ.OpenEphemeralResourceResponse{
		Result:   state,
		Private:  protoResp.Private,
		Deferred: convert.ProtoToDeferred(protoResp.Deferred),
		RenewAt:  protoResp.RenewAt,
	}

	return resp, diags
}

func (c *Client) RenewEphemeralResource(ctx context.Context, request typ.RenewEphemeralResourceRequest) (*typ.RenewEphemeralResourceResponse, typ.Diagnostics) {
	var diags typ.Diagnostics

	protoReq := &tfprotov6.RenewEphemeralResourceRequest{
		TypeName: request.TypeName,
		Private:  request.Private,
	}

	protoResp, err := c.client.RenewEphemeralResource(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return nil, diags
	}

	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return nil, diags
	}

	resp := &typ.RenewEphemeralResourceResponse{
		Private: protoResp.Private,
		RenewAt: protoResp.RenewAt,
	}

	return resp, diags
}

func (c *Client) CloseEphemeralResource(ctx context.Context, request typ.CloseEphemeralResourceRequest) typ.Diagnostics {
	var diags typ.Diagnostics

	protoReq := &tfprotov6.CloseEphemeralResourceRequest{
		TypeName: request.TypeName,
		Private:  request.Private,
	}

	protoResp, err := c.client.CloseEphemeralResource(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return diags
	}

	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return diags
	}

	return diags
}

func (c *Client) ValidateListResourceConfig(ctx context.Context, req typ.ValidateListResourceConfigRequest) typ.Diagnostics {
	var diags typ.Diagnostics

	schema := c.schemas
	lsch, ok := schema.ListResourceTypes[req.TypeName]
	if !ok {
		return typ.ErrorDiagnostics(fmt.Sprintf(`unknown list resource type "%s"`, req.TypeName), nil)
	}
	if !req.Config.Type().HasAttribute("config") {
		return typ.ErrorDiagnostics(`missing required attribute "config"`, nil)
	}

	configSchema := lsch.Block.NestedBlocks["config"]
	config := req.Config.GetAttr("config")
	mp, err := msgpack.Marshal(config, configschema.SchemaBlockImpliedType(configSchema.Block))
	if err != nil {
		return typ.ErrorDiagnostics("msgpack marshal", err)
	}

	protoReq := &tfprotov6.ValidateListResourceConfigRequest{
		TypeName: req.TypeName,
		Config:   &tfprotov6.DynamicValue{MsgPack: mp},
	}
	protoResp, err := c.client.ValidateListResourceConfig(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return diags
	}

	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return diags
	}

	return diags
}

func (c *Client) ListResource(ctx context.Context, req typ.ListResourceRequest) (resp typ.ListResourceResponse, diags typ.Diagnostics) {
	schema := c.schemas

	listSchema, ok := schema.ListResourceTypes[req.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics(fmt.Sprintf(`unknown list resource type "%s"`, req.TypeName), nil)...)
		return
	}

	resourceSchema, ok := schema.ResourceTypes[req.TypeName]
	if !ok || resourceSchema.Identity == nil {
		diags = append(diags, typ.ErrorDiagnostics(fmt.Sprintf(`identitiy schema not found for resource type "%s"`, req.TypeName), nil)...)
		return
	}
	resourceSchemaCty := schema.ResourceTypesCty[req.TypeName]

	if !req.Config.Type().HasAttribute("config") {
		diags = append(diags, typ.ErrorDiagnostics(`missing required attribute "config"`, nil)...)
		return
	}

	config := req.Config.GetAttr("config")
	mp, err := msgpack.Marshal(config, configschema.SchemaBlockImpliedType(listSchema.Block))
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
	}

	protoReq := &tfprotov6.ListResourceRequest{
		TypeName:        req.TypeName,
		Config:          &tfprotov6.DynamicValue{MsgPack: mp},
		IncludeResource: req.IncludeResourceObject,
		Limit:           req.Limit,
	}

	// Start the streaming RPC with a context. The context will be cancelled
	// when this function returns, which will stop the stream if it is still
	// running.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	stream, err := c.client.ListResource(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return
	}

	resp.Result = cty.DynamicVal
	values := make([]cty.Value, 0)

	// Process the stream
	for event := range stream.Results {
		if int64(len(values)) >= req.Limit {
			// If we have reached the limit, we stop receiving events
			break
		}

		if slices.ContainsFunc(event.Diagnostics, func(diag *tfprotov6.Diagnostic) bool {
			return diag != nil && diag.Severity == tfprotov6.DiagnosticSeverityError
		}) {
			// If we have errors, we stop processing and return early
			break
		}

		if slices.ContainsFunc(event.Diagnostics, func(diag *tfprotov6.Diagnostic) bool {
			return diag != nil && diag.Severity == tfprotov6.DiagnosticSeverityWarning
		}) && event.Identity.IdentityData == nil {
			// If we have warnings but no identity data, we continue with the next event
			break
		}

		obj := map[string]cty.Value{
			"display_name": cty.StringVal(event.DisplayName),
			"state":        cty.NullVal(resourceSchemaCty),
			"identity":     cty.NullVal(configschema.SchemaNestedAttributeTypeImpliedType(resourceSchema.Identity)),
		}

		// Handle identity data - it must be present
		if event.Identity == nil || event.Identity.IdentityData == nil {
			diags = append(diags, typ.ErrorDiagnostics(fmt.Sprintf("missing identity data in ListResource event for %s", req.TypeName), nil)...)
		} else {
			identityVal, err := decodeDynamicValue(event.Identity.IdentityData, configschema.SchemaNestedAttributeTypeImpliedType(resourceSchema.Identity))
			if err != nil {
				diags = append(diags, typ.ErrorDiagnostics(err.Error(), nil)...)
			} else {
				obj["identity"] = identityVal
			}
		}

		// Handle resource object if present and requested
		if event.Resource != nil && req.IncludeResourceObject {
			// Use the ResourceTypes schema for the resource object
			resourceObj, err := decodeDynamicValue(event.Resource, resourceSchemaCty)
			if err != nil {
				diags = append(diags, typ.ErrorDiagnostics(err.Error(), nil)...)
			} else {
				obj["state"] = resourceObj
			}
		}

		if diags.HasErrors() {
			// If validation errors occurred, we stop processing and return early
			break
		}

		values = append(values, cty.ObjectVal(obj))
	}

	// The provider result of a list resource is always a list, but
	// we will wrap that list in an object with a single attribute "data",
	// so that we can differentiate between a list resource instance (list.aws_instance.test[index])
	// and the elements of the result of a list resource instance (list.aws_instance.test.data[index])
	resp.Result = cty.ObjectVal(map[string]cty.Value{
		"data":   cty.TupleVal(values),
		"config": config,
	})
	return resp, diags
}

func (c *Client) ValidateActionConfig(ctx context.Context, req typ.ValidateActionConfigRequest) (diags typ.Diagnostics) {
	schema := c.schemas

	actionSchema, ok := schema.ActionsCty[req.TypeName]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics(fmt.Sprintf(`unknown action type "%s"`, req.TypeName), nil)...)
		return
	}

	mp, err := msgpack.Marshal(req.Config, actionSchema)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return
	}

	protoReq := &tfprotov6.ValidateActionConfigRequest{
		ActionType: req.TypeName,
		Config:     &tfprotov6.DynamicValue{MsgPack: mp},
	}

	protoResp, err := c.client.ValidateActionConfig(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return
	}

	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return diags
	}

	return diags
}

func (c *Client) PlanAction(ctx context.Context, req typ.PlanActionRequest) (diags typ.Diagnostics, resp typ.PlanActionResponse) {
	schema := c.schemas

	actionSchema, ok := schema.ActionsCty[req.ActionType]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics(fmt.Sprintf(`unknown action type "%s"`, req.ActionType), nil)...)
		return
	}

	mp, err := msgpack.Marshal(req.ProposedActionData, actionSchema)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return
	}

	protoReq := &tfprotov6.PlanActionRequest{
		ActionType: req.ActionType,
		Config:     &tfprotov6.DynamicValue{MsgPack: mp},
		ClientCapabilities: &tfprotov6.PlanActionClientCapabilities{
			DeferralAllowed: req.ClientCapabilities.DeferralAllowed,
		},
	}

	protoResp, err := c.client.PlanAction(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return
	}

	respDiags := convert.DecodeDiagnostics(protoResp.Diagnostics)
	diags = append(diags, respDiags...)
	if diags.HasErrors() {
		return
	}

	if def := protoResp.Deferred; def != nil {
		resp.Deferred = &typ.Deferred{
			Reason: typ.DeferredReason(def.Reason.String()),
		}
	}

	return
}

func (c *Client) InvokeAction(ctx context.Context, req typ.InvokeActionRequest) (diags typ.Diagnostics, resp typ.InvokeActionResponse) {
	schema := c.schemas

	actionSchema, ok := schema.ActionsCty[req.ActionType]
	if !ok {
		diags = append(diags, typ.ErrorDiagnostics(fmt.Sprintf(`unknown action type "%s"`, req.ActionType), nil)...)
		return
	}

	mp, err := msgpack.Marshal(req.PlannedActionData, actionSchema)
	if err != nil {
		diags = append(diags, typ.ErrorDiagnostics("msgpack marshal", err)...)
		return
	}

	protoReq := &tfprotov6.InvokeActionRequest{
		ActionType:         req.ActionType,
		Config:             &tfprotov6.DynamicValue{MsgPack: mp},
		ClientCapabilities: &tfprotov6.InvokeActionClientCapabilities{},
	}

	stream, err := c.client.InvokeAction(ctx, protoReq)
	if err != nil {
		diags = append(diags, typ.RPCErrorDiagnostics(err)...)
		return
	}

	resp.Events = func(yield func(typ.InvokeActionEvent) bool) {
		for evt := range stream.Events {
			switch ev := evt.Type.(type) {
			case *tfprotov6.ProgressInvokeActionEventType:
				yield(typ.InvokeActionEvent_Progress{
					Message: ev.Message,
				})
			case *tfprotov6.CompletedInvokeActionEventType:
				yield(typ.InvokeActionEvent_Completed{
					Diagnostics: convert.DecodeDiagnostics(ev.Diagnostics),
				})
			default:
				panic(fmt.Sprintf("unexpected event type %T in InvokeAction response", evt.Type))
			}
		}
	}

	return
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
