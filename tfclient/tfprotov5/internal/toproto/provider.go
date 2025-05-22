package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func GetMetadata_Request(in *tfprotov5.GetMetadataRequest) *tfplugin5.GetMetadata_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.GetMetadata_Request{}

	return req
}

func GetProviderSchema_Request(in *tfprotov5.GetProviderSchemaRequest) *tfplugin5.GetProviderSchema_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.GetProviderSchema_Request{}

	return req
}

func GetResourceIdentitySchemas_Request(in *tfprotov5.GetResourceIdentitySchemasRequest) *tfplugin5.GetResourceIdentitySchemas_Request {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetResourceIdentitySchemas_Request{}

	return resp
}

func PrepareProviderConfig_Request(in *tfprotov5.PrepareProviderConfigRequest) *tfplugin5.PrepareProviderConfig_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.PrepareProviderConfig_Request{
		Config: DynamicValue(in.Config),
	}

	return req
}

func Configure_Request(in *tfprotov5.ConfigureProviderRequest) *tfplugin5.Configure_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.Configure_Request{
		TerraformVersion: in.TerraformVersion,
		Config:           DynamicValue(in.Config),
	}

	if in.ClientCapabilities != nil {
		req.ClientCapabilities = &tfplugin5.ClientCapabilities{
			DeferralAllowed: in.ClientCapabilities.DeferralAllowed,
		}
	}

	return req
}

func Stop_Request(in *tfprotov5.StopProviderRequest) *tfplugin5.Stop_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.Stop_Request{}

	return req
}
