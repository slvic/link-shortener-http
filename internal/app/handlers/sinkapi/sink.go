package sinkapi

import (
	"context"

	gw "github.com/go-sink/sink/pkg/sink/v1"
)

// Sink accepts a link and shortens it, returns shortened link.
func (s SinkAPI) Sink(ctx context.Context, req *gw.SinkRequest) (*gw.SinkResponse, error) {
	shortened, err := s.shortener.Shorten(ctx, req.Url)
	if err != nil {
		return nil, err
	}

	return &gw.SinkResponse{Url: shortened}, nil
}
