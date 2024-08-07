package pluginapiv1

import (
	context "context"
	"fmt"
	"log/slog"
	"time"

	goplugin "github.com/hashicorp/go-plugin"
	"github.com/hashicorp/hcl/v2"
	grpc "google.golang.org/grpc"

	"github.com/blackstork-io/fabric/pkg/diagnostics"
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
	return func(ctx context.Context, params *plugin.ProvideContentParams) (result *plugin.ContentResult, diags diagnostics.Diag) {
		p.logger.DebugContext(ctx, "Calling content provider", "name", name)
		defer func(start time.Time) {
			p.logger.DebugContext(ctx, "Called content provider", "name", name, "took", time.Since(start))
		}(time.Now())
		if params == nil {
			diags.Add("Content provider error", "Nil params")
			return
		}
		cfgEncoded, diag := encodeCtyValue(params.Config)
		diags.Extend(diag)
		argsEncoded, diag := encodeCtyValue(params.Args)
		diags.Extend(diag)
		if diags.HasErrors() {
			return
		}
		res, err := client.ProvideContent(ctx, &ProvideContentRequest{
			Provider:    name,
			Config:      cfgEncoded,
			Args:        argsEncoded,
			DataContext: encodeMapData(params.DataContext),
			ContentId:   params.ContentID,
		}, p.callOptions()...)
		if diags.AppendErr(err, "Failed to generate content") {
			return
		}
		result = decodeContentResult(res.GetResult())
		diags.Extend(decodeDiagnosticList(res.GetDiagnostics()))
		return result, diags
	}
}

func (p *grpcPlugin) clientDataFunc(name string, client PluginServiceClient) plugin.RetrieveDataFunc {
	return func(ctx context.Context, params *plugin.RetrieveDataParams) (data plugin.Data, diags diagnostics.Diag) {
		p.logger.DebugContext(ctx, "Calling data source", "name", name)
		defer func(start time.Time) {
			p.logger.DebugContext(ctx, "Called data source", "name", name, "took", time.Since(start))
		}(time.Now())
		if params == nil {
			diags.Add("Data source error", "Nil params")
			return
		}
		cfgEncoded, diag := encodeCtyValue(params.Config)
		diags.Extend(diag)
		argsEncoded, diag := encodeCtyValue(params.Args)
		diags.Extend(diag)

		res, err := client.RetrieveData(ctx, &RetrieveDataRequest{
			Source: name,
			Config: cfgEncoded,
			Args:   argsEncoded,
		}, p.callOptions()...)
		if diags.AppendErr(err, "Failed to fetch data") {
			return
		}
		data = decodeData(res.GetData())
		diags.Extend(decodeDiagnosticList(res.GetDiagnostics()))
		return
	}
}

func (p *grpcPlugin) clientPublishFunc(name string, client PluginServiceClient) plugin.PublishFunc {
	return func(ctx context.Context, params *plugin.PublishParams) (diags diagnostics.Diag) {
		p.logger.DebugContext(ctx, "Calling publisher", "name", name)
		defer func(start time.Time) {
			p.logger.DebugContext(ctx, "Called publisher", "name", name, "took", time.Since(start))
		}(time.Now())
		if params == nil {
			diags.Add("Publisher error", "Nil params")
			return
		}
		argsEncoded, diag := encodeCtyValue(params.Args)
		diags.Extend(diag)
		cfgEncoded, diag := encodeCtyValue(params.Config)
		diags.Extend(diag)
		datactx := encodeMapData(params.DataContext)
		format := encodeOutputFormat(params.Format)
		res, err := client.Publish(ctx, &PublishRequest{
			Publisher:    name,
			Config:       cfgEncoded,
			Args:         argsEncoded,
			DataContext:  datactx,
			Format:       format,
			DocumentName: params.DocumentName,
		}, p.callOptions()...)

		if diags.AppendErr(err, "Failed to publish") {
			return diagnostics.Diag{{
				Severity: hcl.DiagError,
				Summary:  "Failed to publish",
				Detail:   err.Error(),
			}}
		}
		diags.Extend(decodeDiagnosticList(res.GetDiagnostics()))
		return
	}
}
