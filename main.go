package main

import (
	"context"
	"embed"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/owenthereal/jqplay/jq"
)

//go:embed all:public
var PublicFS embed.FS

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	conf, err := Load()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading config")
	}

	log.Info().
		Str("version", jq.Version).
		Str("path", jq.Path).
		Msg("initialized jq")

	log.Info().
		Str("host", conf.Host).
		Str("port", conf.Port).
		Msg("starting server")

	if err := newServer(context.Background(), conf); err != nil {
		log.Fatal().Err(err).Msg("error starting server")
	}
}
