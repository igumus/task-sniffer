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

	var keyword *Keyword
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			keyword = NewKeyword(tc.name)
			require.Equal(t, keyword.buildRegexString(), tc.expected)
		})
	}
}

func TestKeywordMatch(t *testing.T) {
	testcases := []struct {
		name   string
		inputs []struct {
			value       string
			result      string
			shouldMatch bool
		}
	}{
		{
			name: "todo",
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
	}

	var keyword *Keyword
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			keyword = NewKeyword(tc.name)
			for _, input := range tc.inputs {
				groups := keyword.Match(input.value)
				hasmatch := len(groups) > 0
				require.Equal(t, input.shouldMatch, hasmatch)
				if hasmatch {
					require.Equal(t, input.result, groups[3])
				}
			}
		})
	}
}
