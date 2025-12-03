package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ValidateActionConfigRequest(in *tfprotov5.ValidateActionConfigRequest) *tfplugin5.ValidateActionConfig_Request {
	return &tfplugin5.ValidateActionConfig_Request{
		ActionType: in.ActionType,
		Config:     DynamicValue(in.Config),
	}
}

func PlanActionRequest(in *tfprotov5.PlanActionRequest) *tfplugin5.PlanAction_Request {
	resp := &tfplugin5.PlanAction_Request{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: PlanActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}

func InvokeActionRequest(in *tfprotov5.InvokeActionRequest) *tfplugin5.InvokeAction_Request {
	resp := &tfplugin5.InvokeAction_Request{
		ActionType:         in.ActionType,
		Config:             DynamicValue(in.Config),
		ClientCapabilities: InvokeActionClientCapabilities(in.ClientCapabilities),
	}

	return resp
}
