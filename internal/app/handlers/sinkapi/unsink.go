package sinkapi

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	gw "github.com/go-sink/sink/pkg/sink/v1"
)

// Unsink accepts a link and unshortens it, returns an unshortened link.
func (s SinkAPI) Unsink(ctx context.Context, req *gw.UnsinkRequest) (*gw.UnsinkResponse, error) {
	unshortened, err := s.shortener.Unshort(ctx, req.Url)
	if err != nil {
		return nil, err
	}

	if err = grpc.SetHeader(ctx, metadata.Pairs("x-http-code", "301", "x-http-location", unshortened)); err != nil {
		return nil, err
	}

	return &gw.UnsinkResponse{Url: unshortened}, nil
}
