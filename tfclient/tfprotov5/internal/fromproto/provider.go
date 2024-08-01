package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func GetMetadataRequest(in *tfplugin5.GetMetadata_Request) *tfprotov5.GetMetadataRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.GetMetadataRequest{}

	return resp
}

func GetMetadataResponse(in *tfplugin5.GetMetadata_Response) (*tfprotov5.GetMetadataResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.GetMetadataResponse{
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

func GetProviderSchemaRequest(in *tfplugin5.GetProviderSchema_Request) *tfprotov5.GetProviderSchemaRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.GetProviderSchemaRequest{}

	return resp
}

func GetProviderSchemaResponse(in *tfplugin5.GetProviderSchema_Response) (*tfprotov5.GetProviderSchemaResponse, error) {
	provider, err := Schema(in.Provider)
	if err != nil {
		return nil, err
	}

	providerMeta, err := Schema(in.ProviderMeta)
	if err != nil {
		return nil, err
	}

	resourceSchemas := make(map[string]*tfprotov5.Schema, len(in.ResourceSchemas))
	for k, v := range in.ResourceSchemas {
		schema, err := Schema(v)
		if err != nil {
			return nil, err
		}
		resourceSchemas[k] = schema
	}

	dataSourceSchemas := make(map[string]*tfprotov5.Schema, len(in.DataSourceSchemas))
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

	resp := &tfprotov5.GetProviderSchemaResponse{
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

func PrepareProviderConfigRequest(in *tfplugin5.PrepareProviderConfig_Request) *tfprotov5.PrepareProviderConfigRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.PrepareProviderConfigRequest{
		Config: DynamicValue(in.Config),
	}

	return resp
}

func PrepareProviderConfigResponse(in *tfplugin5.PrepareProviderConfig_Response) (*tfprotov5.PrepareProviderConfigResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.PrepareProviderConfigResponse{
		PreparedConfig: DynamicValue(in.PreparedConfig),
		Diagnostics:    diags,
	}

	return resp, nil
}

func ConfigureProviderRequest(in *tfplugin5.Configure_Request) *tfprotov5.ConfigureProviderRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.ConfigureProviderRequest{
		Config:             DynamicValue(in.Config),
		TerraformVersion:   in.TerraformVersion,
		ClientCapabilities: ConfigureProviderClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func ConfigureProviderResponse(in *tfplugin5.Configure_Response) (*tfprotov5.ConfigureProviderResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	return &tfprotov5.ConfigureProviderResponse{
		Diagnostics: diags,
	}, nil
}

func StopProviderRequest(in *tfplugin5.Stop_Request) *tfprotov5.StopProviderRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.StopProviderRequest{}

	return resp
}

func StopProviderResponse(in *tfplugin5.Stop_Response) *tfprotov5.StopProviderResponse {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.StopProviderResponse{
		Error: in.Error,
	}

	return resp
}
