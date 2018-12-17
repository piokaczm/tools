package config

import (
	"io/ioutil"
	"time"
)

type Config struct {
	Topics   []string
	Services []Service
	Parsers  []Parser
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

func New(parsers ...Parser) *Config {
	return &Config{
		make([]string, 0),
		make([]Service, 0),
		parsers,
	}
}

func (c *Config) ReadConfig(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	for _, parse := range c.Parsers {
		services, err := parse(data)
		if err != nil {
			return err
		}

		c.Services = append(c.Services, services...)
	}

	return nil
}
