package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func GetMetadata_DataSourceMetadata(in *tfprotov5.DataSourceMetadata) *tfplugin5.GetMetadata_DataSourceMetadata {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetMetadata_DataSourceMetadata{
		TypeName: in.TypeName,
	}

	return resp
}

func ValidateDataSourceConfig_Request(in *tfprotov5.ValidateDataSourceConfigRequest) *tfplugin5.ValidateDataSourceConfig_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ValidateDataSourceConfig_Request{
		TypeName: in.TypeName,
		Config:   DynamicValue(in.Config),
	}

	return req
}

func ValidateDataSourceConfig_Response(in *tfprotov5.ValidateDataSourceConfigResponse) *tfplugin5.ValidateDataSourceConfig_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ValidateDataSourceConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}

	return resp
}

func ReadDataSource_Request(in *tfprotov5.ReadDataSourceRequest) *tfplugin5.ReadDataSource_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ReadDataSource_Request{
		TypeName:     in.TypeName,
		Config:       DynamicValue(in.Config),
		ProviderMeta: DynamicValue(in.Config),
	}

	if in.ClientCapabilities != nil {
		req.ClientCapabilities = &tfplugin5.ClientCapabilities{
			DeferralAllowed: in.ClientCapabilities.DeferralAllowed,
		}
	}

	return req
}

func ReadDataSource_Response(in *tfprotov5.ReadDataSourceResponse) *tfplugin5.ReadDataSource_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ReadDataSource_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
		State:       DynamicValue(in.State),
		Deferred:    Deferred(in.Deferred),
	}

	return resp
}
