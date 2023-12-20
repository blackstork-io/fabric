package content

import (
	"encoding/json"
	"fmt"
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Interface of the plugin

type Plugin interface {
	Execute(attrs, dict any) (string, error)
}

type RPCServer struct {
	Impl Plugin
}

type Args struct {
	Attrs   []byte
	Content []byte
}

func (s *RPCServer) Execute(input Args, resp *[]byte) (err error) {
	result, err := s.Impl.Execute(input.Attrs, input.Content)
	if err != nil {
		return err
	}
	*resp, err = json.Marshal(
		TextStruct{
			Text: string(result),
		},
	)
	return
}

// Adapter between plugin interface and net/rpc

type RPCClient struct {
	client *rpc.Client
}

var _ Plugin = (*RPCClient)(nil)

func (c *RPCClient) Execute(attrs, dict any) (res string, err error) {
	attrsBytes, err := json.Marshal(attrs)
	if err != nil {
		return
	}
	contentBytes, err := json.Marshal(dict)
	if err != nil {
		return
	}
	var response []byte
	err = c.client.Call(
		"Plugin.Execute",
		Args{
			Attrs:   attrsBytes,
			Content: contentBytes,
		},
		&response)
	if err != nil {
		return
	}
	var text TextStruct
	err = json.Unmarshal(response, &text)
	if err != nil {
		return
	}
	return text.Text, nil
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

type TextStruct struct {
	Text string `json:"text"`
}

func GetText(attrs any) (text string, err error) {
	m, ok := attrs.(map[string]any)
	if !ok {
		err = fmt.Errorf("failed to parse")
		return
	}
	text, ok = m["text"].(string)
	if !ok {
		err = fmt.Errorf("failed to parse")
		return
	}
	return
}
