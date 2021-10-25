package config

import (
	"fmt"

	"github.com/hashicorp/hcl"
)

// Config data structure.
type Config struct {
	Database Database `hcl:"database,block"`
	App      App      `hcl:"app,block"`
}

// App data structure.
type App struct {
	AllowedOrigins []string `hcl:"allowed_origins"`
	HTTPAddr       string   `hcl:"http_addr"`
	GRPCAddr       string   `hcl:"grpc_addr"`
}

// Database data structure.
type Database struct {
	Host     string `hcl:"host"`
	Port     int    `hcl:"port"`
	User     string `hcl:"user"`
	Password string `hcl:"password"`
	Database string `hcl:"database"`
	SSLMode  string `hcl:"sslmode"`
}

// Parse config file.
func Parse(bs []byte) (Config, error) {
	result := &Config{}

	if err := hcl.Unmarshal(bs, result); err != nil {
		return *result, fmt.Errorf("error parsing config: %v", err)
	}

	return *result, nil
}
