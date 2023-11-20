package config

import (
	"path"
)

var (
	keyword_todo     Pattern = newKeyword("todo")
	keyword_fixme    Pattern = newKeyword("fixme")
	default_keywords         = []Pattern{keyword_todo, keyword_fixme}

	// default exclusions
	exclusion_gitignore Pattern = newExclusion(".gitignore")
	exclusion_makefile  Pattern = newExclusion("Makefile")
	exclusion_markdown  Pattern = newExclusion("*.md")
	default_exclusions          = []Pattern{exclusion_gitignore, exclusion_markdown, exclusion_makefile}

	// default configuration
	default_configuration = &config{
		keywords:   default_keywords,
		exclusions: default_exclusions,
	}
)

type ConfigOption func(*config)

type Config interface {
	Keywords() []Pattern
	Exclusions() []Pattern
	Exclude(string)
	Name() string
	Path() string
}

func Load(location, branch string) (Config, error) {
	ret := default_configuration
	ret.path = location
	var err error = nil
	if err = checkRepository(ret.path); err != nil {
		return nil, err
	}
	if ret.url, err = readRepositoryURL(ret.path, branch); err != nil {
		return nil, err
	}

	exclusions := readRepositoryExclusions(ret.path)
	if len(exclusions) > 0 {
		for _, exclusion := range exclusions {
			ret.exclusions = append(ret.exclusions, exclusion)
		}
	}

	return ret, nil
}

type config struct {
	path       string
	url        string
	keywords   []Pattern
	exclusions []Pattern
}

func (c *config) Keywords() []Pattern {
	return c.keywords
}

func (c *config) Exclude(name string) {
	c.exclusions = append(c.exclusions, newExclusion(name))
}

func (c *config) Path() string {
	return c.path
}

func (c *config) Name() string {
	return path.Base(c.path)
}

func (c *config) Exclusions() []Pattern {
	return c.exclusions
}
