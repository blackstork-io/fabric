package pluginapiv1

import (
	context "context"

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
	})
}

type grpcServer struct {
	schema *plugin.Schema
	UnimplementedPluginServiceServer
}

func (srv *grpcServer) GetSchema(ctx context.Context, req *GetSchemaRequest) (*GetSchemaResponse, error) {
	schema, err := encodeSchema(srv.schema)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to encode schema: %v", err)
	}
	return &GetSchemaResponse{Schema: schema}, nil
}

func (srv *grpcServer) RetrieveData(ctx context.Context, req *RetrieveDataRequest) (*RetrieveDataResponse, error) {
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
	provider := req.GetProvider()
	cfg, err := decodeCtyValue(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeCtyValue(req.GetArgs())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode args: %v", err)
	}
	datactx := decodeMapData(req.GetDataContext())
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
	publisher := req.GetPublisher()
	cfg, err := decodeCtyValue(req.GetConfig())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode config: %v", err)
	}
	args, err := decodeCtyValue(req.GetArgs())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to decode args: %v", err)
	}
	datactx := decodeMapData(req.GetDataContext())
	format := decodeOutputFormat(req.GetFormat())
	diags := srv.schema.Publish(ctx, publisher, &plugin.PublishParams{
		Config:      cfg,
		Args:        args,
		DataContext: datactx,
		Format:      format,
	})
	return &PublishResponse{
		Diagnostics: encodeDiagnosticList(diags),
	}, nil
}
