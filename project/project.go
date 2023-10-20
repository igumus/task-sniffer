package project

import (
	"bufio"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/igumus/task-sniffer/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Project interface {
	Tasks() error
}

type project struct {
	logger zerolog.Logger
	cfg    config.Config
}

func New(cfg config.Config) Project {
	ret := &project{
		cfg: cfg,
	}
	ret.logger = log.Logger.With().Str("project", cfg.Name()).Logger()
	return ret
}

func (p *project) listFiles() ([]string, error) {
	ret := make([]string, 0)
	return ret, filepath.Walk(p.cfg.Path(), func(path string, info fs.FileInfo, err error) error {
		if p.cfg.Path() != path {
			file := info.Name()
			isDir := info.IsDir()
			if strings.HasPrefix(file, ".") {
				if isDir {
					return filepath.SkipDir
				}
			} else {
				include := true
				for _, exclusion := range p.cfg.Exclusions() {
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
			}
		}
		return nil
	})
}

func (p *project) process(file string) {
	readFile, err := os.Open(file)
	if err != nil {
		p.logger.Error().Str("file", file).Err(err).Msg("reading file failed")
		return
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	line := ""
	p.logger.Debug().
		Str("file", path.Base(file)).
		Msg("Processing")
	for index := 1; fileScanner.Scan(); index++ {
		line = strings.TrimSpace(fileScanner.Text())
		if len(line) > 0 && strings.HasPrefix(line, "//") {
			for _, keyword := range p.cfg.Keywords() {
				groups := keyword.Match(line)
				if len(groups) > 0 {
					p.logger.Info().
						Int("row", index).
						Str("file", path.Base(file)).
						Str("keyword", keyword.Name()).
						Str("name", groups[3]).
						Msg("Issue")
				}
			}
		}
	}
}

func (p *project) Tasks() error {
	files, err := p.listFiles()
	if err != nil {
		log.Err(err)
		return err
	}
	p.logger.Debug().Int("count", len(files)).Msg("files")
	for _, file := range files {
		p.process(file)
	}
	return nil
}
