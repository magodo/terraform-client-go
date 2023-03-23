package tfclient

import (
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type RawClient struct {
	pluginClient *plugin.Client

	// Either one of below will be non nil
	v5client tfprotov5.ProviderServer
	v6client tfprotov6.ProviderServer
}

// AsV5Client returns the v5 client if the linked provider is running in protocol v5, otherwise return nil
func (c *RawClient) AsV5Client() tfprotov5.ProviderServer {
	return c.v5client
}

// AsV6Client returns the v6 client if the linked provider is running in protocol v6, otherwise return nil
func (c *RawClient) AsV6Client() tfprotov6.ProviderServer {
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
