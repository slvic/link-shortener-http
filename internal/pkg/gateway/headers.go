package gateway

import (
	"context"
	"net/http"
	"strconv"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"
)

// ReplaceHeaders of a grpc request with that of http's
func ReplaceHeaders() func(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	return func(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
		md, ok := runtime.ServerMetadataFromContext(ctx)
		if !ok {
			return nil
		}

		if vals := md.HeaderMD.Get("x-http-location"); len(vals) > 0 {
			delete(md.HeaderMD, "x-http-location")
			delete(w.Header(), "Grpc-Metadata-x-http-location")

			w.Header().Set("location", vals[0])
		}

		// set http status code
		if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
			code, err := strconv.Atoi(vals[0])
			if err != nil {
				return err
			}

			// delete the headers to not expose any grpc-metadata in http response
			delete(md.HeaderMD, "x-http-code")
			delete(w.Header(), "Grpc-Metadata-x-http-code")

			w.WriteHeader(code)
		}

		return nil
	}
}
