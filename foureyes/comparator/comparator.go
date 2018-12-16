package comparator

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNoMatches = errors.New("comparator: no matches") // TODO: rethink no match handling, err is weird, empty string is weird...

type Matcher struct {
	topics []string
}

type Source interface {
	Name() string
	Topics() ([][]string, error)
}

func New(topics []string) *Matcher {
	return &Matcher{topics}
}

func (m *Matcher) Match(source Source) (string, error) {
	sourceTopics, err := source.Topics()
	if err != nil {
		return "", err
	}

	intersection, exists := m.findMatches(sourceTopics)
	if !exists {
		return "", ErrNoMatches
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
	appended := make(map[string]struct{})

	for _, list := range s {
		for _, t := range list {
			if !m.include(t) {
				continue
			}

			if _, ok := appended[t]; ok {
				continue
			}

			appended[t] = struct{}{}
			matches = append(matches, t)
		}
	}

	return matches, len(matches) > 0
}

func (m *Matcher) include(e string) bool {
	for _, t := range m.topics {
		if strings.ToLower(t) == strings.ToLower(e) {
			return true
		}
	}
	return false
}
