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

func GetMetadata_Response(in *tfprotov5.GetMetadataResponse) *tfplugin5.GetMetadata_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetMetadata_Response{
		DataSources:        make([]*tfplugin5.GetMetadata_DataSourceMetadata, 0, len(in.DataSources)),
		Diagnostics:        Diagnostics(in.Diagnostics),
		Functions:          make([]*tfplugin5.GetMetadata_FunctionMetadata, 0, len(in.Functions)),
		Resources:          make([]*tfplugin5.GetMetadata_ResourceMetadata, 0, len(in.Resources)),
		ServerCapabilities: ServerCapabilities(in.ServerCapabilities),
	}

	for _, datasource := range in.DataSources {
		resp.DataSources = append(resp.DataSources, GetMetadata_DataSourceMetadata(&datasource))
	}

	for _, function := range in.Functions {
		resp.Functions = append(resp.Functions, GetMetadata_FunctionMetadata(&function))
	}

	for _, resource := range in.Resources {
		resp.Resources = append(resp.Resources, GetMetadata_ResourceMetadata(&resource))
	}

	return resp
}

func GetProviderSchema_Request(in *tfprotov5.GetProviderSchemaRequest) *tfplugin5.GetProviderSchema_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.GetProviderSchema_Request{}

	return req
}

func GetProviderSchema_Response(in *tfprotov5.GetProviderSchemaResponse) *tfplugin5.GetProviderSchema_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetProviderSchema_Response{
		DataSourceSchemas:  make(map[string]*tfplugin5.Schema, len(in.DataSourceSchemas)),
		Diagnostics:        Diagnostics(in.Diagnostics),
		Functions:          make(map[string]*tfplugin5.Function, len(in.Functions)),
		Provider:           Schema(in.Provider),
		ProviderMeta:       Schema(in.ProviderMeta),
		ResourceSchemas:    make(map[string]*tfplugin5.Schema, len(in.ResourceSchemas)),
		ServerCapabilities: ServerCapabilities(in.ServerCapabilities),
	}

	for name, schema := range in.ResourceSchemas {
		resp.ResourceSchemas[name] = Schema(schema)
	}

	for name, schema := range in.DataSourceSchemas {
		resp.DataSourceSchemas[name] = Schema(schema)
	}

	for name, function := range in.Functions {
		resp.Functions[name] = Function(function)
	}

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

func PrepareProviderConfig_Response(in *tfprotov5.PrepareProviderConfigResponse) *tfplugin5.PrepareProviderConfig_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.PrepareProviderConfig_Response{
		Diagnostics:    Diagnostics(in.Diagnostics),
		PreparedConfig: DynamicValue(in.PreparedConfig),
	}

	return resp
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

func Configure_Response(in *tfprotov5.ConfigureProviderResponse) *tfplugin5.Configure_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.Configure_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}

	return resp
}

func Stop_Request(in *tfprotov5.StopProviderRequest) *tfplugin5.Stop_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.Stop_Request{}

	return req
}

func Stop_Response(in *tfprotov5.StopProviderResponse) *tfplugin5.Stop_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.Stop_Response{
		Error: in.Error,
	}

	return resp
}
