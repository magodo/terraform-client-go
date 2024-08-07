package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ServerCapabilities(in *tfprotov6.ServerCapabilities) *tfplugin6.ServerCapabilities {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ServerCapabilities{
		GetProviderSchemaOptional: in.GetProviderSchemaOptional,
		MoveResourceState:         in.MoveResourceState,
		PlanDestroy:               in.PlanDestroy,
	}

	return resp
}
