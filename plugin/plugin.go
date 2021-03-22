package plugin

import (
	"net/rpc"

	hcplugin "github.com/hashicorp/go-plugin"
)

// ServicePack is the interface that we're exposing as a plugin.
type ServicePack interface {
	Greet() string
	// TODO: RunAllProbes
}

// Here is an implementation that talks over RPC
type ServicePackRPC struct{ client *rpc.Client }

func (g *ServicePackRPC) Greet() string {
	var resp string
	err := g.client.Call("Plugin.Greet", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that ServicePackRPC talks to, conforming to
// the requirements of net/rpc
type ServicePackRPCServer struct {
	// This is the real implementation
	Impl ServicePack
}

func (s *ServicePackRPCServer) Greet(args interface{}, resp *string) error {
	*resp = s.Impl.Greet()
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type ServicePackPlugin struct {
	// Impl Injection
	Impl ServicePack
}

func (p *ServicePackPlugin) Server(*hcplugin.MuxBroker) (interface{}, error) {
	return &ServicePackRPCServer{Impl: p.Impl}, nil
}

func (ServicePackPlugin) Client(b *hcplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ServicePackRPC{client: c}, nil
}
