package config

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
}

func Load() Config {
	return default_configuration
}

type config struct {
	keywords   []Pattern
	exclusions []Pattern
}

func (c *config) Keywords() []Pattern {
	return c.keywords
}

func (c *config) Exclusions() []Pattern {
	return c.exclusions
}
