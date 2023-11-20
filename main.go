package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/igumus/task-sniffer/config"
	"github.com/igumus/task-sniffer/project"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// todo: something
func main() {
	folder := flag.String("path", ".", "folder to scan")
	branch := flag.String("branch", "main", "folder to scan")
	debug := flag.Bool("debug", false, "sets log level to debug")
	flag.Parse()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	path, err := filepath.Abs(*folder)
	if err != nil {
		log.Fatal().Err(err).Msg("failed")
	} else {
		log.Logger = log.Logger.With().Str("project", path).Logger()
	}

	cfg, err := config.Load(path, *branch)
	if err != nil {
		log.Fatal().Err(err).Msg("failed")
	}

	project := project.New()
	project.ListFiles(context.Background(), cfg)
}
