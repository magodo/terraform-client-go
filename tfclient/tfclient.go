package tfclient

import (
	"crypto/tls"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov5/tf5client"
	"github.com/magodo/terraform-client-go/tfclient/tfprotov6/tf6client"
	"google.golang.org/grpc"
)

type Client struct {
	pluginClient *plugin.Client

	// Either one of below will be non nil
	v5client tfprotov5.ProviderServer
	v6client tfprotov6.ProviderServer
}

type Option struct {
	// One of the following must be set, but not both.
	//
	// Cmd is the unstarted subprocess for starting the plugin. If this is
	// set, then the Client starts the plugin process on its own and connects
	// to it.
	//
	// Reattach is configuration for reattaching to an existing plugin process
	// that is already running. This isn't common.
	Cmd      *exec.Cmd
	Reattach *plugin.ReattachConfig

	// SecureConfig is configuration for verifying the integrity of the
	// executable. It can not be used with Reattach.
	SecureConfig *plugin.SecureConfig

	// TLSConfig is used to enable TLS on the RPC client.
	TLSConfig *tls.Config

	// The minimum and maximum port to use for communicating with
	// the subprocess. If not set, this defaults to 10,000 and 25,000
	// respectively.
	MinPort, MaxPort uint

	// StartTimeout is the timeout to wait for the plugin to say it
	// has started successfully.
	StartTimeout time.Duration

	// If non-nil, then the stderr of the client will be written to here
	// (as well as the log). This is the original os.Stderr of the subprocess.
	// This isn't the output of synced stderr.
	Stderr io.Writer

	// SyncStdout, SyncStderr can be set to override the
	// respective os.Std* values in the plugin. Care should be taken to
	// avoid races here. If these are nil, then this will be set to
	// ioutil.Discard.
	SyncStdout io.Writer
	SyncStderr io.Writer

	// Logger is the logger that the client will used. If none is provided,
	// it will default to hclog's default logger.
	Logger hclog.Logger

	// AutoMTLS has the client and server automatically negotiate mTLS for
	// transport authentication. This ensures that only the original client will
	// be allowed to connect to the server, and all other connections will be
	// rejected. The client will also refuse to connect to any server that isn't
	// the original instance started by the client.
	//
	// In this mode of operation, the client generates a one-time use tls
	// certificate, sends the public x.509 certificate to the new server, and
	// the server generates a one-time use tls certificate, and sends the public
	// x.509 certificate back to the client. These are used to authenticate all
	// rpc connections between the client and server.
	//
	// Setting AutoMTLS to true implies that the server must support the
	// protocol, and correctly negotiate the tls certificates, or a connection
	// failure will result.
	//
	// The client should not set TLSConfig, nor should the server set a
	// TLSProvider, because AutoMTLS implies that a new certificate and tls
	// configuration will be generated at startup.
	//
	// You cannot Reattach to a server with this option enabled.
	AutoMTLS bool

	// GRPCDialOptions allows plugin users to pass custom grpc.DialOption
	// to create gRPC connections. This only affects plugins using the gRPC
	// protocol.
	GRPCDialOptions []grpc.DialOption
}

// New spins up an un-configured provider server, whose lifecycle is managed by the client, so make sure to call the "Kill" method on exit.
func New(opts Option) (*Client, error) {
	// handshake is the HandshakeConfig used to configure clients and servers.
	// Reference: https://github.com/hashicorp/terraform/blob/a9230c9e7582c353c224cf0f4832d472ce042c0d/internal/plugin/serve.go#L22
	handshake := plugin.HandshakeConfig{
		MagicCookieKey:   "TF_PLUGIN_MAGIC_COOKIE",
		MagicCookieValue: "d602bf8f470bc67ca7faa0386276bbdd4330efaf76d1a219cb4d6991ca9872b2",
	}

	// versionedPlugins includes both protocol 5 and 6 because this is the function
	// called in providerFactory (command/meta_providers.go) to set up the initial
	// plugin client config.
	// Reference: https://github.com/hashicorp/terraform/blob/a9230c9e7582c353c224cf0f4832d472ce042c0d/internal/plugin/plugin.go#L11
	versionedPlugins := map[int]plugin.PluginSet{
		5: {"provider": &tf5client.GRPCClientPlugin{}},
		6: {"provider": &tf6client.GRPCClientPlugin{}},
	}

	config := &plugin.ClientConfig{
		HandshakeConfig:  handshake,
		VersionedPlugins: versionedPlugins,
		Cmd:              opts.Cmd,
		Reattach:         opts.Reattach,
		SecureConfig:     opts.SecureConfig,
		TLSConfig:        opts.TLSConfig,
		Managed:          false,
		MinPort:          opts.MinPort,
		MaxPort:          opts.MaxPort,
		StartTimeout:     opts.StartTimeout,
		Stderr:           opts.Stderr,
		SyncStdout:       opts.SyncStdout,
		SyncStderr:       opts.SyncStderr,
		AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
		Logger:           opts.Logger,
		AutoMTLS:         opts.AutoMTLS,
		GRPCDialOptions:  opts.GRPCDialOptions,
	}

	var client Client

	pclient := plugin.NewClient(config)

	client.pluginClient = pclient

	rpcClient, err := pclient.Client()
	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense("provider")
	if err != nil {
		return nil, err
	}

	protoVer := pclient.NegotiatedVersion()
	switch protoVer {
	case 5:
		p := raw.(tfprotov5.ProviderServer)
		client.v5client = p
		return &client, nil
	case 6:
		p := raw.(tfprotov6.ProviderServer)
		client.v6client = p
		return &client, nil
	default:
		return nil, fmt.Errorf("unsupported protocol version %d", protoVer)
	}
}

// AsV5Client returns the v5 client if the linked provider is running in protocol v5, otherwise return nil
func (c *Client) AsV5Client() tfprotov5.ProviderServer {
	return c.v5client
}

// AsV6Client returns the v6 client if the linked provider is running in protocol v6, otherwise return nil
func (c *Client) AsV6Client() tfprotov6.ProviderServer {
	return c.v6client
}

// Kill ends the executing subprocess (if it is running) and perform any cleanup
// tasks necessary such as capturing any remaining logs and so on.
//
// This method blocks until the process successfully exits.
//
// This method can safely be called multiple times.
func (c *Client) Kill() {
	c.pluginClient.Kill()
}
