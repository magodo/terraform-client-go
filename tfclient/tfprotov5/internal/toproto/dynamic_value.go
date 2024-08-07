package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func DynamicValue(in *tfprotov5.DynamicValue) *tfplugin5.DynamicValue {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.DynamicValue{
		Msgpack: in.MsgPack,
		Json:    in.JSON,
	}

	return resp
}

func CtyType(in tftypes.Type) []byte {
	if in == nil {
		return nil
	}

	// MarshalJSON is always error safe.
	// nolint:staticcheck // Intended first-party usage
	resp, _ := in.MarshalJSON()

	return resp
}
