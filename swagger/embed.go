package swagger

import (
	"embed"
	"io/fs"
)

// swaggerUI built files generating with `make generate` or `make download-swagger`.
//go:embed swagger-ui
var swaggerUI embed.FS

const swaggerUISubPath = "swagger-ui"

//go:embed sink/v1/sink.swagger.json
var sinkSwaggerJSON []byte

// GetSwaggerUI returns a file system in which the Swagger UI is embedded.
func GetSwaggerUI() fs.FS {
	swaggerFS, err := fs.Sub(swaggerUI, swaggerUISubPath)
	if err != nil {
		panic(err) // the application won't compile without go:embed folder
	}
	return swaggerFS
}

// GetSinkSwaggerJSON returns swagger.json byte slice of Sink API service.
func GetSinkSwaggerJSON() []byte {
	return sinkSwaggerJSON
}
