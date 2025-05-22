package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func FunctionMetadata(in *tfplugin5.GetMetadata_FunctionMetadata) *tfprotov5.FunctionMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov5.FunctionMetadata{
		Name: in.Name,
	}
}

func CallFunctionResponse(in *tfplugin5.CallFunction_Response) *tfprotov5.CallFunctionResponse {
	if in == nil {
		return nil
	}

	resp := &tfprotov5.CallFunctionResponse{
		Error:  FunctionError(in.Error),
		Result: DynamicValue(in.Result),
	}

	return resp
}

func GetFunctionsResponse(in *tfplugin5.GetFunctions_Response) (*tfprotov5.GetFunctionsResponse, error) {
	if in == nil {
		return nil, nil
	}

	diags, err := Diagnostics(in.Diagnostics)
	if err != nil {
		return nil, err
	}

	funcs, err := Functions(in.Functions)
	if err != nil {
		return nil, err
	}

	resp := &tfprotov5.GetFunctionsResponse{
		Diagnostics: diags,
		Functions:   funcs,
	}

	return resp, nil
}

func FunctionError(in *tfplugin5.FunctionError) *tfprotov5.FunctionError {
	if in == nil {
		return nil
	}
	return &tfprotov5.FunctionError{
		Text:             in.Text,
		FunctionArgument: in.FunctionArgument,
	}
}

func Functions(in map[string]*tfplugin5.Function) (map[string]*tfprotov5.Function, error) {
	if in == nil {
		return nil, nil
	}

	out := map[string]*tfprotov5.Function{}
	for k, v := range in {
		params, err := FunctionParameters(v.Parameters)
		if err != nil {
			return nil, err
		}
		param, err := FunctionParameter(v.VariadicParameter)
		if err != nil {
			return nil, err
		}
		ret, err := FunctionReturn(v.Return)
		if err != nil {
			return nil, err
		}

		out[k] = &tfprotov5.Function{
			Parameters:         params,
			VariadicParameter:  param,
			Return:             ret,
			Summary:            v.Summary,
			Description:        v.Description,
			DescriptionKind:    tfprotov5.StringKind(v.DescriptionKind),
			DeprecationMessage: v.DeprecationMessage,
		}
	}

	return out, nil
}

func FunctionParameters(in []*tfplugin5.Function_Parameter) ([]*tfprotov5.FunctionParameter, error) {
	if in == nil {
		return nil, nil
	}

	var out []*tfprotov5.FunctionParameter
	for _, p := range in {
		pp, err := FunctionParameter(p)
		if err != nil {
			return nil, err
		}
		out = append(out, pp)
	}
	return out, nil
}

func FunctionParameter(in *tfplugin5.Function_Parameter) (*tfprotov5.FunctionParameter, error) {
	if in == nil {
		return nil, nil
	}

	typ, err := tftypes.ParseJSONType(in.Type)
	if err != nil {
		return nil, err
	}

	return &tfprotov5.FunctionParameter{
		AllowNullValue:     in.AllowNullValue,
		AllowUnknownValues: in.AllowUnknownValues,
		Description:        in.Description,
		DescriptionKind:    tfprotov5.StringKind(in.DescriptionKind),
		Name:               in.Name,
		Type:               typ,
	}, nil
}

func FunctionReturn(in *tfplugin5.Function_Return) (*tfprotov5.FunctionReturn, error) {
	if in == nil {
		return nil, nil
	}

	typ, err := tftypes.ParseJSONType(in.Type)
	if err != nil {
		return nil, err
	}

	return &tfprotov5.FunctionReturn{
		Type: typ,
	}, nil
}
