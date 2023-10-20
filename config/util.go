package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	ini "gopkg.in/ini.v1"
)

var (
	ErrNoGitRepository = errors.New("Project is not a git repository")
	ErrNoGitRemoteAddr = errors.New("Project does not contain remote url address")
)

func checkRepository(location string) error {
	dotGit := path.Join(location, ".git")
	stat, err := os.Stat(dotGit)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNoGitRepository
		}
		return err
	}
	if !stat.IsDir() {
		return ErrNoGitRepository
	}
	return nil
}

func readRepositoryExclusions(location string) []Pattern {
	ret := make([]Pattern, 0)
	gitIgnore := path.Join(location, ".gitignore")
	file, err := os.Open(gitIgnore)
	if err != nil {
		return ret
	}
	defer file.Close()
	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)
	line := ""
	for fileScanner.Scan() {
		line = strings.ReplaceAll(strings.TrimSpace(fileScanner.Text()), "/", "")
		if len(line) > 0 && !strings.HasPrefix(line, "#") {
			ret = append(ret, newExclusion(line))
		}
	}
	return ret
}

func readRepositoryURL(location, mainBranch string) (string, error) {
	config := path.Join(location, ".git", "config")
	cfg, err := ini.Load(config)
	if err != nil {
		return "", err
	}

	branchSectionName := "branch \"" + mainBranch + "\""
	branchSection := cfg.Section(branchSectionName)
	if branchSection == nil {
		return "", fmt.Errorf("Project does not contain %s branch", mainBranch)
	}

	remote := branchSection.Key("remote")
	if remote == nil {
		return "", fmt.Errorf("Project does not contain remote information in %s", mainBranch)
	}

	sectionName := "remote \"" + remote.String() + "\""
	section := cfg.Section(sectionName)
	if section == nil {
		return "", fmt.Errorf("Project does not contain %s remote section", sectionName)
	}

	url := section.Key("url")
	if url == nil {
		return "", ErrNoGitRemoteAddr
	}
	return url.String(), nil
}
