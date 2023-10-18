package config

import (
	"regexp"
	"strings"
)

func revert(ch rune) string {
	if ch >= 65 && ch < 97 {
		return strings.ToLower(string(ch))
	}
	return strings.ToUpper(string(ch))
}

type Keyword struct {
	name     string
	compiled *regexp.Regexp
}

func NewKeyword(name string) *Keyword {
	ret := &Keyword{
		name: strings.TrimSpace(name),
	}
	ret.compile()
	return ret
}

func (k *Keyword) Name() string {
	return k.name
}

func (k *Keyword) Match(str string) []string {
	return k.compiled.FindStringSubmatch(str)
}

func (k *Keyword) buildRegexString() string {
	builder := strings.Builder{}
	builder.WriteString("^(.*")
	for _, ch := range k.name {
		builder.WriteString("[")
		builder.WriteString(string(ch))
		builder.WriteString(revert(ch))
		builder.WriteString("]+")
	}

	builder.WriteString(")([: ]*)*(.*)$")

	return builder.String()
}

func (k *Keyword) compile() {
	k.compiled = regexp.MustCompile(k.buildRegexString())
}
