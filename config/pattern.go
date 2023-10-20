package config

import (
	"regexp"
	"strings"
)

type Pattern interface {
	Name() string
	Match(string) []string
	String() string
}

type patternKind uint8

const (
	unknown   patternKind = iota
	keyword   patternKind = iota
	exclusion patternKind = iota
)

type pattern struct {
	name     string
	compiled *regexp.Regexp
	kind     patternKind
}

func newPattern(kind patternKind, name string) *pattern {
	ret := &pattern{
		name: strings.TrimSpace(name),
		kind: kind,
	}
	ret.compile()
	return ret
}

func newKeyword(name string) *pattern {
	return newPattern(keyword, name)
}

func newExclusion(name string) *pattern {
	return newPattern(exclusion, name)
}

func (k *pattern) revert(ch rune) string {
	if ch >= 65 && ch < 97 {
		return strings.ToLower(string(ch))
	}
	return strings.ToUpper(string(ch))
}

func (k *pattern) buildRegexString() string {
	builder := strings.Builder{}
	builder.WriteString("^")
	switch k.kind {
	case keyword:
		builder.WriteString("(.*")
		for _, ch := range k.name {
			builder.WriteString("[")
			builder.WriteString(string(ch))
			builder.WriteString(k.revert(ch))
			builder.WriteString("]+")
		}
		builder.WriteString(")([: ]*)*(.*)")
		break
	case exclusion:
		value := k.name
		value = strings.ReplaceAll(value, "?", "")
		value = strings.ReplaceAll(value, ".", "\\.")
		value = strings.ReplaceAll(value, "*", "(.*)")
		builder.WriteString(value)
		break
	default:
		panic("not reachable")
	}
	builder.WriteString("$")

	return builder.String()
}

func (k *pattern) compile() {
	regex := k.buildRegexString()
	k.compiled = regexp.MustCompile(regex)
}

func (k *pattern) String() string {
	switch k.kind {
	case keyword:
		return strings.Title(k.name)
	default:
		return k.name
	}
}

func (k *pattern) Name() string {
	return k.name
}

func (k *pattern) Match(str string) []string {
	return k.compiled.FindStringSubmatch(str)
}
