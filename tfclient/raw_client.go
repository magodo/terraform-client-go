package tfclient

import (
	"github.com/hashicorp/go-plugin"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/tf5client"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/tf6client"
)

type RawClient struct {
	pluginClient *plugin.Client

	// Either one of below will be non nil
	v5client tf5client.TFProtoV5Client
	v6client tf6client.TFProtoV6Client
}

// AsV5Client returns the v5 client if the linked provider is running in protocol v5, otherwise return nil
func (c *RawClient) AsV5Client() tf5client.TFProtoV5Client {
	return c.v5client
}

// AsV6Client returns the v6 client if the linked provider is running in protocol v6, otherwise return nil
func (c *RawClient) AsV6Client() tf6client.TFProtoV6Client {
	return c.v6client
}

// Kill ends the executing subprocess (if it is running) and perform any cleanup
// tasks necessary such as capturing any remaining logs and so on.
//
// This method blocks until the process successfully exits.
//
// This method can safely be called multiple times.
func (c *RawClient) Kill() {
	c.pluginClient.Kill()
}
