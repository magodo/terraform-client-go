package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

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

func GetProviderSchemaResponse(in *tfplugin6.GetProviderSchema_Response) (*tfprotov6.GetProviderSchemaResponse, error) {
	provider, err := Schema(in.Provider)
	if err != nil {
		return nil, err
	}

	var providerMeta *tfprotov6.Schema
	if in.ProviderMeta != nil {
		providerMeta, err = Schema(in.ProviderMeta)
		if err != nil {
			return nil, err
		}
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

	ephemeralResourceSchemas := make(map[string]*tfprotov6.Schema, len(in.EphemeralResourceSchemas))
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

	resp := &tfprotov6.GetProviderSchemaResponse{
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

func GetResourceIdentitySchemasResponse(in *tfplugin6.GetResourceIdentitySchemas_Response) (*tfprotov6.GetResourceIdentitySchemasResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.GetResourceIdentitySchemasResponse{
		Diagnostics:     diags,
		IdentitySchemas: make(map[string]*tfprotov6.ResourceIdentitySchema, len(in.IdentitySchemas)),
	}

	for name, schema := range in.IdentitySchemas {
		sch, err := ResourceIdentitySchema(schema)
		if err != nil {
			return nil, err
		}
		resp.IdentitySchemas[name] = sch
	}

	return resp, nil
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

func ConfigureProviderResponse(in *tfplugin6.ConfigureProvider_Response) (*tfprotov6.ConfigureProviderResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	return &tfprotov6.ConfigureProviderResponse{
		Diagnostics: diags,
	}, nil
}

func StopProviderResponse(in *tfplugin6.StopProvider_Response) (*tfprotov6.StopProviderResponse, error) {
	return &tfprotov6.StopProviderResponse{
		Error: in.Error,
	}, nil
}
