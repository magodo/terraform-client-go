package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func StringKind(in tfprotov6.StringKind) tfplugin6.StringKind {
	return tfplugin6.StringKind(in)
}
