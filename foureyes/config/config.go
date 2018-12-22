package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type Config struct {
	Topics        []string `yaml:"topics"`
	SlackUsername string   `yaml:"slack_username"`
	SlackConfig   SlackConfig
}

func New() *Config {
	return &Config{
		Topics: make([]string, 0),
	}
}

func (c *Config) ReadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	slack, err := slackParser(data)
	if err != nil {
		return err
	}
	c.SlackConfig = slack.Config

	return yaml.Unmarshal(data, c)
}
