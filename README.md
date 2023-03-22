# terraform-client-go

terraform-client-go provides low-level Go bindings for the Terraform plugin protocol, only as the **client**.

## Why

There is a hashicorp official project [https://github.com/hashicorp/terraform-plugin-go](https://github.com/hashicorp/terraform-plugin-go) that provides low-level Go bindings for the Terraform plugin protocol, while [it only support the **server** part](https://github.com/hashicorp/terraform-plugin-go/blob/edfd37f2ba46017fc4fec2bdfed2fd0ec091a333/tfprotov5/tf5server/plugin.go#L36-L37).

There is also a personal project from Mart [https://github.com/apparentlymart/terraform-provider](https://github.com/apparentlymart/terraform-provider) that does the similar thing as this project, while it seems not being actively maintained, and some of its dependencies that being exposed in some forms (e.g. user facing types) seems not the hashicorp offically recommended ones (e.g. https://github.com/rpcplugin/go against https://github.com/hashicorp/terraform-plugin-go, https://github.com/apparentlymart/terraform-schema-go against https://github.com/hashicorp/terraform-json).

## What

This project provides two layers of client abstractions:

- Raw client: The protocol specific, thin client. It simply builds on top of the `terraform-plugin-go` by implementing the `tfprotov{5,6}.ProviderServer` as the client, that is required by `terraform-plugin-go`. Therefore, the type system being exposed is defined by `terraform-plugin-go`.

    Example: https://github.com/magodo/terraform-client-go/tree/main/examples/raw

- Normalized client: The normalized client regardless of the protocol version being used (this is how terraform core handles protocol diff). It is built on top of the *raw client* by adding a type conversion layer to convert between the `terraform-plugin-go` and `cty` type systems. As a result, the type system being exposed by this client is defined by [`terraform-json`](https://github.com/hashicorp/terraform-json).

    Example: https://github.com/magodo/terraform-client-go/tree/main/examples/client

## How

There are a lot of code duplication&adoption from different sources, for a reason:

- https://github.com/hashicorp/terraform-plugin-go: The `tfprotov{5|6}/internal/{from|to}proto` is duplicated for type conversion between protobuf generated types and `terraform-plugin-go` types
- https://github.com/apparentlymart/terraform-provider: The diagnostics type definition and conversion for protobuf generated diagnostics type is duplicated
- https://github.com/hashicorp/terraform:
    - The normalized client interface (with a little difference on the signatures) 
    - The schema implied type based on hcldec
    - The type conversion between terraform core types and protobuf generated types is duplicated, but adopted for conversion between `terraform-json` types and `terraform-plugin-go` types.
    - The client interface implementations for the two protocols