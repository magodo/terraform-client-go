package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func ActionSchema(in *tfplugin5.ActionSchema) (*tfprotov5.ActionSchema, error) {
	schema, err := Schema(in.GetSchema())
	if err != nil {
		return nil, err
	}
	resp := &tfprotov5.ActionSchema{
		Schema: schema,
	}

	return resp, nil
}

