package project

import (
	"bufio"
	"context"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/igumus/task-sniffer/config"
	"github.com/igumus/task-sniffer/reporter"
	"github.com/rs/zerolog/log"
)

func listFiles(cfg config.Config) ([]string, error) {
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

func processFile(cfg config.Config, src string, modify bool, ch chan<- reporter.Task) {
	var (
		err     error
		inFile  *os.File
		outFile *os.File
	)

	inFile, err = os.OpenFile(src, os.O_RDWR, 0777)
	if err != nil {
		log.Error().Str("file", src).Err(err).Msg("reading file failed")
		return
	}
	defer inFile.Close()

	outFile, err = os.OpenFile(src, os.O_RDWR, 0777)
	if err != nil {
		log.Error().Str("file", src).Err(err).Msg("reading file failed")
		return
	}
	defer outFile.Close()

	fileScanner := bufio.NewScanner(inFile)
	fileScanner.Split(bufio.ScanLines)
	found := false
	rawLine := ""
	line := ""
	name := path.Base(src)
	log.Debug().Str("file", name).Msg("Processing")
	for index := 1; fileScanner.Scan(); index++ {
		found = false
		rawLine = fileScanner.Text()
		line = strings.TrimSpace(rawLine)
		if len(line) > 0 && (strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#")) {
			for _, keyword := range cfg.Keywords() {
				groups := keyword.Match(line)
				if len(groups) > 0 {
					ch <- reporter.Task{
						Row:     index,
						Keyword: keyword.Name(),
						File:    name,
						Desc:    groups[3],
					}
					if modify {
						found = true
						break
					}
				}
			}
		}

		if found {
			if !modify {
				outFile.WriteString(rawLine)
				outFile.WriteString("\n")
			}
		} else {
			outFile.WriteString(rawLine)
			outFile.WriteString("\n")
		}
	}
}

// todo: implement details
func Report(ctx context.Context, cfg config.Config, reportFunc reporter.ReportFunc) {
	files, err := listFiles(cfg)
	if err != nil {
		log.Err(err).Msg("listing files failed")
		return
	}
	tasks := make(chan reporter.Task)
	go func() {
		for _, file := range files {
			processFile(cfg, file, cfg.Modify(), tasks)
		}
		close(tasks)
	}()
	reportFunc(ctx, tasks)
}
