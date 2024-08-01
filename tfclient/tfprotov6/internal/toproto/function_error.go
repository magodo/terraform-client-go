package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func FunctionError(in *tfprotov6.FunctionError) *tfplugin6.FunctionError {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.FunctionError{
		FunctionArgument: in.FunctionArgument,
		Text:             ForceValidUTF8(in.Text),
	}

	return resp
}
