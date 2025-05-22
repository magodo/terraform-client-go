package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/internal/tfplugin5"
)

func CallFunction_Request(in *tfprotov5.CallFunctionRequest) *tfplugin5.CallFunction_Request {
	if in == nil {
		return nil
	}

	var arguments []*tfplugin5.DynamicValue
	for _, v := range in.Arguments {
		arguments = append(arguments, DynamicValue(v))
	}

	req := &tfplugin5.CallFunction_Request{
		Name:      in.Name,
		Arguments: arguments,
	}

	return req
}

func Function(in *tfprotov5.Function) *tfplugin5.Function {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.Function{
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		DeprecationMessage: in.DeprecationMessage,
		Parameters:         make([]*tfplugin5.Function_Parameter, 0, len(in.Parameters)),
		Return:             Function_Return(in.Return),
		Summary:            in.Summary,
		VariadicParameter:  Function_Parameter(in.VariadicParameter),
	}

	for _, parameter := range in.Parameters {
		resp.Parameters = append(resp.Parameters, Function_Parameter(parameter))
	}

	return resp
}

func Function_Parameter(in *tfprotov5.FunctionParameter) *tfplugin5.Function_Parameter {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.Function_Parameter{
		AllowNullValue:     in.AllowNullValue,
		AllowUnknownValues: in.AllowUnknownValues,
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		Name:               in.Name,
		Type:               CtyType(in.Type),
	}

	return resp
}

func Function_Return(in *tfprotov5.FunctionReturn) *tfplugin5.Function_Return {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.Function_Return{
		Type: CtyType(in.Type),
	}

	return resp
}

func GetFunctions_Request(in *tfprotov5.GetFunctionsRequest) *tfplugin5.GetFunctions_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin5.GetFunctions_Request{}

	return req
}

func GetMetadata_FunctionMetadata(in *tfprotov5.FunctionMetadata) *tfplugin5.GetMetadata_FunctionMetadata {
	if in == nil {
		return nil
	}

	resp := &tfplugin5.GetMetadata_FunctionMetadata{
		Name: in.Name,
	}

	return resp
}
