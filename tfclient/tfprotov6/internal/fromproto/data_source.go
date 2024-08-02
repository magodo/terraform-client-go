package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func DataSourceMetadata(in *tfplugin6.GetMetadata_DataSourceMetadata) *tfprotov6.DataSourceMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov6.DataSourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateDataResourceConfigRequest(in *tfplugin6.ValidateDataResourceConfig_Request) *tfprotov6.ValidateDataResourceConfigRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ValidateDataResourceConfigRequest{
		Config:   DynamicValue(in.Config),
		TypeName: in.TypeName,
	}

	return resp
}

func ValidateDataResourceConfigResponse(in *tfplugin6.ValidateDataResourceConfig_Response) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.ValidateDataResourceConfigResponse{
		Diagnostics: diags,
	}

	return resp, nil
}

func ReadDataSourceRequest(in *tfplugin6.ReadDataSource_Request) *tfprotov6.ReadDataSourceRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.ReadDataSourceRequest{
		Config:             DynamicValue(in.Config),
		ProviderMeta:       DynamicValue(in.ProviderMeta),
		TypeName:           in.TypeName,
		ClientCapabilities: ReadDataSourceClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func ReadDataSourceResponse(in *tfplugin6.ReadDataSource_Response) (*tfprotov6.ReadDataSourceResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.ReadDataSourceResponse{
		State:       DynamicValue(in.State),
		Diagnostics: diags,
		Deferred:    Deferred(in.Deferred),
	}

	return resp, nil
}
