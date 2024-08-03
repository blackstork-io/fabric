package pluginapiv1

import (
	context "context"
	"log/slog"

	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"github.com/blackstork-io/fabric/plugin"
)

func Serve(schema *plugin.Schema) {
	goplugin.Serve(&goplugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins: map[string]goplugin.Plugin{
			schema.Name: &grpcPlugin{schema: schema},
		},
		GRPCServer: func(opts []grpc.ServerOption) *grpc.Server {
			opts = append(opts, grpc.MaxRecvMsgSize(defaultMsgSize))
			return grpc.NewServer(opts...)
		},
		Logger: loggerForGoplugin(),
	})
}

type grpcServer struct {
	schema *plugin.Schema
	UnimplementedPluginServiceServer
}

func (srv *grpcServer) GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error) {
	slog.DebugContext(ctx, "GetSchema")
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "GetSchema done", "panic", r)
			panic(r)
		} else {
			slog.DebugContext(ctx, "GetSchema done")
		}
	}()
	schema, err := encodeSchema(srv.schema)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to encode schema: %v", err)
	}
	return &GetSchemaResponse{Schema: schema}, nil
}

func (srv *grpcServer) RetrieveData(ctx context.Context, req *RetrieveDataRequest) (*RetrieveDataResponse, error) {
	slog.DebugContext(ctx, "RetrieveData")
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "RetrieveData done", "panic", r)
			panic(r)
		} else {
			slog.DebugContext(ctx, "RetrieveData done")
		}
	}()
	source := req.GetSource()
	cfg, err := decodeCtyValue(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeCtyValue(req.GetArgs())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode args: %v", err)
	}
	data, diags := srv.schema.RetrieveData(ctx, source, &plugin.RetrieveDataParams{
		Config: cfg,
		Args:   args,
	})
	return &RetrieveDataResponse{
		Data:        encodeData(data),
		Diagnostics: encodeDiagnosticList(diags),
	}, nil
}

func (srv *grpcServer) ProvideContent(ctx context.Context, req *ProvideContentRequest) (*ProvideContentResponse, error) {
	slog.DebugContext(ctx, "ProvideContent")
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "ProvideContent done", "panic", r)
			panic(r)
		} else {
			slog.DebugContext(ctx, "ProvideContent done")
		}
	}()
	provider := req.GetProvider()
	cfg, err := decodeCtyValue(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeCtyValue(req.GetArgs())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode args: %v", err)
	}
	datactx := decodeMapData(req.GetDataContext().GetValue())
	result, diags := srv.schema.ProvideContent(ctx, provider, &plugin.ProvideContentParams{
		Config:      cfg,
		Args:        args,
		DataContext: datactx,
		ContentID:   req.GetContentId(),
	})
	return &ProvideContentResponse{
		Result:      encodeContentResult(result),
		Diagnostics: encodeDiagnosticList(diags),
	}, nil
}

func (srv *grpcServer) Publish(ctx context.Context, req *PublishRequest) (*PublishResponse, error) {
	slog.DebugContext(ctx, "Publish")
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "Publish done", "panic", r)
			panic(r)
		} else {
			slog.DebugContext(ctx, "Publish done")
		}
	}()
	publisher := req.GetPublisher()
	cfg, err := decodeCtyValue(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeCtyValue(req.GetArgs())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode args: %v", err)
	}
	datactx := decodeMapData(req.GetDataContext().GetValue())
	format := decodeOutputFormat(req.GetFormat())
	diags := srv.schema.Publish(ctx, publisher, &plugin.PublishParams{
		Config:       cfg,
		Args:         args,
		DataContext:  datactx,
		Format:       format,
		DocumentName: req.GetDocumentName(),
	})
	return &PublishResponse{
		Diagnostics: encodeDiagnosticList(diags),
	}, nil
}
