package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ResourceIdentityData(in *tfprotov6.ResourceIdentityData) *tfplugin6.ResourceIdentityData {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.ResourceIdentityData{
		IdentityData: DynamicValue(in.IdentityData),
	}

	return resp
}
