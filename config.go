package main

import (
	"github.com/joeshaw/envdecode"

	"github.com/owenthereal/jqplay/jq"
)

type Config struct {
	Host      string `env:"HOST,default=0.0.0.0"`
	Port      string `env:"PORT,default=8080"`
	GinMode   string `env:"GIN_MODE,default=debug"`
	AssetHost string `env:"ASSET_HOST"`
	JQVer     string
}

func (c *Config) IsProd() bool {
	return c.GinMode == "release"
}

func Load() (*Config, error) {
	conf := &Config{}
	err := envdecode.Decode(conf)
	if err != nil {
		return nil, err
	}

	conf.JQVer = jq.Version

	return conf, nil
}
