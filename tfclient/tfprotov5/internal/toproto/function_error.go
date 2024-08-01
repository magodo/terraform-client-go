package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func FunctionError(in *tfprotov5.FunctionError) *tfplugin5.FunctionError {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.FunctionError{
		FunctionArgument: in.FunctionArgument,
		Text:             ForceValidUTF8(in.Text),
	}

	return resp
}
