package plugin

import (
	"net/rpc"

	hcplugin "github.com/hashicorp/go-plugin"
)

// ServicePack is the interface that we're exposing as a plugin.
type ServicePack interface {
	RunProbes() error
}

// ServicePackRPC is an implementation that talks over RPC
type ServicePackRPC struct{ client *rpc.Client }

// RunProbes returns a message
func (g *ServicePackRPC) RunProbes() error {
	var err error
	return g.client.Call("Plugin.RunProbes", new(interface{}), &err)
}

// ServicePackRPCServer is the RPC server that ServicePackRPC talks to, conforming to
// the requirements of net/rpc
type ServicePackRPCServer struct {
	// This is the real implementation
	Impl ServicePack
}

// RunProbes is a wrapper for interface implementation
func (s *ServicePackRPCServer) RunProbes(args interface{}, resp *error) error {
	*resp = s.Impl.RunProbes()
	return *resp
}

// ServicePackPlugin is the implementation of plugin.Plugin so we can serve/consume this
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

// Server implements RPC server
func (p *ServicePackPlugin) Server(*hcplugin.MuxBroker) (interface{}, error) {
	return &ServicePackRPCServer{Impl: p.Impl}, nil
}

// Client implements RPC client
func (ServicePackPlugin) Client(b *hcplugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &ServicePackRPC{client: c}, nil
}
