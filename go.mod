module github.com/magodo/terraform-client-go

go 1.23.0

toolchain go1.24.1

require (
	github.com/apparentlymart/go-dump v0.0.0-20180507223929-23540a00eaa3
	github.com/davecgh/go-spew v1.1.1
	github.com/evanphx/json-patch/v5 v5.6.0
	github.com/google/go-cmp v0.7.0
	github.com/hashicorp/go-hclog v1.5.0
	github.com/hashicorp/go-plugin v1.6.3
	github.com/hashicorp/go-version v1.7.0
	github.com/hashicorp/hc-install v0.9.0
	github.com/hashicorp/hcl/v2 v2.11.1
	github.com/hashicorp/terraform-exec v0.21.0
	github.com/hashicorp/terraform-json v0.25.0
	github.com/hashicorp/terraform-plugin-go v0.28.0
	github.com/zclconf/go-cty v1.16.3
	github.com/zclconf/go-cty-debug v0.0.0-20240209213017-b8d9e32151be
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/vmihailenco/msgpack/v5 v5.4.1 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	golang.org/x/mod v0.21.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto v0.0.0-20230110181048-76db0878b65f // indirect
)

// Remove this once https://github.com/hashicorp/terraform-json/issues/164 is resolved.
replace github.com/hashicorp/terraform-json => github.com/magodo/terraform-json v0.13.1-0.20250523033318-1c9170c5e727
