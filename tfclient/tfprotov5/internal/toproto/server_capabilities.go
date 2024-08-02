package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ServerCapabilities(in *tfprotov5.ServerCapabilities) *tfplugin5.ServerCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.ServerCapabilities{
		GetProviderSchemaOptional: in.GetProviderSchemaOptional,
		MoveResourceState:         in.MoveResourceState,
		PlanDestroy:               in.PlanDestroy,
	}

	return resp
}
