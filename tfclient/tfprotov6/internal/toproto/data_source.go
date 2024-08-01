package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func GetMetadata_DataSourceMetadata(in *tfprotov6.DataSourceMetadata) *tfplugin6.GetMetadata_DataSourceMetadata {
	if in == nil {
		return nil
	}

	return &tfplugin6.GetMetadata_DataSourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateDataResourceConfig_Request(in *tfprotov6.ValidateDataResourceConfigRequest) *tfplugin6.ValidateDataResourceConfig_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin6.ValidateDataResourceConfig_Request{
		TypeName: in.TypeName,
		Config:   DynamicValue(in.Config),
	}

	return req
}

func ValidateDataResourceConfig_Response(in *tfprotov6.ValidateDataResourceConfigResponse) *tfplugin6.ValidateDataResourceConfig_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ValidateDataResourceConfig_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
	}

	return resp
}

func ReadDataSource_Request(in *tfprotov6.ReadDataSourceRequest) *tfplugin6.ReadDataSource_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin6.ReadDataSource_Request{
		TypeName:     in.TypeName,
		Config:       DynamicValue(in.Config),
		ProviderMeta: DynamicValue(in.ProviderMeta),
	}

	if in.ClientCapabilities != nil {
		req.ClientCapabilities = &tfplugin6.ClientCapabilities{
			DeferralAllowed: in.ClientCapabilities.DeferralAllowed,
		}
	}

	return req
}

func ReadDataSource_Response(in *tfprotov6.ReadDataSourceResponse) *tfplugin6.ReadDataSource_Response {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ReadDataSource_Response{
		Diagnostics: Diagnostics(in.Diagnostics),
		State:       DynamicValue(in.State),
		Deferred:    Deferred(in.Deferred),
	}

	return resp
}
