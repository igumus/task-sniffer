package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeywordRegexString(t *testing.T) {
	testcases := []struct {
		name     string
		expected string
	}{

		{
			name:     "todo",
			expected: "^(.*[tT]+[oO]+[dD]+[oO]+)([: ]*)*(.*)$",
		},
		{
			name:     "Todo",
			expected: "^(.*[Tt]+[oO]+[dD]+[oO]+)([: ]*)*(.*)$",
		},
		{
			name:     "fixme",
			expected: "^(.*[fF]+[iI]+[xX]+[mM]+[eE]+)([: ]*)*(.*)$",
		},
		{
			name:     "FixMe",
			expected: "^(.*[Ff]+[iI]+[xX]+[Mm]+[eE]+)([: ]*)*(.*)$",
		},
	}

	var keyword *Pattern
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			keyword = newKeyword(tc.name)
			require.Equal(t, keyword.buildRegexString(), tc.expected)
		})
	}
}

func TestKeywordMatch(t *testing.T) {
	Keyword_Todo := newKeyword("todo")
	Keyword_Fixme := newKeyword("fixme")

	testcases := []struct {
		kw     *Pattern
		inputs []struct {
			value       string
			result      string
			shouldMatch bool
		}
	}{
		{
			kw: Keyword_Todo,
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
			kw: Keyword_Fixme,
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
		t.Run(tc.kw.name, func(t *testing.T) {
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
