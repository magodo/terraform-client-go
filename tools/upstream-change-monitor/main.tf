// This module monitors the depending upstream files
// Run `terraform apply` or `terraform plan` to see any change comparing to the main branch of the depending repos.

terraform {
  required_providers {
    restful = {
      source = "magodo/restful"
    }
  }
}

variable "token" {
  type = string
}

provider "restful" {
  base_url = "https://api.github.com"
  security = {
    http = {
      token = {
        token = var.token
      }
    }
  }
}

locals {
  terraform-plugin-go-commit = "v0.23.0"
  terraform-plugin-go = toset([
    "tfprotov5/internal/tfplugin5/tfplugin5.proto",
    "tfprotov6/internal/tfplugin6/tfplugin6.proto",

    "tfprotov5/internal/fromproto/attribute_path.go",
    "tfprotov5/internal/fromproto/data_source.go",
    "tfprotov5/internal/fromproto/diagnostic.go",
    "tfprotov5/internal/fromproto/function.go",
    "tfprotov5/internal/fromproto/provider.go",
    "tfprotov5/internal/fromproto/resource.go",
    "tfprotov5/internal/fromproto/schema.go",
    "tfprotov5/internal/fromproto/server_capabilities.go",
    "tfprotov5/internal/fromproto/state.go",
    "tfprotov5/internal/fromproto/string_kind.go",
    "tfprotov5/internal/fromproto/types.go",

    "tfprotov5/internal/toproto/attribute_path.go",
    "tfprotov5/internal/toproto/data_source.go",
    "tfprotov5/internal/toproto/diagnostic.go",
    "tfprotov5/internal/toproto/dynamic_value.go",
    "tfprotov5/internal/toproto/function.go",
    "tfprotov5/internal/toproto/provider.go",
    "tfprotov5/internal/toproto/resource.go",
    "tfprotov5/internal/toproto/schema.go",
    "tfprotov5/internal/toproto/server_capabilities.go",
    "tfprotov5/internal/toproto/state.go",
    "tfprotov5/internal/toproto/string_kind.go",

    "tfprotov6/internal/fromproto/attribute_path.go",
    "tfprotov6/internal/fromproto/data_source.go",
    "tfprotov6/internal/fromproto/diagnostic.go",
    "tfprotov6/internal/fromproto/function.go",
    "tfprotov6/internal/fromproto/provider.go",
    "tfprotov6/internal/fromproto/resource.go",
    "tfprotov6/internal/fromproto/schema.go",
    "tfprotov6/internal/fromproto/server_capabilities.go",
    "tfprotov6/internal/fromproto/state.go",
    "tfprotov6/internal/fromproto/string_kind.go",
    "tfprotov6/internal/fromproto/types.go",

    "tfprotov6/internal/toproto/attribute_path.go",
    "tfprotov6/internal/toproto/data_source.go",
    "tfprotov6/internal/toproto/diagnostic.go",
    "tfprotov6/internal/toproto/dynamic_value.go",
    "tfprotov6/internal/toproto/function.go",
    "tfprotov6/internal/toproto/provider.go",
    "tfprotov6/internal/toproto/resource.go",
    "tfprotov6/internal/toproto/schema.go",
    "tfprotov6/internal/toproto/server_capabilities.go",
    "tfprotov6/internal/toproto/state.go",
    "tfprotov6/internal/toproto/string_kind.go",
  ])
  terraform-commit = "v1.10.0-alpha20240730"
  terraform = toset([
    "internal/providers/provider.go",
    "internal/providers/functions.go",
    "internal/configs/configschema/decoder_spec.go",
    "internal/configs/configschema/empty_value.go",
    "internal/configs/configschema/implied_type.go",
    "internal/plugin/convert/schema.go",
    "internal/plugin/convert/deferred.go",
    "internal/plugin/grpc_provider.go",
    "internal/plugin6/convert/schema.go",
    "internal/plugin6/convert/deferred.go",
    "internal/plugin6/grpc_provider.go",
  ])
}

data "restful_resource" "terraform-plugin-go-used" {
  for_each        = local.terraform-plugin-go
  id              = "/repos/hashicorp/terraform-plugin-go/contents/${each.value}"
  allow_not_exist = true
  query = {
    ref = [local.terraform-plugin-go-commit]
  }
}

data "restful_resource" "terraform-plugin-go-main" {
  for_each        = local.terraform-plugin-go
  id              = "/repos/hashicorp/terraform-plugin-go/contents/${each.value}"
  allow_not_exist = true
  query = {
    ref = ["main"]
  }
}

data "restful_resource" "terraform-used" {
  for_each        = local.terraform
  id              = "/repos/hashicorp/terraform/contents/${each.value}"
  allow_not_exist = true
  query = {
    ref = [local.terraform-commit]
  }
}

data "restful_resource" "terraform-main" {
  for_each        = local.terraform
  id              = "/repos/hashicorp/terraform/contents/${each.value}"
  allow_not_exist = true
  query = {
    ref = ["main"]
  }
}

output "terraform-plugin-go-changes" {
  value = [
    for f in local.terraform-plugin-go : f if
    data.restful_resource.terraform-plugin-go-used[f].output != null && data.restful_resource.terraform-plugin-go-main[f].output != null ?
    # Content diff
    data.restful_resource.terraform-plugin-go-used[f].output.sha != data.restful_resource.terraform-plugin-go-main[f].output.sha
    :
    # Existance diff
    data.restful_resource.terraform-plugin-go-used[f].output == null != data.restful_resource.terraform-plugin-go-main[f].output == null
  ]
}

output "terraform-changes" {
  value = [
    for f in local.terraform : f if
    data.restful_resource.terraform-used[f].output != null && data.restful_resource.terraform-main[f].output != null ?
    # Content diff
    data.restful_resource.terraform-used[f].output.sha != data.restful_resource.terraform-main[f].output.sha
    :
    # Existance diff
    data.restful_resource.terraform-used[f].output == null != data.restful_resource.terraform-main[f].output == null
  ]
}
