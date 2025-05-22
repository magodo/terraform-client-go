package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func DataSourceMetadata(in *tfplugin5.GetMetadata_DataSourceMetadata) *tfprotov5.DataSourceMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov5.DataSourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateDataSourceConfigResponse(in *tfplugin5.ValidateDataSourceConfig_Response) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	return &tfprotov5.ValidateDataSourceConfigResponse{
		Diagnostics: diags,
	}, nil
}

func ReadDataSourceResponse(in *tfplugin5.ReadDataSource_Response) (*tfprotov5.ReadDataSourceResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}
	resp := &tfprotov5.ReadDataSourceResponse{
		Diagnostics: diags,
		State:       DynamicValue(in.State),
		Deferred:    Deferred(in.Deferred),
	}

	return resp, nil
}
