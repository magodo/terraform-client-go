package fromproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func FunctionMetadata(in *tfplugin6.GetMetadata_FunctionMetadata) *tfprotov6.FunctionMetadata {
	if in == nil {
		return nil
	}

	return &tfprotov6.FunctionMetadata{
		Name: in.Name,
	}
}

func CallFunctionRequest(in *tfplugin6.CallFunction_Request) *tfprotov6.CallFunctionRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.CallFunctionRequest{
		Arguments: make([]*tfprotov6.DynamicValue, 0, len(in.Arguments)),
		Name:      in.Name,
	}

	for _, argument := range in.Arguments {
		resp.Arguments = append(resp.Arguments, DynamicValue(argument))
	}

	return resp
}

func CallFunctionResponse(in *tfplugin6.CallFunction_Response) *tfprotov6.CallFunctionResponse {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.CallFunctionResponse{
		Error:  FunctionError(in.Error),
		Result: DynamicValue(in.Result),
	}

	return resp
}

func GetFunctionsRequest(in *tfplugin6.GetFunctions_Request) *tfprotov6.GetFunctionsRequest {
	if in == nil {
		return nil
	}

	resp := &tfprotov6.GetFunctionsRequest{}

	return resp
}

func GetFunctionsResponse(in *tfplugin6.GetFunctions_Response) (*tfprotov6.GetFunctionsResponse, error) {
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

	resp := &tfprotov6.GetFunctionsResponse{
		Diagnostics: diags,
		Functions:   funcs,
	}

	return resp, nil
}

func FunctionError(in *tfplugin6.FunctionError) *tfprotov6.FunctionError {
	if in == nil {
		return nil
	}
	return &tfprotov6.FunctionError{
		Text:             in.Text,
		FunctionArgument: in.FunctionArgument,
	}
}

func Functions(in map[string]*tfplugin6.Function) (map[string]*tfprotov6.Function, error) {
	if in == nil {
		return nil, nil
	}

	out := map[string]*tfprotov6.Function{}
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

		out[k] = &tfprotov6.Function{
			Parameters:         params,
			VariadicParameter:  param,
			Return:             ret,
			Summary:            v.Summary,
			Description:        v.Description,
			DescriptionKind:    tfprotov6.StringKind(v.DescriptionKind),
			DeprecationMessage: v.DeprecationMessage,
		}
	}

	return out, nil
}

func FunctionParameters(in []*tfplugin6.Function_Parameter) ([]*tfprotov6.FunctionParameter, error) {
	if in == nil {
		return nil, nil
	}

	var out []*tfprotov6.FunctionParameter
	for _, p := range in {
		pp, err := FunctionParameter(p)
		if err != nil {
			return nil, err
		}
		out = append(out, pp)
	}
	return out, nil
}

func FunctionParameter(in *tfplugin6.Function_Parameter) (*tfprotov6.FunctionParameter, error) {
	if in == nil {
		return nil, nil
	}

	typ, err := tftypes.ParseJSONType(in.Type)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.FunctionParameter{
		AllowNullValue:     in.AllowNullValue,
		AllowUnknownValues: in.AllowUnknownValues,
		Description:        in.Description,
		DescriptionKind:    tfprotov6.StringKind(in.DescriptionKind),
		Name:               in.Name,
		Type:               typ,
	}, nil
}

func FunctionReturn(in *tfplugin6.Function_Return) (*tfprotov6.FunctionReturn, error) {
	if in == nil {
		return nil, nil
	}

	typ, err := tftypes.ParseJSONType(in.Type)
	if err != nil {
		return nil, err
	}

	return &tfprotov6.FunctionReturn{
		Type: typ,
	}, nil
}
