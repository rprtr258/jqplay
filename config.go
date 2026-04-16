package main

import (
	"cmp"
	"os"

	"github.com/rprtr258/jqplay/jq"
)

type Config struct {
	Host      string
	Port      string
	Env       string
	AssetHost string
	JQVer     string
}

func (c *Config) IsProd() bool {
	return c.Env == "production"
}

func Load() Config {
	return Config{
		Host:      cmp.Or(os.Getenv("HOST"), "0.0.0.0"),
		Port:      cmp.Or(os.Getenv("PORT"), "8080"),
		Env:       cmp.Or(os.Getenv("ENV"), "development"),
		AssetHost: os.Getenv("ASSET_HOST"),
		JQVer:     jq.Version,
	}
}
