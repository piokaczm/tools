package docker

import (
	"errors"
	"fmt"
	"io/ioutil"

	errorss "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ErrNoPath is returned when path for dictionary config is not defined.
var ErrNoPath = errors.New("no path for dictionary config provided")

// YAMLDictionary is an impelemntation of Dictionary based on YAML caonfig file.
type YAMLDictionary struct {
	translations map[string]Command
}

// NewYAMLDictionary returns a new YAMLDictionary with transtaltions defined in a config located in `path`.
func NewYAMLDictionary(path string) (*YAMLDictionary, error) {
	if path == "" {
		return nil, ErrNoPath
	}

	y := &YAMLDictionary{translations: make(map[string]Command)}
	err := y.readConfig(path)
	return y, err
}

func (y *YAMLDictionary) Add(name string, cmd Command) error {
	return nil
}

// Get returns a translation for a provided name.
func (y *YAMLDictionary) Get(name string) (Command, error) {
	cmd, ok := y.translations[name]
	if !ok {
		return cmd, fmt.Errorf("command %q does not exist", name)
	}

	return cmd, nil
}

func (y *YAMLDictionary) readConfig(path string) error {
	rawConf, err := ioutil.ReadFile(path)
	if err != nil {
		return errorss.Wrap(err, "couldn't read config file")
	}

	var commands struct {
		Translations []Command `yaml:"translations"`
	}

	err = yaml.Unmarshal(rawConf, &commands)
	if err != nil {
		return errorss.Wrap(err, "couldn't unmarshall config file")
	}

	for _, c := range commands.Translations {
		y.translations[c.Name] = c
	}
	return nil
}
