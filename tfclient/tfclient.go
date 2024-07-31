package tfclient

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
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

// New creates a normalized client. It spins up an un-configured provider server, whose lifecycle is managed by the client, so make sure to call the "Kill" method on exit.
func New(opts Option) (Client, error) {
	c, v, err := newRaw(opts)
	if err != nil {
		return nil, err
	}
	switch v {
	case 5:
		return tf5client.New(c.pluginClient, c.v5client)
	case 6:
		return tf6client.New(c.pluginClient, c.v6client)
	default:
		return nil, fmt.Errorf("unsupported protocol version %d", v)
	}
}

// NewRaw creates a raw client. It spins up an un-configured provider server, whose lifecycle is managed by the client, so make sure to call the "Kill" method on exit.
func NewRaw(opts Option) (*RawClient, error) {
	c, _, err := newRaw(opts)
	return c, err
}

func newRaw(opts Option) (*RawClient, int, error) {
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

	if reattach := opts.Reattach; reattach == nil {
		config.VersionedPlugins = versionedPlugins
	} else {
		// Reference: https://github.com/hashicorp/terraform/blob/15ecdb66c84cd8202b0ae3d34c44cb4bbece5444/internal/command/meta_providers.go#L425
		if pv := reattach.ProtocolVersion; pv == 0 {
			config.Plugins = versionedPlugins[5]
		} else if plugins, ok := versionedPlugins[pv]; ok {
			config.Plugins = plugins
		} else {
			return nil, 0, fmt.Errorf("no supported plugins for protocol %d", pv)
		}
	}

	var client RawClient

	pclient := plugin.NewClient(config)

	client.pluginClient = pclient

	rpcClient, err := pclient.Client()
	if err != nil {
		return nil, 0, err
	}

	raw, err := rpcClient.Dispense("provider")
	if err != nil {
		return nil, 0, err
	}

	protoVer := pclient.NegotiatedVersion()
	switch protoVer {
	case 5:
		p := raw.(tfprotov5.ProviderServer)
		client.v5client = p
		return &client, 5, nil
	case 6:
		p := raw.(tfprotov6.ProviderServer)
		client.v6client = p
		return &client, 6, nil
	default:
		return nil, 0, fmt.Errorf("unsupported protocol version %d", protoVer)
	}
}

func ParseReattach(in string) (*plugin.ReattachConfig, error) {
	if in == "" {
		return nil, nil
	}

	type reattachConfig struct {
		Protocol        string
		ProtocolVersion int
		Addr            struct {
			Network string
			String  string
		}
		Pid  int
		Test bool
	}
	var m map[string]reattachConfig
	err := json.Unmarshal([]byte(in), &m)
	if err != nil {
		return nil, fmt.Errorf("Invalid format for TF_REATTACH_PROVIDERS: %w", err)
	}
	if len(m) != 1 {
		return nil, fmt.Errorf("expect only one of provider specified in the TF_REATTACH_PROVIDERS, got=%d", len(m))
	}

	var c reattachConfig
	var p string
	for k, v := range m {
		c = v
		p = k
	}

	var addr net.Addr
	switch c.Addr.Network {
	case "unix":
		addr, err = net.ResolveUnixAddr("unix", c.Addr.String)
		if err != nil {
			return nil, fmt.Errorf("Invalid unix socket path %q: %w", c.Addr.String, err)
		}
	case "tcp":
		addr, err = net.ResolveTCPAddr("tcp", c.Addr.String)
		if err != nil {
			return nil, fmt.Errorf("Invalid TCP address %q: %w", c.Addr.String, err)
		}
	default:
		return nil, fmt.Errorf("Unknown address type %q for %q", c.Addr.Network, p)
	}
	return &plugin.ReattachConfig{
		Protocol:        plugin.Protocol(c.Protocol),
		ProtocolVersion: c.ProtocolVersion,
		Pid:             c.Pid,
		Test:            c.Test,
		Addr:            addr,
	}, nil
}
