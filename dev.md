# Guide for Updating `terraform-plugin-go`

Since this project is the client implementation of the terraform protocol, it requires regular updates of the `terraform-plugin-go`. This section describes the process of this update. 

## 1. Update the Dependency

This is a regular Go dependency update:

```shell
$ go get github.com/hashicorp/terraform-plugin-go@<version>
$ go mod tidy
```

## 2. Update the Protocol

This project is vendoring and generating the protobuf Go code, therefore we need to copy the corresponding proto definition from the target version of `terraform-plugin-go` repo. These proto files are copied to `./tfclient/tfprotov[5|6]/internal/tfplugin[5|6]`.

After copying the proto files, regenerate the Go code via `make protoc`.

## 3. Update the Clients

There are two sets of clients need to update: the gRPC client and the user facing client.

### 3.1 Update the gRPC Client

These are the clients defined at `./tfclient/tfprotov[5|6]/tf[5|6]client/grpc_client.go`. They implemented the protobuf server(s). So the types being used are `github.com/hashicorp/terraform-plugin-go/tfprotov[5|6]`. For any new features added in the protobuf, add them to this client first.

#### 3.1.1 Update the gRPC Type Conversions

During the gRPC client updates above, there requires type conversion between the `github.com/hashicorp/terraform-plugin-go/tfprotov[5|6]` and the underlying protobuf types `github.com/magodo/terraform-client-go/tfclient/tfprotov[5|6]/internal/tfplugin[5|6]`. These conversions are defined at `./tfclient/tfprotov[5|6]/internal/[from|to]proto`.

When implementing them, the step is a bit anti-intuitive. E.g. when implementing the `fromproto`, the step is:

- Mimic the same file in `github.com/hashicorp/terraform-plugin-go/tfprotov5/internal/toproto` (**Note this is `toproto`**)
- Revert each function's input and output types

The sanity behind this is as below:

```
  ┌──────────────────────────┐  
  │    terraform-client-go   │  
  │            ▲             │  
  │            │ output      │  
  │         fromproto        │  
  │            │ input*      │  
  │          gRPC  |         │  
  └────────────┬───|─────────┘  
               │   |            
               │   |         
  ┌────────────┼───|─────────┐  
  │            ▲   |         │  
  │          gRPC  |         │  
  │            │ output*     │  
  │         toproto          │  
  │            │ input       │  
  │    terraform-plugin-go   │  
  └──────────────────────────┘  
```
                                
The code we are referring to is as the gRPC server, while we are as the gRPC client, therefore we look at `toproto` for `fromproto`. The same opposite direction applies to the input/output types of the function.

### 3.2 Update the User Facing Client

These are the clients defined at `./tfclient/tfprotov[5|6]/tf[5|6]client/client.go`. They are built on top of the gRPC client above, but faces to the API users, which means:

- Types of `typ` packages are used as the function signature
- The `typ.Diagnostics` is used as a second return value (if any), for better user experience

The code is mimicing the files in `terraform/internal/plugin(6)/grpc_provider.go`.

## 4. (Opt.) Add New Command

If the update introduces a new RPC interface, which is meaningful to the users, consider to implement a new `cmd` under `./cmd/` folder. This is also a good way to test the whole code path.
