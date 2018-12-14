package comparator

import (
	"fmt"
	"strings"
)

type Matcher struct {
	topics []string
}

type Source interface {
	Name() string
	Topics() [][]string
}

func New(topics []string) *Matcher {
	return &Matcher{topics}
}

func (m *Matcher) Match(source Source) (string, error) {
	intersection, exists := m.findMatches(source.Topics())
	if !exists {
		return "", fmt.Errorf("no matches") // probably better to just return bool, easier for consumers to check...
	}

	return m.buildMsg(intersection, source.Name()), nil
}

func (m *Matcher) buildMsg(in []string, name string) string {
	return fmt.Sprintf(
		"Hi there, I've found that someone is talking about %s in %s!",
		strings.Join(in, ", "),
		name,
	)
}

func (m *Matcher) findMatches(s [][]string) ([]string, bool) {
	var matches []string

	for _, list := range s {
		for _, t := range list {
			if !m.include(t) {
				continue
			}

			matches = append(matches, t)
		}
	}

	return matches, len(matches) > 0
}

func (m *Matcher) include(e string) bool {
	for _, t := range m.topics {
		if t == e {
			return true
		}
	}
	return false
}
