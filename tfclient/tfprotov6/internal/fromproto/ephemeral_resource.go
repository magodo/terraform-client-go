package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func GetMetadata_EphemeralResourceMetadata(in *tfplugin6.GetMetadata_EphemeralResourceMetadata) *tfprotov6.EphemeralResourceMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov6.EphemeralResourceMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateEphemeralResourceConfig_Response(in *tfplugin6.ValidateEphemeralResourceConfig_Response) (*tfprotov6.ValidateEphemeralResourceConfigResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.ValidateEphemeralResourceConfigResponse{
		Diagnostics: diags,
	}, nil
}

func OpenEphemeralResource_Response(in *tfplugin6.OpenEphemeralResource_Response) (*tfprotov6.OpenEphemeralResourceResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.OpenEphemeralResourceResponse{
		Result:      DynamicValue(in.Result),
		Diagnostics: diags,
		Private:     in.Private,
		RenewAt:     Timestamp(in.RenewAt),
		Deferred:    Deferred(in.Deferred),
	}, nil
}

func RenewEphemeralResource_Response(in *tfplugin6.RenewEphemeralResource_Response) (*tfprotov6.RenewEphemeralResourceResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.RenewEphemeralResourceResponse{
		Diagnostics: diags,
		Private:     in.Private,
		RenewAt:     Timestamp(in.RenewAt),
	}, nil
}

func CloseEphemeralResource_Response(in *tfplugin6.CloseEphemeralResource_Response) (*tfprotov6.CloseEphemeralResourceResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.CloseEphemeralResourceResponse{
		Diagnostics: diags,
	}, nil
}
