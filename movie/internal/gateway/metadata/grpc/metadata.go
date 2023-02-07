package grpc

import (
	"context"

	"github.com/luispinto23/movieexample/gen"
	"github.com/luispinto23/movieexample/internal/grpcutil"
	"github.com/luispinto23/movieexample/metadata/pkg/model"
	"github.com/luispinto23/movieexample/pkg/discovery"
)

// Gateway defines a movie metadata gRPC gateway.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a movie metadata service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry: registry}
}

// Get returns movie metadata by a movie id.
func (g *Gateway) Get(ctx context.Context, id string) (*model.Metadata, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "metadata", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	metadataRequest := gen.GetMetadataRequest{MovieId: id}

	client := gen.NewMetadataServiceClient(conn)
	resp, err := client.GetMetadata(ctx, &metadataRequest)
	if err != nil {
		return nil, err
	}

	return model.MetadataFromProto(resp.Metadata), nil
}
