package pluginapiv1

import (
	context "context"
	"log/slog"

	goplugin "github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	"github.com/blackstork-io/fabric/plugin"
	astv1 "github.com/blackstork-io/fabric/plugin/ast/v1"
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
	schema, diags := encodeSchema(srv.schema)
	if diags.HasErrors() {
		return nil, status.Errorf(codes.Internal, "failed to encode schema: %v", diags)
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
	cfg, err := decodeBlock(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeBlock(req.GetArgs())
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
	cfg, err := decodeBlock(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeBlock(req.GetArgs())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode args: %v", err)
	}
	datactx := decodeMapData(req.GetDataContext().GetValue())
	result, diags := srv.schema.ProvideContent(ctx, provider, &plugin.ProvideContentParams{
		Config:      cfg,
		Args:        args,
		DataContext: datactx,
	})
	return &ProvideContentResponse{
		Result:      astv1.EncodeNode(result.AsNode()),
		Diagnostics: encodeDiagnosticList(diags),
	}, nil
}

func (srv *grpcServer) Publish(ctx context.Context, req *PublishRequest) (*PublishResponse, error) {
	slog.DebugContext(ctx, "Publish")
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "Publishing failed", "panic", r)
			panic(r)
		} else {
			slog.DebugContext(ctx, "Publish done")
		}
	}()
	publisher := req.GetPublisher()
	cfg, err := decodeBlock(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeBlock(req.GetArgs())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode args: %v", err)
	}
	datactx := decodeMapData(req.GetDataContext().GetValue())
	diags := srv.schema.Publish(ctx, publisher, &plugin.PublishParams{
		Config:       cfg,
		Args:         args,
		DataContext:  datactx,
		DocumentName: req.GetDocumentName(),
		Document:     astv1.DecodeNode(req.GetDocument()),
	})
	return &PublishResponse{
		Diagnostics: encodeDiagnosticList(diags),
	}, nil
}

func (srv *grpcServer) PublisherInfo(ctx context.Context, req *PublisherInfoRequest) (*PublisherInfoResponse, error) {
	slog.DebugContext(ctx, "Publisher info")
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "Publisher info failed", "panic", r)
			panic(r)
		} else {
			slog.DebugContext(ctx, "Publisher info done")
		}
	}()
	publisher := req.GetPublisher()
	cfg, err := decodeBlock(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeBlock(req.GetArgs())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode args: %v", err)
	}

	info, diags := srv.schema.PublisherInfo(ctx, publisher, &plugin.PublisherInfoParams{
		Config: cfg,
		Args:   args,
	})
	return &PublisherInfoResponse{
		PublisherInfo: encodePublisherInfo(info),
		Diagnostics:   encodeDiagnosticList(diags),
	}, nil
}

func (srv *grpcServer) RenderNode(ctx context.Context, req *RenderNodeRequest) (*RenderNodeResponse, error) {
	slog.DebugContext(ctx, "RenderNode")
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "RenderNode done", "panic", r)
			panic(r)
		} else {
			slog.DebugContext(ctx, "RenderNode done")
		}
	}()

	result, diags := srv.schema.RenderNode(ctx, decodeRenderNodeParams(req))

	return &RenderNodeResponse{
		SubtreeReplacement: astv1.EncodeNode(result),
		Diagnostics:        encodeDiagnosticList(diags),
	}, nil
}
