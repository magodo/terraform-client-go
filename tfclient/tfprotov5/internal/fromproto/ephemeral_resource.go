package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func GetMetadata_EphemeralResourceMetadata(in *tfplugin5.GetMetadata_EphemeralResourceMetadata) *tfprotov5.EphemeralResourceMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov5.EphemeralResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateEphemeralResourceConfig_Response(in *tfplugin5.ValidateEphemeralResourceConfig_Response) (*tfprotov5.ValidateEphemeralResourceConfigResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov5.ValidateEphemeralResourceConfigResponse{
		Diagnostics: diags,
	}, nil
}

func OpenEphemeralResource_Response(in *tfplugin5.OpenEphemeralResource_Response) (*tfprotov5.OpenEphemeralResourceResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov5.OpenEphemeralResourceResponse{
		Result:      DynamicValue(in.Result),
		Diagnostics: diags,
		Private:     in.Private,
		RenewAt:     Timestamp(in.RenewAt),
		Deferred:    Deferred(in.Deferred),
	}, nil
}

func RenewEphemeralResource_Response(in *tfplugin5.RenewEphemeralResource_Response) (*tfprotov5.RenewEphemeralResourceResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov5.RenewEphemeralResourceResponse{
		Diagnostics: diags,
		Private:     in.Private,
		RenewAt:     Timestamp(in.RenewAt),
	}, nil
}

func CloseEphemeralResource_Response(in *tfplugin5.CloseEphemeralResource_Response) (*tfprotov5.CloseEphemeralResourceResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov5.CloseEphemeralResourceResponse{
		Diagnostics: diags,
	}, nil
}
