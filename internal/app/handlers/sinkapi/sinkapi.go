package sinkapi

import (
	"context"
)

// URLShortener interface.
type URLShortener interface {
	Shorten(ctx context.Context, link string) (shortenedLink string, err error)
	Unshort(ctx context.Context, shortened string) (unshortenedLink string, err error)
}

// SinkAPI data structure.
type SinkAPI struct {
	shortener URLShortener
}

// New creates a new instance if SinkAPI.
func New(shortener URLShortener) SinkAPI {
	return SinkAPI{
		shortener: shortener,
	}
}
