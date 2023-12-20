package data

import (
	"encoding/json"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Interface of the plugin

type Plugin interface {
	Execute(input any) (result any, err error)
}

type RPCServer struct {
	Impl Plugin
}

func (s *RPCServer) Execute(input []byte, resp *[]byte) (err error) {
	result, err := s.Impl.Execute(input)
	if err != nil {
		return err
	}
	*resp, err = json.Marshal(Result{
		Result: result,
	})
	return
}

// Adapter between plugin interface and net/rpc

type RPCClient struct {
	client *rpc.Client
}

var _ Plugin = (*RPCClient)(nil)

func (c *RPCClient) Execute(attrs any) (res any, err error) {
	attrsBytes, err := json.Marshal(attrs)
	if err != nil {
		return
	}
	var response []byte
	err = c.client.Call("Plugin.Execute", attrsBytes, &response)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &res)
	return
}

// The go-plugin plugin, combines all above into one interface

type GoPlugin struct {
	Impl Plugin
}

var _ plugin.Plugin = (*GoPlugin)(nil)

func (p *GoPlugin) Server(_ *plugin.MuxBroker) (any, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (p *GoPlugin) Client(_ *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &RPCClient{client: c}, nil
}

type Result struct {
	Result any `json:"result"`
}
