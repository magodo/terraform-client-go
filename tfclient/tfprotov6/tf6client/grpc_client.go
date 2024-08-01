package tf6client

import (
	"context"

	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/fromproto"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/toproto"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type GRPCClient struct {
	client tfplugin6.ProviderClient
}

var _ tfprotov6.ProviderServer = &GRPCClient{}

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
