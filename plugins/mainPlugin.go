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
	Impl plugininterface.PluginRPCSer
}

func (s *RPCServer) Call(args plugininterface.ArgsSer, res *plugininterface.ResultSer) error {
	*res = s.Impl.Call(args)
	return nil
}

func (s *RPCServer) GetPlugins(_ struct{}, res *[]plugininterface.PluginSer) (err error) {
	*res = s.Impl.GetPlugins()
	return nil
}

// Adapter between plugin interface and net/rpc

type RPCClient struct {
	client *rpc.Client
}

// Call implements plugininterface.PluginRPCSer.
func (c *RPCClient) Call(args plugininterface.ArgsSer) (res plugininterface.ResultSer) {
	err := c.client.Call("Plugin.Call", args, &res)
	if err != nil {
		res.Diags = append(res.Diags, &plugininterface.RemoteDiag{
			Severity: hcl.DiagError,
			Summary:  "RPC Call error",
			Detail:   err.Error(),
		})
	}
	return
}

// GetPlugins implements plugininterface.PluginRPC.
func (c *RPCClient) GetPlugins() (res []plugininterface.PluginSer) {
	log.Println("RPCClient GetPlugins")
	err := c.client.Call("Plugin.GetPlugins", struct{}{}, &res)
	if err != nil {
		// TODO: hmmm, add diags/ error?
		log.Println("RPCClient GetPlugins Error:", err)
	}
	return
}

var _ plugininterface.PluginRPCSer = (*RPCClient)(nil)

// The go-plugin plugin, combines all above into one interface

type GoPlugin struct {
	Impl plugininterface.PluginRPCSer
}

var _ plugin.Plugin = (*GoPlugin)(nil)

func (p *GoPlugin) Server(_ *plugin.MuxBroker) (any, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (p *GoPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &RPCClient{client: c}, nil
}
