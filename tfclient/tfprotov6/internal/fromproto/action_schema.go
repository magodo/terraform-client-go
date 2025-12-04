package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func ActionSchema(in *tfplugin6.ActionSchema) (*tfprotov6.ActionSchema, error) {
	schema, err := Schema(in.GetSchema())
	if err != nil {
		return nil, err
	}
	resp := &tfprotov6.ActionSchema{
		Schema: schema,
	}
	return resp, nil
}
