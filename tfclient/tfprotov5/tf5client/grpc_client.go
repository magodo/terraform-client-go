package tf5client

import (
	"context"
	"io"

	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/fromproto"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/toproto"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

type GRPCClient struct {
	client tfplugin5.ProviderClient
}

var _ tfprotov5.ProviderServer = &GRPCClient{}
var _ tfprotov5.ActionServer = &GRPCClient{}
var _ tfprotov5.ListResourceServer = &GRPCClient{}

// ApplyResourceChange implements tfprotov5.ProviderServer
func (c *GRPCClient) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	r := toproto.ApplyResourceChange_Request(req)
	resp, err := c.client.ApplyResourceChange(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ApplyResourceChangeResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ImportResourceState implements tfprotov5.ProviderServer
func (c *GRPCClient) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	r := toproto.ImportResourceState_Request(req)
	resp, err := c.client.ImportResourceState(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ImportResourceStateResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// PlanResourceChange implements tfprotov5.ProviderServer
func (c *GRPCClient) PlanResourceChange(ctx context.Context, req *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	r := toproto.PlanResourceChange_Request(req)
	resp, err := c.client.PlanResourceChange(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.PlanResourceChangeResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ReadResource implements tfprotov5.ProviderServer
func (c *GRPCClient) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	r := toproto.ReadResource_Request(req)
	resp, err := c.client.ReadResource(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ReadResourceResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// UpgradeResourceState implements tfprotov5.ProviderServer
func (c *GRPCClient) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	r := toproto.UpgradeResourceState_Request(req)
	resp, err := c.client.UpgradeResourceState(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.UpgradeResourceStateResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ValidateResourceTypeConfig implements tfprotov5.ProviderServer
func (c *GRPCClient) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	r := toproto.ValidateResourceTypeConfig_Request(req)
	resp, err := c.client.ValidateResourceTypeConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ValidateResourceTypeConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ReadDataSource implements tfprotov5.ProviderServer
func (c *GRPCClient) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	r := toproto.ReadDataSource_Request(req)
	resp, err := c.client.ReadDataSource(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ReadDataSourceResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ValidateDataSourceConfig implements tfprotov5.ProviderServer
func (c *GRPCClient) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	r := toproto.ValidateDataSourceConfig_Request(req)
	resp, err := c.client.ValidateDataSourceConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ValidateDataSourceConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ConfigureProvider implements tfprotov5.ProviderServer
func (c *GRPCClient) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	r := toproto.Configure_Request(req)
	resp, err := c.client.Configure(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ConfigureProviderResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetProviderSchema implements tfprotov5.ProviderServer
func (c *GRPCClient) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	r := toproto.GetProviderSchema_Request(req)
	resp, err := c.client.GetSchema(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.GetProviderSchemaResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// PrepareProviderConfig implements tfprotov5.ProviderServer
func (c *GRPCClient) PrepareProviderConfig(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	r := toproto.PrepareProviderConfig_Request(req)
	resp, err := c.client.PrepareProviderConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.PrepareProviderConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// StopProvider implements tfprotov5.ProviderServer
func (c *GRPCClient) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	r := toproto.Stop_Request(req)
	resp, err := c.client.Stop(ctx, r)
	if err != nil {
		return nil, err
	}
	ret := fromproto.StopProviderResponse(resp)
	return ret, nil
}

func (c *GRPCClient) CallFunction(ctx context.Context, req *tfprotov5.CallFunctionRequest) (*tfprotov5.CallFunctionResponse, error) {
	r := toproto.CallFunction_Request(req)
	resp, err := c.client.CallFunction(ctx, r)
	if err != nil {
		return nil, err
	}
	ret := fromproto.CallFunctionResponse(resp)
	return ret, nil
}

func (c *GRPCClient) GetFunctions(ctx context.Context, req *tfprotov5.GetFunctionsRequest) (*tfprotov5.GetFunctionsResponse, error) {
	r := toproto.GetFunctions_Request(req)
	resp, err := c.client.GetFunctions(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.GetFunctionsResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) GetMetadata(ctx context.Context, req *tfprotov5.GetMetadataRequest) (*tfprotov5.GetMetadataResponse, error) {
	r := toproto.GetMetadata_Request(req)
	resp, err := c.client.GetMetadata(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.GetMetadataResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) MoveResourceState(ctx context.Context, req *tfprotov5.MoveResourceStateRequest) (*tfprotov5.MoveResourceStateResponse, error) {
	r := toproto.MoveResourceState_Request(req)
	resp, err := c.client.MoveResourceState(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.MoveResourceStateResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) GetResourceIdentitySchemas(ctx context.Context, req *tfprotov5.GetResourceIdentitySchemasRequest) (*tfprotov5.GetResourceIdentitySchemasResponse, error) {
	r := toproto.GetResourceIdentitySchemas_Request(req)
	resp, err := c.client.GetResourceIdentitySchemas(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.GetResourceIdentitySchemasResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) UpgradeResourceIdentity(ctx context.Context, req *tfprotov5.UpgradeResourceIdentityRequest) (*tfprotov5.UpgradeResourceIdentityResponse, error) {
	r := toproto.UpgradeResourceIdentity_Request(req)
	resp, err := c.client.UpgradeResourceIdentity(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.UpgradeResourceIdentityResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov5.ValidateEphemeralResourceConfigRequest) (*tfprotov5.ValidateEphemeralResourceConfigResponse, error) {
	r := toproto.ValidateEphemeralResourceConfigRequest(req)
	resp, err := c.client.ValidateEphemeralResourceConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ValidateEphemeralResourceConfig_Response(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) OpenEphemeralResource(ctx context.Context, req *tfprotov5.OpenEphemeralResourceRequest) (*tfprotov5.OpenEphemeralResourceResponse, error) {
	r := toproto.OpenEphemeralResourceRequest(req)
	resp, err := c.client.OpenEphemeralResource(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.OpenEphemeralResource_Response(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) RenewEphemeralResource(ctx context.Context, req *tfprotov5.RenewEphemeralResourceRequest) (*tfprotov5.RenewEphemeralResourceResponse, error) {
	r := toproto.RenewEphemeralResourceRequest(req)
	resp, err := c.client.RenewEphemeralResource(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.RenewEphemeralResource_Response(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) CloseEphemeralResource(ctx context.Context, req *tfprotov5.CloseEphemeralResourceRequest) (*tfprotov5.CloseEphemeralResourceResponse, error) {
	r := toproto.CloseEphemeralResourceRequest(req)
	resp, err := c.client.CloseEphemeralResource(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.CloseEphemeralResource_Response(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) ValidateListResourceConfig(ctx context.Context, req *tfprotov5.ValidateListResourceConfigRequest) (*tfprotov5.ValidateListResourceConfigResponse, error) {
	r := toproto.ValidateListResourceConfigRequest(req)
	resp, err := c.client.ValidateListResourceConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	return fromproto.ValidateListResourceConfig_Response(resp)
}

func (c *GRPCClient) ListResource(ctx context.Context, req *tfprotov5.ListResourceRequest) (*tfprotov5.ListResourceServerStream, error) {
	r := toproto.ListResourceRequest(req)
	resp, err := c.client.ListResource(ctx, r)
	if err != nil {
		return nil, err
	}

	var result tfprotov5.ListResourceServerStream
	result.Results = func(yield func(tfprotov5.ListResourceResult) bool) {
		for {
			event, err := resp.Recv()
			if err != nil {
				break
			}
			evt, err := fromproto.ListResource_ListResourceEvent(event)
			if err != nil {
				// TODO: need a better error handling
				continue
			}
			if !yield(*evt) {
				break
			}
		}
	}
	return &result, nil
}

func (c *GRPCClient) ValidateActionConfig(ctx context.Context, req *tfprotov5.ValidateActionConfigRequest) (*tfprotov5.ValidateActionConfigResponse, error) {
	r := toproto.ValidateActionConfigRequest(req)
	resp, err := c.client.ValidateActionConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	return fromproto.ValidateActionConfig_Response(resp)
}

func (c *GRPCClient) PlanAction(ctx context.Context, req *tfprotov5.PlanActionRequest) (*tfprotov5.PlanActionResponse, error) {
	r := toproto.PlanActionRequest(req)
	resp, err := c.client.PlanAction(ctx, r)
	if err != nil {
		return nil, err
	}
	return fromproto.PlanAction_Response(resp)
}

func (c *GRPCClient) InvokeAction(ctx context.Context, req *tfprotov5.InvokeActionRequest) (*tfprotov5.InvokeActionServerStream, error) {
	r := toproto.InvokeActionRequest(req)
	resp, err := c.client.InvokeAction(ctx, r)
	if err != nil {
		return nil, err
	}
	var result tfprotov5.InvokeActionServerStream
	result.Events = func(yield func(tfprotov5.InvokeActionEvent) bool) {
		for {
			event, err := resp.Recv()

			var evt *tfprotov5.InvokeActionEvent
			// This follows the same logic as terraform/internal/plugin/grpc_provider.go does.
			if err == io.EOF {
				break
			}
			if err != nil {
				// We handle this by returning a finished response with the error
				// If the client errors we won't be receiving any more events.
				evt = &tfprotov5.InvokeActionEvent{
					Type: tfprotov5.CompletedInvokeActionEventType{
						Diagnostics: []*tfprotov5.Diagnostic{
							{
								Severity: tfprotov5.DiagnosticSeverityError,
								Summary:  "rpc error",
								Detail:   err.Error(),
							},
						},
					},
				}
			} else {
				evt, err = fromproto.InvokeAction_InvokeActionEvent(event)
				if err != nil {
					// TODO: need a better error handling
					continue
				}
			}
			if !yield(*evt) {
				break
			}
		}
	}
	return &result, nil
}
