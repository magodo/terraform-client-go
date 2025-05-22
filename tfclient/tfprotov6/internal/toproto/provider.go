package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func GetMetadata_Request(in *tfprotov6.GetMetadataRequest) *tfplugin6.GetMetadata_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin6.GetMetadata_Request{}

	return req
}

func GetProviderSchema_Request(in *tfprotov6.GetProviderSchemaRequest) *tfplugin6.GetProviderSchema_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin6.GetProviderSchema_Request{}

	return req
}

func GetResourceIdentitySchemas_Request(in *tfprotov6.GetResourceIdentitySchemasRequest) *tfplugin6.GetResourceIdentitySchemas_Request {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.GetResourceIdentitySchemas_Request{}

	return resp
}

func ValidateProviderConfig_Request(in *tfprotov6.ValidateProviderConfigRequest) *tfplugin6.ValidateProviderConfig_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin6.ValidateProviderConfig_Request{
		Config: DynamicValue(in.Config),
	}

	return req
}

func ConfigureProvider_Request(in *tfprotov6.ConfigureProviderRequest) *tfplugin6.ConfigureProvider_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin6.ConfigureProvider_Request{
		TerraformVersion: in.TerraformVersion,
		Config:           DynamicValue(in.Config),
	}

	if in.ClientCapabilities != nil {
		req.ClientCapabilities = &tfplugin6.ClientCapabilities{
			DeferralAllowed: in.ClientCapabilities.DeferralAllowed,
		}
	}

	return req
}

func StopProvider_Request(in *tfprotov6.StopProviderRequest) *tfplugin6.StopProvider_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin6.StopProvider_Request{}

	return req
}
