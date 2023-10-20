package main

import (
	"flag"
	"os"

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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	cfg, err := config.Load(*folder, *branch)
	if err != nil {
		log.Fatal().Err(err).Msg("failed")
	}

	project := project.New(cfg)
	if err := project.Tasks(); err != nil {
		log.Fatal().Err(err)
	}
}
