package tf6client

import (
	"context"
	"errors"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/internal/tfplugin6"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

type GRPCClientPlugin struct {
	plugin.Plugin
}

func (p *GRPCClientPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return nil, errors.New("terraform-client-go only implements gRPC clients")
}

func (p *GRPCClientPlugin) Client(*plugin.MuxBroker, *rpc.Client) (interface{}, error) {
	return nil, errors.New("terraform-client-go only implements gRPC clients")
}

func (p *GRPCClientPlugin) GRPCClient(_ context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{
		client: tfplugin6.NewProviderClient(c),
	}, nil
}

func (p *GRPCClientPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	return errors.New("terraform-client-go only implements gRPC clients")
}
