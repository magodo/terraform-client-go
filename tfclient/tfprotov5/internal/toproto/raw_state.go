package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func RawState(in *tfprotov5.RawState) *tfplugin5.RawState {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.RawState{
		Json:    in.JSON,
		Flatmap: in.Flatmap,
	}

	return resp
}
