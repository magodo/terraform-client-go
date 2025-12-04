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

func ReadDataSource_Request(in *tfprotov5.ReadDataSourceRequest) *tfplugin5.ReadDataSource_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.ReadDataSource_Request{
		TypeName:           in.TypeName,
		Config:             DynamicValue(in.Config),
		ProviderMeta:       DynamicValue(in.Config),
		ClientCapabilities: ReadDataSourceClientCapabilities(in.ClientCapabilities),
	}

	return req
}
