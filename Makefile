tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go
	cd buildtools; go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

protoc:
	@cd tfclient/tfprotov5/internal/tfplugin5 && \
		protoc \
			--proto_path=. \
			--go_out=. \
			--go_opt=paths=source_relative \
			--go-grpc_out=. \
			--go-grpc_opt=paths=source_relative \
			tfplugin5.proto
	@cd tfclient/tfprotov6/internal/tfplugin6 && \
		protoc \
			--proto_path=. \
			--go_out=. \
			--go_opt=paths=source_relative \
			--go-grpc_out=. \
			--go-grpc_opt=paths=source_relative \
			tfplugin6.proto

.PHONY: protoc tools
