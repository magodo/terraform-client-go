package tf6client

import (
	"context"
	"io"

	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/fromproto"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/toproto"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type GRPCClient struct {
	client tfplugin6.ProviderClient
}

var _ tfprotov6.ProviderServer = &GRPCClient{}
var _ tfprotov6.ListResourceServer = &GRPCClient{}
var _ tfprotov6.ActionServer = &GRPCClient{}

// ApplyResourceChange implements tfprotov6.ProviderServer
func (c *GRPCClient) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
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

// ImportResourceState implements tfprotov6.ProviderServer
func (c *GRPCClient) ImportResourceState(ctx context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
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

// PlanResourceChange implements tfprotov6.ProviderServer
func (c *GRPCClient) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
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

// ReadResource implements tfprotov6.ProviderServer
func (c *GRPCClient) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
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

// UpgradeResourceState implements tfprotov6.ProviderServer
func (c *GRPCClient) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
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

// ValidateResourceConfig implements tfprotov6.ProviderServer
func (c *GRPCClient) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	r := toproto.ValidateResourceConfig_Request(req)
	resp, err := c.client.ValidateResourceConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ValidateResourceConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ReadDataSource implements tfprotov6.ProviderServer
func (c *GRPCClient) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
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

// ValidateDataResourceConfig implements tfprotov6.ProviderServer
func (c *GRPCClient) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	r := toproto.ValidateDataResourceConfig_Request(req)
	resp, err := c.client.ValidateDataResourceConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ValidateDataResourceConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ConfigureProvider implements tfprotov6.ProviderServer
func (c *GRPCClient) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	r := toproto.ConfigureProvider_Request(req)
	resp, err := c.client.ConfigureProvider(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ConfigureProviderResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetProviderSchema implements tfprotov6.ProviderServer
func (c *GRPCClient) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	r := toproto.GetProviderSchema_Request(req)
	resp, err := c.client.GetProviderSchema(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.GetProviderSchemaResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// StopProvider implements tfprotov6.ProviderServer
func (c *GRPCClient) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	r := toproto.StopProvider_Request(req)
	resp, err := c.client.StopProvider(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.StopProviderResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ValidateProviderConfig implements tfprotov6.ProviderServer
func (c *GRPCClient) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	r := toproto.ValidateProviderConfig_Request(req)
	resp, err := c.client.ValidateProviderConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.ValidateProviderConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *GRPCClient) CallFunction(ctx context.Context, req *tfprotov6.CallFunctionRequest) (*tfprotov6.CallFunctionResponse, error) {
	r := toproto.CallFunction_Request(req)
	resp, err := c.client.CallFunction(ctx, r)
	if err != nil {
		return nil, err
	}
	ret := fromproto.CallFunctionResponse(resp)
	return ret, nil
}

func (c *GRPCClient) GetFunctions(ctx context.Context, req *tfprotov6.GetFunctionsRequest) (*tfprotov6.GetFunctionsResponse, error) {
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

func (c *GRPCClient) GetMetadata(ctx context.Context, req *tfprotov6.GetMetadataRequest) (*tfprotov6.GetMetadataResponse, error) {
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

func (c *GRPCClient) MoveResourceState(ctx context.Context, req *tfprotov6.MoveResourceStateRequest) (*tfprotov6.MoveResourceStateResponse, error) {
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

func (c *GRPCClient) GetResourceIdentitySchemas(ctx context.Context, req *tfprotov6.GetResourceIdentitySchemasRequest) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
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

func (c *GRPCClient) UpgradeResourceIdentity(ctx context.Context, req *tfprotov6.UpgradeResourceIdentityRequest) (*tfprotov6.UpgradeResourceIdentityResponse, error) {
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

func (c *GRPCClient) ValidateEphemeralResourceConfig(ctx context.Context, req *tfprotov6.ValidateEphemeralResourceConfigRequest) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
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

func (c *GRPCClient) OpenEphemeralResource(ctx context.Context, req *tfprotov6.OpenEphemeralResourceRequest) (*tfprotov6.OpenEphemeralResourceResponse, error) {
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

func (c *GRPCClient) RenewEphemeralResource(ctx context.Context, req *tfprotov6.RenewEphemeralResourceRequest) (*tfprotov6.RenewEphemeralResourceResponse, error) {
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

func (c *GRPCClient) CloseEphemeralResource(ctx context.Context, req *tfprotov6.CloseEphemeralResourceRequest) (*tfprotov6.CloseEphemeralResourceResponse, error) {
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

func (c *GRPCClient) ValidateListResourceConfig(ctx context.Context, req *tfprotov6.ValidateListResourceConfigRequest) (*tfprotov6.ValidateListResourceConfigResponse, error) {
	r := toproto.ValidateListResourceConfigRequest(req)
	resp, err := c.client.ValidateListResourceConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	return fromproto.ValidateListResourceConfig_Response(resp)
}

func (c *GRPCClient) ListResource(ctx context.Context, req *tfprotov6.ListResourceRequest) (*tfprotov6.ListResourceServerStream, error) {
	r := toproto.ListResourceRequest(req)
	resp, err := c.client.ListResource(ctx, r)
	if err != nil {
		return nil, err
	}

	var result tfprotov6.ListResourceServerStream
	result.Results = func(yield func(tfprotov6.ListResourceResult) bool) {
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

func (c *GRPCClient) ValidateActionConfig(ctx context.Context, req *tfprotov6.ValidateActionConfigRequest) (*tfprotov6.ValidateActionConfigResponse, error) {
	r := toproto.ValidateActionConfigRequest(req)
	resp, err := c.client.ValidateActionConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	return fromproto.ValidateActionConfig_Response(resp)
}

func (c *GRPCClient) PlanAction(ctx context.Context, req *tfprotov6.PlanActionRequest) (*tfprotov6.PlanActionResponse, error) {
	r := toproto.PlanActionRequest(req)
	resp, err := c.client.PlanAction(ctx, r)
	if err != nil {
		return nil, err
	}
	return fromproto.PlanAction_Response(resp)
}

func (c *GRPCClient) InvokeAction(ctx context.Context, req *tfprotov6.InvokeActionRequest) (*tfprotov6.InvokeActionServerStream, error) {
	r := toproto.InvokeActionRequest(req)
	resp, err := c.client.InvokeAction(ctx, r)
	if err != nil {
		return nil, err
	}
	var result tfprotov6.InvokeActionServerStream
	result.Events = func(yield func(tfprotov6.InvokeActionEvent) bool) {
		for {
			event, err := resp.Recv()

			var evt *tfprotov6.InvokeActionEvent
			// This follows the same logic as terraform/internal/plugin/grpc_provider.go does.
			if err == io.EOF {
				break
			}
			if err != nil {
				// We handle this by returning a finished response with the error
				// If the client errors we won't be receiving any more events.
				evt = &tfprotov6.InvokeActionEvent{
					Type: tfprotov6.CompletedInvokeActionEventType{
						Diagnostics: []*tfprotov6.Diagnostic{
							{
								Severity: tfprotov6.DiagnosticSeverityError,
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
