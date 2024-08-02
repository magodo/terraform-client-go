package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func RawState(in *tfprotov6.RawState) *tfplugin6.RawState {
	if in == nil {
		return nil
	}
	return &tfplugin6.RawState{Json: in.JSON, Flatmap: in.Flatmap}
}
