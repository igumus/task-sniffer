package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	"github.com/igumus/task-sniffer/config"
	"github.com/igumus/task-sniffer/project"
	"github.com/igumus/task-sniffer/reporter"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// todo: something
func main() {
	folder := flag.String("path", ".", "folder to scan")
	branch := flag.String("branch", "main", "folder to scan")
	debug := flag.Bool("debug", false, "sets log level to debug")
	reporterName := flag.String("reporter", "console", "[console]")
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

	var (
		report    reporter.ReportFunc
		reportErr error
	)

	if *reporterName == "console" {
		report, reportErr = reporter.DefaultReporter(cfg)
	} else {
		log.Warn().Str("reporter", *reporterName).Msg("unknown reporter, fallback to default repoter")
		report, reportErr = reporter.DefaultReporter(cfg)
	}
	if reportErr != nil {
		log.Fatal().Err(reportErr).Msg("failed to create reporter")
	}

	project.Report(context.Background(), cfg, report)
}
