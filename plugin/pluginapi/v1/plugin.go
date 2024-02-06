package pluginapiv1

import (
	context "context"
	"fmt"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	grpc "google.golang.org/grpc"

	"github.com/blackstork-io/fabric/plugin"
)

var handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PLUGINS_FOR",
	MagicCookieValue: "fabric",
}

type grpcPlugin struct {
	goplugin.Plugin
	schema *plugin.Schema
}

func (p *grpcPlugin) GRPCServer(broker *goplugin.GRPCBroker, s *grpc.Server) error {
	RegisterPluginServiceServer(s, &grpcServer{
		schema: p.schema,
	})
	return nil
}

func (p *grpcPlugin) GRPCClient(ctx context.Context, broker *goplugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	client := NewPluginServiceClient(c)
	res, err := client.GetSchema(ctx, &GetSchemaRequest{})
	if err != nil {
		return nil, err
	}
	schema, err := decodeSchema(res.Schema)
	if err != nil {
		return nil, err
	}
	for name, ds := range schema.DataSources {
		if ds == nil {
			return nil, fmt.Errorf("nil data source")
		}
		ds.DataFunc = p.clientDataFunc(name, client)
	}
	for name, cg := range schema.ContentProviders {
		if cg == nil {
			return nil, fmt.Errorf("nil content provider")
		}
		cg.ContentFunc = p.clientGenerateFunc(name, client)
	}
	return schema, nil
}

func (p *grpcPlugin) clientGenerateFunc(name string, client PluginServiceClient) plugin.ProvideContentFunc {
	return func(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.Content, hcl.Diagnostics) {
		if params == nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Nil params",
				Detail:   "Nil params",
			}}
		}
		cfgEncoded, err := encodeCtyValue(params.Config)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to encode config",
				Detail:   err.Error(),
			}}
		}
		argsEncoded, err := encodeCtyValue(params.Args)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to encode args",
				Detail:   err.Error(),
			}}
		}
		res, err := client.ProvideContent(ctx, &ProvideContentRequest{
			Provider:    name,
			Config:      cfgEncoded,
			Args:        argsEncoded,
			DataContext: encodeMapData(params.DataContext),
		})
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to generate content",
				Detail:   err.Error(),
			}}
		}
		content := decodeContent(res.Content)
		diags := decodeDiagnosticList(res.Diagnostics)
		return content, diags
	}
}

func (p *grpcPlugin) clientDataFunc(name string, client PluginServiceClient) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
		if params == nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Nil params",
				Detail:   "Nil params",
			}}
		}
		cfgEncoded, err := encodeCtyValue(params.Config)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to encode config",
				Detail:   err.Error(),
			}}
		}
		argsEncoded, err := encodeCtyValue(params.Args)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to encode args",
				Detail:   err.Error(),
			}}
		}

		res, err := client.RetrieveData(ctx, &RetrieveDataRequest{
			Source: name,
			Config: cfgEncoded,
			Args:   argsEncoded,
		})
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to fetch data",
				Detail:   err.Error(),
			}}
		}
		data := decodeData(res.Data)
		diags := decodeDiagnosticList(res.Diagnostics)
		return data, diags
	}
}
