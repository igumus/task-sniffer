package project

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/igumus/task-sniffer/config"
	"github.com/rs/zerolog/log"
)

type Project interface {
	ListFiles(context.Context, config.Config)
}

type project struct{}

func New() Project {
	return &project{}
}

func (p *project) listFiles(cfg config.Config) ([]string, error) {
	ret := make([]string, 0)
	return ret, filepath.Walk(cfg.Path(), func(path string, info fs.FileInfo, err error) error {
		if cfg.Path() == path {
			return nil
		}
		file := info.Name()
		isDir := info.IsDir()
		isHidden := strings.HasPrefix(file, ".")
		if isHidden && isDir {
			if isDir {
				log.Debug().Str("folder", file).Msg("ignore hidden folders")
				return filepath.SkipDir
			} else {
				log.Debug().Str("file", file).Msg("ignore hidden files")
				return nil
			}
		}

		include := true
		for _, exclusion := range cfg.Exclusions() {
			include = true
			if len(exclusion.Match(file)) > 0 {
				if isDir {
					log.Debug().Str("folder", file).Str("filter", exclusion.Name()).Msg("exclude folder")
					return filepath.SkipDir
				} else {
					log.Debug().Str("file", file).Str("filter", exclusion.Name()).Msg("exclude file")
					include = false
					break
				}
			}
		}
		if include && !isDir {
			ret = append(ret, path)
		}
		return nil
	})
}

func (p *project) ListFiles(ctx context.Context, cfg config.Config) {
	files, err := p.listFiles(cfg)
	if err != nil {
		log.Err(err).Msg("listing files failed")
		return
	}
	for _, file := range files {
		log.Info().Str("file", file).Msg("Listing files")
	}
}
