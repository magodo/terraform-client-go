package toproto

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
)

func CallFunction_Request(in *tfprotov6.CallFunctionRequest) *tfplugin6.CallFunction_Request {
	if in == nil {
		return nil
	}

	var arguments []*tfplugin6.DynamicValue
	for _, v := range in.Arguments {
		arguments = append(arguments, DynamicValue(v))
	}

	req := &tfplugin6.CallFunction_Request{
		Name:      in.Name,
		Arguments: arguments,
	}

	return req
}

func Function(in *tfprotov6.Function) *tfplugin6.Function {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.Function{
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		DeprecationMessage: in.DeprecationMessage,
		Parameters:         make([]*tfplugin6.Function_Parameter, 0, len(in.Parameters)),
		Return:             Function_Return(in.Return),
		Summary:            in.Summary,
		VariadicParameter:  Function_Parameter(in.VariadicParameter),
	}

	for _, parameter := range in.Parameters {
		resp.Parameters = append(resp.Parameters, Function_Parameter(parameter))
	}

	return resp
}

func Function_Parameter(in *tfprotov6.FunctionParameter) *tfplugin6.Function_Parameter {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.Function_Parameter{
		AllowNullValue:     in.AllowNullValue,
		AllowUnknownValues: in.AllowUnknownValues,
		Description:        in.Description,
		DescriptionKind:    StringKind(in.DescriptionKind),
		Name:               in.Name,
		Type:               CtyType(in.Type),
	}

	return resp
}

func Function_Return(in *tfprotov6.FunctionReturn) *tfplugin6.Function_Return {
	if in == nil {
		return nil
	}

	resp := &tfplugin6.Function_Return{
		Type: CtyType(in.Type),
	}

	return resp
}

func GetFunctions_Request(in *tfprotov6.GetFunctionsRequest) *tfplugin6.GetFunctions_Request {
	if in == nil {
		return nil
	}

	req := &tfplugin6.GetFunctions_Request{}

	return req
}

func GetMetadata_FunctionMetadata(in *tfprotov6.FunctionMetadata) *tfplugin6.GetMetadata_FunctionMetadata {
	if in == nil {
		return nil
	}

	return &tfplugin6.GetMetadata_FunctionMetadata{
		Name: in.Name,
	}
}
