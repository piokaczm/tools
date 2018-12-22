package config

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	ErrMalformedInterval = errors.New("config: interval is malformed, use format <number><unit> (eg. 10s)")

	timeDict = map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
	}
)

type Slack struct {
	Config SlackConfig `yaml:"slack"`
}

type SlackConfig struct {
	ApiToken string          `yaml:"api_token"`
	Channels []*SlackChannel `yaml:"channels"`
}

type SlackChannel struct {
	Name           string `yaml:"name"`
	IntervalString string `yaml:"interval"`
	IntervalTime   time.Duration
}

func slackParser(data []byte) (*Slack, error) {
	s := &Slack{}
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	for _, ch := range s.Config.Channels {
		parseInterval(ch)
	}

	return s, nil
}

func parseInterval(s *SlackChannel) error {
	for sign, unit := range timeDict {
		if strings.HasSuffix(s.IntervalString, sign) {
			stringWithoutUnit := s.IntervalString[0 : len(s.IntervalString)-1]
			val, err := strconv.Atoi(stringWithoutUnit)
			if err != nil {
				return err // TODO: wrap it
			}

			s.IntervalTime = time.Duration(val) * unit
			return nil
		}
	}

	return ErrMalformedInterval
}
