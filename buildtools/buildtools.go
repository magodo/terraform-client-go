//go:build tools

package buildtools

import (
	// Protocol Buffers compiler plugin for Go gRPC. This tool is versioned
	// separately from google.golang.org/grpc.
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
