package config

import (
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
)

type Slack struct {
	s struct {
		apiToken string         `yaml:"api_token"`
		channels []slackChannel `yaml:"channels"`
	} `yaml:"slack"`
}

type slackChannel struct {
	name     string `yaml:"name"`
	interval string `yaml:"interval"`
}

func SlackParser(data []byte) ([]Service, error) {
	var channels []Service

	var s interface{}
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}

	spew.Dump(s)
	return channels, nil
}
