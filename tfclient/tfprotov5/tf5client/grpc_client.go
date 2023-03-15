package tf5client

import (
	"context"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/fromproto"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/toproto"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

type GRPCClient struct {
	client tfplugin5.ProviderClient
}

var _ tfprotov5.ProviderServer = &GRPCClient{}

// ApplyResourceChange implements tfprotov5.ProviderServer
func (c *GRPCClient) ApplyResourceChange(ctx context.Context, req *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	r, err := toproto.ApplyResourceChange_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.ImportResourceState_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.PlanResourceChange_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.ReadResource_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.UpgradeResourceState_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.ValidateResourceTypeConfig_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.ReadDataSource_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.ValidateDataSourceConfig_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.Configure_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.GetProviderSchema_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.PrepareProviderConfig_Request(req)
	if err != nil {
		return nil, err
	}
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
	r, err := toproto.Stop_Request(req)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Stop(ctx, r)
	if err != nil {
		return nil, err
	}
	ret, err := fromproto.StopProviderResponse(resp)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
