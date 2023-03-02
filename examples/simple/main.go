package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfprotov5/tf5client"
	"github.com/magodo/terraform-client-go/tfprotov6/tf6client"
)

// Handshake is the HandshakeConfig used to configure clients and servers.
// Copied from: https://github.com/hashicorp/terraform/blob/a9230c9e7582c353c224cf0f4832d472ce042c0d/internal/plugin/serve.go#L22
var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  4,
	MagicCookieKey:   "TF_PLUGIN_MAGIC_COOKIE",
	MagicCookieValue: "d602bf8f470bc67ca7faa0386276bbdd4330efaf76d1a219cb4d6991ca9872b2",
}

// VersionedPlugins includes both protocol 5 and 6 because this is the function
// called in providerFactory (command/meta_providers.go) to set up the initial
// plugin client config.
// Copied from: https://github.com/hashicorp/terraform/blob/a9230c9e7582c353c224cf0f4832d472ce042c0d/internal/plugin/plugin.go#L11
var VersionedPlugins = map[int]plugin.PluginSet{
	5: {"provider": &tf5client.GRPCClientPlugin{}},
	6: {"provider": &tf6client.GRPCClientPlugin{}},
}

func main() {
	pluginPath := flag.String("path", "", "The path to the plugin")
	flag.Parse()
	client, err := buildClient(*pluginPath)
	if err != nil {
		log.Fatal(err)
	}
	switch client := client.(type) {
	case tfprotov5.ProviderServer:
		resp, err := client.GetProviderSchema(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Provider version: %d\n", resp.Provider.Version)
	case tfprotov6.ProviderServer:
		resp, err := client.GetProviderSchema(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Provider version: %d\n", resp.Provider.Version)
	}
}

func buildClient(pluginPath string) (interface{}, error) {
	config := &plugin.ClientConfig{
		HandshakeConfig:  Handshake,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Cmd:              exec.Command(pluginPath),
		//AutoMTLS:         true,
		VersionedPlugins: VersionedPlugins,
	}

	client := plugin.NewClient(config)
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense("provider")
	if err != nil {
		return nil, err
	}

	protoVer := client.NegotiatedVersion()
	switch protoVer {
	case 5:
		p := raw.(tfprotov5.ProviderServer)
		return p, nil
	case 6:
		p := raw.(tfprotov6.ProviderServer)
		return p, nil
	default:
		panic("unsupported protocol version")
	}
}
