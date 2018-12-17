package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlackParse(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		path := "./fixtures/config.yml"
		conf := New(SlackParser)

		err := conf.ReadConfig(path)
		assert.Nil(t, err)

		srs := conf.Services
		assert.Len(t, srs, 1)
		assert.Len(t, srs[0].Sources(), 2)
		assert.Equal(t, srs[0].Sources()[0].Name(), "random")
		// assert.Equal(t, srs[0].Sources()[0].Interval(), 5*time.Second)
		assert.Equal(t, srs[0].Sources()[1].Name(), "general")
		// assert.Equal(t, srs[0].Sources()[1].Interval(), 10*time.Second)
	})
}
