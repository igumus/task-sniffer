package config

import (
	"regexp"
	"strings"
)

type patternKind uint8

const (
	unknown   patternKind = iota
	keyword   patternKind = iota
	exclusion patternKind = iota
)

type Pattern struct {
	name     string
	compiled *regexp.Regexp
	kind     patternKind
}

func newPattern(kind patternKind, name string) *Pattern {
	ret := &Pattern{
		name: strings.TrimSpace(name),
		kind: kind,
	}
	ret.compile()
	return ret
}

func newKeyword(name string) *Pattern {
	return newPattern(keyword, name)
}

func newExclusion(name string) *Pattern {
	return newPattern(exclusion, name)
}

func (k *Pattern) Name() string {
	return k.name
}

func (k *Pattern) Match(str string) []string {
	return k.compiled.FindStringSubmatch(str)
}

func (k *Pattern) revert(ch rune) string {
	if ch >= 65 && ch < 97 {
		return strings.ToLower(string(ch))
	}
	return strings.ToUpper(string(ch))
}

func (k *Pattern) buildRegexString() string {
	builder := strings.Builder{}
	builder.WriteString("^(.*")
	for _, ch := range k.name {
		builder.WriteString("[")
		builder.WriteString(string(ch))
		builder.WriteString(k.revert(ch))
		builder.WriteString("]+")
	}

	builder.WriteString(")([: ]*)*(.*)$")

	return builder.String()
}

func (k *Pattern) compile() {
	k.compiled = regexp.MustCompile(k.buildRegexString())
}
