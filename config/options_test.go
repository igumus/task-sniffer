package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCheckPatternExists(t *testing.T, name string, items []Pattern) {
	for _, item := range items {
		if item.Name() == name {
			return
		}
	}
	assert.Failf(t, "Pattern Not Exists", "`%s` not found in items", name)
}

func TestConfigCreation(t *testing.T) {
	const defaultExclusionCount = 3
	exclusions := default_configuration.Exclusions()
	require.Equal(t, len(exclusions), defaultExclusionCount)
	testCheckPatternExists(t, "Makefile", exclusions)
	testCheckPatternExists(t, "*.md", exclusions)
	testCheckPatternExists(t, ".gitignore", exclusions)

	const defaultKeywordCount = 2
	keywords := default_configuration.Keywords()
	require.Equal(t, len(keywords), defaultKeywordCount)
	testCheckPatternExists(t, "todo", keywords)
	testCheckPatternExists(t, "fixme", keywords)

}
