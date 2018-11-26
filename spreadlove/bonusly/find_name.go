package bonusly

import (
	"strings"
)

type nameFilter func(string) bool

func (c *Client) FindName(target string) ([]string, error) {
	me, err := c.me()
	if err != nil {
		return nil, err
	}

	rawUsers, err := c.users()
	if err != nil {
		return nil, err
	}

	return parseUsers(rawUsers, me.Result.Username, containsFilter(target))
}

func containsFilter(target string) nameFilter {
	return func(name string) bool {
		return !strings.Contains(name, target)
	}
}
