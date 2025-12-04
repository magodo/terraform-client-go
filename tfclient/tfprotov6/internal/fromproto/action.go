package fromproto

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ActionMetadata(in *tfplugin6.GetMetadata_ActionMetadata) *tfprotov6.ActionMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov6.ActionMetadata{
		TypeName: in.TypeName,
	}
}

func ValidateActionConfig_Response(in *tfplugin6.ValidateActionConfig_Response) (*tfprotov6.ValidateActionConfigResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.ValidateActionConfigResponse{
		Diagnostics: diags,
	}, nil
}

func PlanAction_Response(in *tfplugin6.PlanAction_Response) (*tfprotov6.PlanActionResponse, error) {
	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov6.PlanActionResponse{
		Diagnostics: diags,
		Deferred:    Deferred(in.Deferred),
	}

	return resp, nil
}

func InvokeAction_InvokeActionEvent(in *tfplugin6.InvokeAction_Event) (*tfprotov6.InvokeActionEvent, error) {
	switch event := (in.Type).(type) {
	case *tfplugin6.InvokeAction_Event_Progress_:
		return &tfprotov6.InvokeActionEvent{
			Type: &tfprotov6.ProgressInvokeActionEventType{
				Message: event.Progress.GetMessage(),
			},
		}, nil
	case *tfplugin6.InvokeAction_Event_Completed_:
		diags, err := Diagnostics(event.Completed.GetDiagnostics())
		if err != nil {
			return nil, err
		}
		return &tfprotov6.InvokeActionEvent{
			Type: &tfprotov6.CompletedInvokeActionEventType{
				Diagnostics: diags,
			},
		}, nil
	}

	// It is not currently possible to create tfprotov6.InvokeActionEventType
	// implementations outside the tfprotov6 package. If this panic was reached,
	// it implies that a new event type was introduced and needs to be implemented
	// as a new case above.
	panic(fmt.Sprintf("unimplemented tfprotov6.InvokeActionEventType type: %T", in.Type))
}
