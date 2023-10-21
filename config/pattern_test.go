package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPatternRegexStringGeneration(t *testing.T) {
	testcases := []struct {
		name     string
		kind     patternKind
		expected string
	}{

		{
			name:     "todo",
			kind:     keyword,
			expected: "^(.*[tT]+[oO]+[dD]+[oO]+)([: ]*)*(.*)$",
		},
		{
			name:     "Todo",
			kind:     keyword,
			expected: "^(.*[Tt]+[oO]+[dD]+[oO]+)([: ]*)*(.*)$",
		},
		{
			name:     "fixme",
			kind:     keyword,
			expected: "^(.*[fF]+[iI]+[xX]+[mM]+[eE]+)([: ]*)*(.*)$",
		},
		{
			name:     "FixMe",
			kind:     keyword,
			expected: "^(.*[Ff]+[iI]+[xX]+[Mm]+[eE]+)([: ]*)*(.*)$",
		},
		{
			name:     ".gitignore",
			kind:     exclusion,
			expected: "^\\.gitignore$",
		},
	}

	var kw *pattern
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			kw = newPattern(tc.kind, tc.name)
			require.Equal(t, kw.buildRegexString(), tc.expected)
		})
	}
}

func TestExclusionMatch(t *testing.T) {
	testcases := []struct {
		kw          Pattern
		input       string
		result      string
		shouldMatch bool
	}{
		{kw: exclusion_gitignore, input: "", result: "", shouldMatch: false},
		{kw: exclusion_gitignore, input: ".gitignore", result: ".gitignore", shouldMatch: true},
		{kw: exclusion_gitignore, input: ".itignore", result: "", shouldMatch: false},
		{kw: exclusion_makefile, input: "", result: "", shouldMatch: false},
		{kw: exclusion_makefile, input: "Makefile", result: "Makefile", shouldMatch: true},
		{kw: exclusion_makefile, input: "Makeflie", result: "", shouldMatch: false},
		{kw: exclusion_makefile, input: "makefile", result: "", shouldMatch: false},
	}

	for _, tc := range testcases {
		t.Run(tc.kw.Name(), func(t *testing.T) {
			ret := tc.kw.Match(tc.input)
			if tc.shouldMatch {
				require.Greater(t, len(ret), 0)
				require.Equal(t, tc.result, ret[0])
			} else {
				require.Equal(t, len(ret), 0)
			}
		})
	}
}

func TestKeywordMatch(t *testing.T) {
	testcases := []struct {
		kw     Pattern
		inputs []struct {
			value       string
			result      string
			shouldMatch bool
		}
	}{
		{
			kw: keyword_todo,
			inputs: []struct {
				value       string
				result      string
				shouldMatch bool
			}{
				{value: "// todo: item0", shouldMatch: true, result: "item0"},
				{value: "// ToDo: item1", shouldMatch: true, result: "item1"},
				{value: "// ToDooooooo: item2", shouldMatch: true, result: "item2"},
				{value: "// ToDo item3", shouldMatch: true, result: "item3"},
				{value: "// ToDo:", shouldMatch: true, result: ""},
				{value: "//todo:item4", shouldMatch: true, result: "item4"},
				{value: "# ToDo:item5", shouldMatch: true, result: "item5"},
				{value: "-- ToDo:item6", shouldMatch: true, result: "item6"},
				{value: "// ToxDooooooo: item2", shouldMatch: false},
			},
		},
		{
			kw: keyword_fixme,
			inputs: []struct {
				value       string
				result      string
				shouldMatch bool
			}{
				{value: "// fixme: item0", shouldMatch: true, result: "item0"},
				{value: "// FixMe: item1", shouldMatch: true, result: "item1"},
				{value: "// Fixmeeee: item2", shouldMatch: true, result: "item2"},
				{value: "// FixmE item3", shouldMatch: true, result: "item3"},
				{value: "// Fixme:", shouldMatch: true, result: ""},
				{value: "//fixme:item4", shouldMatch: true, result: "item4"},
				{value: "# FixMe:item5", shouldMatch: true, result: "item5"},
				{value: "-- fixMe:item6", shouldMatch: true, result: "item6"},
				{value: "// fixoMe: item2", shouldMatch: false},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.kw.Name(), func(t *testing.T) {
			for _, input := range tc.inputs {
				groups := tc.kw.Match(input.value)
				hasmatch := len(groups) > 0
				require.Equal(t, input.shouldMatch, hasmatch)
				if hasmatch {
					require.Equal(t, input.result, groups[3])
				}
			}
		})
	}
}
