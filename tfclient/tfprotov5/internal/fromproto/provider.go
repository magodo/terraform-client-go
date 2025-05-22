package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

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
		EphemeralResources: make([]tfprotov5.EphemeralResourceMetadata, 0, len(in.EphemeralResources)),
	}

	for _, datasource := range in.DataSources {
		if v := DataSourceMetadata(datasource); v != nil {
			resp.DataSources = append(resp.DataSources, *v)
		}
	}

	for _, ephemeralResource := range in.EphemeralResources {
		if v := GetMetadata_EphemeralResourceMetadata(ephemeralResource); v != nil {
			resp.EphemeralResources = append(resp.EphemeralResources, *v)
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

func GetProviderSchemaResponse(in *tfplugin5.GetProviderSchema_Response) (*tfprotov5.GetProviderSchemaResponse, error) {
	provider, err := Schema(in.Provider)
	if err != nil {
		return nil, err
	}

	var providerMeta *tfprotov5.Schema
	if in.ProviderMeta != nil {
		providerMeta, err = Schema(in.ProviderMeta)
		if err != nil {
			return nil, err
		}
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

	ephemeralResourceSchemas := make(map[string]*tfprotov5.Schema, len(in.EphemeralResourceSchemas))
	for k, v := range in.EphemeralResourceSchemas {
		schema, err := Schema(v)
		if err != nil {
			return nil, err
		}
		ephemeralResourceSchemas[k] = schema
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
		ServerCapabilities:       ServerCapabilities(in.ServerCapabilities),
		Provider:                 provider,
		ProviderMeta:             providerMeta,
		ResourceSchemas:          resourceSchemas,
		DataSourceSchemas:        dataSourceSchemas,
		Functions:                funcs,
		EphemeralResourceSchemas: ephemeralResourceSchemas,
		Diagnostics:              diags,
	}

	return resp, nil
}

func GetResourceIdentitySchemasResponse(in *tfplugin5.GetResourceIdentitySchemas_Response) (*tfprotov5.GetResourceIdentitySchemasResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.GetResourceIdentitySchemasResponse{
		Diagnostics:     diags,
		IdentitySchemas: make(map[string]*tfprotov5.ResourceIdentitySchema, len(in.IdentitySchemas)),
	}

	for name, schema := range in.IdentitySchemas {
		resp.IdentitySchemas[name], err = ResourceIdentitySchema(schema)
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
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

func ConfigureProviderResponse(in *tfplugin5.Configure_Response) (*tfprotov5.ConfigureProviderResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	return &tfprotov5.ConfigureProviderResponse{
		Diagnostics: diags,
	}, nil
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
