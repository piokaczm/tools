package topics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanerPipelines(t *testing.T) {
	tcs := []struct {
		name     string
		filters  []Filter
		document []string
		expected []string
		assert   func(*testing.T, []string, []string, error)
	}{
		{
			name:     "with only nouns filter",
			filters:  []Filter{OnlyWithNouns},
			document: []string{"The mother of my friend is running like crazy!"},
			expected: []string{"mother friend crazy"},
			assert: func(t *testing.T, exp []string, res []string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, exp, res)
			},
		},
		{
			name:     "with length limit filter",
			filters:  []Filter{NotShorterThan(3)},
			document: []string{"The mother of my friend is running like crazy!"},
			expected: []string{"mother friend running like crazy"},
			assert: func(t *testing.T, exp []string, res []string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, exp, res)
			},
		},
		{
			name:     "with stemming filter",
			filters:  []Filter{WithStemming},
			document: []string{"The mother of my friend is running like crazy!"},
			expected: []string{"the mother of my friend is run like crazi !"},
			assert: func(t *testing.T, exp []string, res []string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, exp, res)
			},
		},
		{
			name:     "with lemmatizing filter",
			filters:  []Filter{WithLemmatizing},
			document: []string{"The mother of my friend is running like crazy!"},
			expected: []string{"The mother of my friend be run like crazy !"},
			assert: func(t *testing.T, exp []string, res []string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, exp, res)
			},
		},
		{
			name:     "with multiple filters",
			filters:  []Filter{Downcase, WithLemmatizing, WithStemming},
			document: []string{"The MotHer of my fRiend is runNing like crazy!"},
			expected: []string{"the mother of my friend be run like crazi !"},
			assert: func(t *testing.T, exp []string, res []string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, exp, res)
			},
		},
		{
			name:     "with names filters",
			filters:  []Filter{Downcase, WithLemmatizing, WithStemming},
			document: []string{"Katowice is a pretty city. I'd love to play with MyApp."},
			expected: []string{"katowic be a pretti citi . i'd love to play with myapp ."},
			assert: func(t *testing.T, exp []string, res []string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, exp, res)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCleaner()
			c.BuildPipeline(tc.filters...)

			clean, err := c.Clean(tc.document)
			tc.assert(t, tc.expected, clean, err)
		})
	}
}
