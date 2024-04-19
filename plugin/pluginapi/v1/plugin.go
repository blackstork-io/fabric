package pluginapiv1

import (
	context "context"
	"fmt"
	"log/slog"
	"time"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	grpc "google.golang.org/grpc"

	"github.com/blackstork-io/fabric/plugin"
)

var defaultMsgSize = 1024 * 1024 * 20

var handshake = goplugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "PLUGINS_FOR",
	MagicCookieValue: "fabric",
}

type grpcPlugin struct {
	goplugin.Plugin
	logger *slog.Logger
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
	for name, pub := range schema.Publishers {
		if pub == nil {
			return nil, fmt.Errorf("nil publisher")
		}
		pub.PublishFunc = p.clientPublishFunc(name, client)
	}
	return schema, nil
}

func (p *grpcPlugin) callOptions() []grpc.CallOption {
	return []grpc.CallOption{
		grpc.MaxCallRecvMsgSize(defaultMsgSize),
		grpc.MaxCallSendMsgSize(defaultMsgSize),
	}
}

func (p *grpcPlugin) clientGenerateFunc(name string, client PluginServiceClient) plugin.ProvideContentFunc {
	return func(ctx context.Context, params *plugin.ProvideContentParams) (*plugin.ContentResult, hcl.Diagnostics) {
		p.logger.Debug("Calling content provider", "name", name)
		defer func(start time.Time) {
			p.logger.Debug("Called content provider", "name", name, "took", time.Since(start))
		}(time.Now())
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
			ContentId:   params.ContentID,
		}, p.callOptions()...)
		if err != nil {
			return nil, hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to generate content",
				Detail:   err.Error(),
			}}
		}
		result := decodeContentResult(res.Result)
		diags := decodeDiagnosticList(res.Diagnostics)
		return result, diags
	}
}

func (p *grpcPlugin) clientDataFunc(name string, client PluginServiceClient) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (plugin.Data, hcl.Diagnostics) {
		p.logger.Debug("Calling data source", "name", name)
		defer func(start time.Time) {
			p.logger.Debug("Called data source", "name", name, "took", time.Since(start))
		}(time.Now())
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
		}, p.callOptions()...)
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

func (p *grpcPlugin) clientPublishFunc(name string, client PluginServiceClient) plugin.PublishFunc {
	return func(ctx context.Context, params *plugin.PublishParams) hcl.Diagnostics {
		p.logger.Debug("Calling publisher", "name", name)
		defer func(start time.Time) {
			p.logger.Debug("Called publisher", "name", name, "took", time.Since(start))
		}(time.Now())
		if params == nil {
			return hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Nil params",
				Detail:   "Nil params",
			}}
		}
		argsEncoded, err := encodeCtyValue(params.Args)
		if err != nil {
			return hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to encode args",
				Detail:   err.Error(),
			}}
		}
		cfgEncoded, err := encodeCtyValue(params.Config)
		if err != nil {
			return hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to encode config",
				Detail:   err.Error(),
			}}
		}
		datactx := encodeMapData(params.DataContext)
		format := encodeOutputFormat(params.Format)
		res, err := client.Publish(ctx, &PublishRequest{
			Publisher:   name,
			Config:      cfgEncoded,
			Args:        argsEncoded,
			DataContext: datactx,
			Format:      format,
		}, p.callOptions()...)
		if err != nil {
			return hcl.Diagnostics{{
				Severity: hcl.DiagError,
				Summary:  "Failed to publish",
				Detail:   err.Error(),
			}}
		}
		return decodeDiagnosticList(res.Diagnostics)
	}
}
