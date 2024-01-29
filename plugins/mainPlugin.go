package plugins

import (
	"log"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"

	plugininterface "github.com/blackstork-io/fabric/plugininterface/v1"
)

// Interface of the plugin

type RPCServer struct {
	Impl plugininterface.PluginRPC
}

func (s *RPCServer) Call(args plugininterface.Args, res *plugininterface.Result) error {
	*res = s.Impl.Call(args)
	return nil
}

func (s *RPCServer) GetPlugins(_ struct{}, res *[]plugininterface.Plugin) error {
	*res = s.Impl.GetPlugins()
	return nil
}

// Adapter between plugin interface and net/rpc

type RPCClient struct {
	client *rpc.Client
}

// Call implements plugininterface.PluginRPC.
func (c *RPCClient) Call(args plugininterface.Args) (res plugininterface.Result) {
	err := c.client.Call("Plugin.Call", args, &res)
	if err != nil {
		res.Diags = res.Diags.Append(&hcl.Diagnostic{
			Severity: hcl.DiagError,
			Summary:  "RPC Call error",
			Detail:   err.Error(),
		})
	}
	return
}

// GetPlugins implements plugininterface.PluginRPC.
func (c *RPCClient) GetPlugins() (res []plugininterface.Plugin) {
	log.Println("RPCClient GetPlugins")
	err := c.client.Call("Plugin.GetPlugins", struct{}{}, &res)
	if err != nil {
		// TODO: hmmm, add diags/ error?
		log.Println("RPCClient GetPlugins Error:", err)
	}
	return
}

var _ plugininterface.PluginRPC = (*RPCClient)(nil)

// The go-plugin plugin, combines all above into one interface

type GoPlugin struct {
	Impl plugininterface.PluginRPC
}

var _ plugin.Plugin = (*GoPlugin)(nil)

func (p *GoPlugin) Server(_ *plugin.MuxBroker) (any, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (p *GoPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &RPCClient{client: c}, nil
}
