package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func RawState(in *tfplugin6.RawState) *tfprotov6.RawState {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.RawState{
		JSON:    in.Json,
		Flatmap: in.Flatmap,
	}

	return resp
}
