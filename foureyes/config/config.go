package config

import (
	"io/ioutil"
	"time"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Topics        []string
	SlackChannels []SlackChannel
}

type Service interface {
	Sources() []ServiceSource
	Auth() string
}

type ServiceSource interface {
	Name() string
	Interval() time.Duration
}

type Parser func([]byte) ([]Service, error)

func New() *Config {
	return &Config{
		make([]string, 0),
		make([]SlackChannel, 0),
	}
}

func (c *Config) ReadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	slackChannels, err := slackParser(data)
	if err != nil {
		return err
	}
	c.SlackChannels = slackChannels

	topics, err := topicsParser(data)
	if err != nil {
		return err
	}

	c.Topics = topics

	return nil
}

func topicsParser(data []byte) ([]string, error) {
	ts := struct {
		Topics []string `yaml:"topics"`
	}{}

	err := yaml.Unmarshal(data, &ts)
	if err != nil {
		return nil, err
	}

	return ts.Topics, nil
}
