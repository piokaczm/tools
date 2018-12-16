package comparator

import (
	"errors"
	"testing"

	"github.com/piokaczm/tools/foureyes/mock"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	tcs := []struct {
		name        string
		topics      []string
		source      Source
		expectedMsg string
		assert      func(*testing.T, string, string, error)
	}{
		{
			name:        "happy path",
			topics:      []string{"topic", "another", "notExistingTopic"},
			source:      mock.NewSource("slack channel #random", [][]string{[]string{"topic", "lameTopic", "another"}}, nil),
			expectedMsg: "Hi there, I've found that someone is talking about topic, another in slack channel #random!",
			assert: func(t *testing.T, expected, msg string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expected, msg)
			},
		},
		{
			name:   "with duplicates within source topics",
			topics: []string{"topic", "another", "notExistingTopic"},
			source: mock.NewSource(
				"slack channel #random",
				[][]string{
					[]string{"topic", "lameTopic", "another"},
					[]string{"topic", "lameTopic", "another"},
				},
				nil),
			expectedMsg: "Hi there, I've found that someone is talking about topic, another in slack channel #random!",
			assert: func(t *testing.T, expected, msg string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expected, msg)
			},
		},
		{
			name:        "with mixed case topics",
			topics:      []string{"ToPiC", "anotherTOPIC", "notExistingTopic"},
			source:      mock.NewSource("slack channel #random", [][]string{[]string{"topiC", "ANOTHERtopic", "lameTOpic"}}, nil),
			expectedMsg: "Hi there, I've found that someone is talking about topiC, ANOTHERtopic in slack channel #random!",
			assert: func(t *testing.T, expected, msg string, err error) {
				assert.Nil(t, err)
				assert.Equal(t, expected, msg)
			},
		},
		{
			name:        "with empty comparator topics",
			topics:      []string{},
			source:      mock.NewSource("slack channel #random", [][]string{[]string{"topiC", "ANOTHERtopic", "lameTOpic"}}, nil),
			expectedMsg: "",
			assert: func(t *testing.T, expected, msg string, err error) {
				assert.EqualError(t, err, "comparator: no matches")
				assert.Equal(t, expected, msg)
			},
		},
		{
			name:        "with empty source topics",
			topics:      []string{"ToPiC", "anotherTOPIC", "notExistingTopic"},
			source:      mock.NewSource("slack channel #random", [][]string{}, nil),
			expectedMsg: "",
			assert: func(t *testing.T, expected, msg string, err error) {
				assert.EqualError(t, err, "comparator: no matches")
				assert.Equal(t, expected, msg)
			},
		},
		{
			name:        "with source error",
			topics:      []string{"ToPiC", "anotherTOPIC", "notExistingTopic"},
			source:      mock.NewSource("slack channel #random", [][]string{}, errors.New("source err")),
			expectedMsg: "",
			assert: func(t *testing.T, expected, msg string, err error) {
				assert.EqualError(t, err, "source err")
				assert.Equal(t, expected, msg)
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			c := New(tc.topics)

			msg, err := c.Match(tc.source)
			tc.assert(t, tc.expectedMsg, msg, err)
		})
	}
}
