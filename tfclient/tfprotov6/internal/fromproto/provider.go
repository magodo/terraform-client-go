package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func GetMetadataRequest(in *tfplugin6.GetMetadata_Request) *tfprotov6.GetMetadataRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.GetMetadataRequest{}

	return resp
}

func GetMetadataResponse(in *tfplugin6.GetMetadata_Response) (*tfprotov6.GetMetadataResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.GetMetadataResponse{
		ServerCapabilities: ServerCapabilities(in.ServerCapabilities),
		Diagnostics:        diags,
	}

	for _, datasource := range in.DataSources {
		if v := DataSourceMetadata(datasource); v != nil {
			resp.DataSources = append(resp.DataSources, *v)
		}
	}

	for _, resource := range in.Resources {
		if v := ResourceMetadata(resource); v != nil {
			resp.Resources = append(resp.Resources, *v)
		}
	}

	for _, f := range in.Functions {
		if v := FunctionMetadata(f); v != nil {
			resp.Functions = append(resp.Functions, *v)
		}
	}

	return resp, nil
}

func GetProviderSchemaRequest(in *tfplugin6.GetProviderSchema_Request) *tfprotov6.GetProviderSchemaRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.GetProviderSchemaRequest{}

	return resp
}

func GetProviderSchemaResponse(in *tfplugin6.GetProviderSchema_Response) (*tfprotov6.GetProviderSchemaResponse, error) {
	provider, err := Schema(in.Provider)
	if err != nil {
		return nil, err
	}

	providerMeta, err := Schema(in.ProviderMeta)
	if err != nil {
		return nil, err
	}

	resourceSchemas := make(map[string]*tfprotov6.Schema, len(in.ResourceSchemas))
	for k, v := range in.ResourceSchemas {
		schema, err := Schema(v)
		if err != nil {
			return nil, err
		}
		resourceSchemas[k] = schema
	}

	dataSourceSchemas := make(map[string]*tfprotov6.Schema, len(in.DataSourceSchemas))
	for k, v := range in.DataSourceSchemas {
		schema, err := Schema(v)
		if err != nil {
			return nil, err
		}
		dataSourceSchemas[k] = schema
	}

	funcs, err := Functions(in.Functions)
	if err != nil {
		return nil, err
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.GetProviderSchemaResponse{
		ServerCapabilities: ServerCapabilities(in.ServerCapabilities),
		Provider:           provider,
		ProviderMeta:       providerMeta,
		ResourceSchemas:    resourceSchemas,
		DataSourceSchemas:  dataSourceSchemas,
		Functions:          funcs,
		Diagnostics:        diags,
	}

	return resp, nil
}

func ValidateProviderConfigRequest(in *tfplugin6.ValidateProviderConfig_Request) *tfprotov6.ValidateProviderConfigRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ValidateProviderConfigRequest{
		Config: DynamicValue(in.Config),
	}

	return resp
}

func ValidateProviderConfigResponse(in *tfplugin6.ValidateProviderConfig_Response) (*tfprotov6.ValidateProviderConfigResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ValidateProviderConfigResponse{
		PreparedConfig: nil, // TODO: in has no that field
		Diagnostics:    diags,
	}, nil
}

func ConfigureProviderRequest(in *tfplugin6.ConfigureProvider_Request) *tfprotov6.ConfigureProviderRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ConfigureProviderRequest{
		Config:             DynamicValue(in.Config),
		TerraformVersion:   in.TerraformVersion,
		ClientCapabilities: ConfigureProviderClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func ConfigureProviderResponse(in *tfplugin6.ConfigureProvider_Response) (*tfprotov6.ConfigureProviderResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ConfigureProviderResponse{
		Diagnostics: diags,
	}, nil
}

func StopProviderRequest(in *tfplugin6.StopProvider_Request) *tfprotov6.StopProviderRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.StopProviderRequest{}

	return resp
}

func StopProviderResponse(in *tfplugin6.StopProvider_Response) (*tfprotov6.StopProviderResponse, error) {
	return &tfprotov6.StopProviderResponse{
		Error: in.Error,
	}, nil
}
