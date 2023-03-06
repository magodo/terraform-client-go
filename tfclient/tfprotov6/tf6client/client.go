package tf6client

import (
	"context"
	fromproto2 "github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/fromproto"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
	toproto2 "github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/toproto"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type GRPCClient struct {
	client tfplugin6.ProviderClient
}

var _ tfprotov6.ProviderServer = &GRPCClient{}

// ApplyResourceChange implements tfprotov6.ProviderServer
func (c *GRPCClient) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	r, err := toproto2.ApplyResourceChange_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ApplyResourceChange(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.ApplyResourceChangeResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ImportResourceState implements tfprotov6.ProviderServer
func (c *GRPCClient) ImportResourceState(ctx context.Context, req *tfprotov6.ImportResourceStateRequest) (*tfprotov6.ImportResourceStateResponse, error) {
	r, err := toproto2.ImportResourceState_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ImportResourceState(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.ImportResourceStateResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// PlanResourceChange implements tfprotov6.ProviderServer
func (c *GRPCClient) PlanResourceChange(ctx context.Context, req *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	r, err := toproto2.PlanResourceChange_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.PlanResourceChange(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.PlanResourceChangeResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ReadResource implements tfprotov6.ProviderServer
func (c *GRPCClient) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	r, err := toproto2.ReadResource_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ReadResource(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.ReadResourceResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// UpgradeResourceState implements tfprotov6.ProviderServer
func (c *GRPCClient) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	r, err := toproto2.UpgradeResourceState_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.UpgradeResourceState(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.UpgradeResourceStateResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ValidateResourceConfig implements tfprotov6.ProviderServer
func (c *GRPCClient) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	r, err := toproto2.ValidateResourceConfig_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ValidateResourceConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.ValidateResourceConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ReadDataSource implements tfprotov6.ProviderServer
func (c *GRPCClient) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	r, err := toproto2.ReadDataSource_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ReadDataSource(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.ReadDataSourceResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ValidateDataResourceConfig implements tfprotov6.ProviderServer
func (c *GRPCClient) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	r, err := toproto2.ValidateDataResourceConfig_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ValidateDataResourceConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.ValidateDataResourceConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ConfigureProvider implements tfprotov6.ProviderServer
func (c *GRPCClient) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	r, err := toproto2.Configure_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ConfigureProvider(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.ConfigureProviderResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// GetProviderSchema implements tfprotov6.ProviderServer
func (c *GRPCClient) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	r, err := toproto2.GetProviderSchema_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.GetProviderSchema(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.GetProviderSchemaResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// StopProvider implements tfprotov6.ProviderServer
func (c *GRPCClient) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	r, err := toproto2.Stop_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.StopProvider(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.StopProviderResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// ValidateProviderConfig implements tfprotov6.ProviderServer
func (c *GRPCClient) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	r, err := toproto2.ValidateProviderConfig_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.ValidateProviderConfig(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto2.ValidateProviderConfigResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
