package server

import (
	"mime"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/go-sink/sink/swagger"
)

const swaggerUIPrefix = "/docs/"
const swaggerJSONPath = "/docs/swagger.json"

func serveSwaggerUI(mux *http.ServeMux) error {
	if err := mime.AddExtensionType(".svg", "image/svg+xml"); err != nil {
		return err
	}

	// Expose files on <host>/docs/
	mux.HandleFunc(swaggerJSONPath, func(w http.ResponseWriter, _ *http.Request) {
		if _, err := w.Write(swagger.GetSinkSwaggerJSON()); err != nil {
			log.Err(err).Msg("error writing swagger.json file: %w")
		}
	})

	mux.Handle(swaggerUIPrefix, http.StripPrefix(swaggerUIPrefix, http.FileServer(http.FS(swagger.GetSwaggerUI()))))

	return nil
}
