package docker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	errorss "github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ErrNoPath is returned when path for dictionary config is not defined.
var ErrNoPath = errors.New("no path for dictionary config provided")

// YAMLDictionary is an impelemntation of Dictionary based on YAML caonfig file.
type YAMLDictionary struct {
	path         string
	translations map[string]Command
}

type commands struct {
	Translations []Command `yaml:"translations"`
}

// NewYAMLDictionary returns a new YAMLDictionary with transtaltions defined in a config located in `path`.
func NewYAMLDictionary(path string) (*YAMLDictionary, error) {
	if path == "" {
		return nil, ErrNoPath
	}

	y := &YAMLDictionary{
		translations: make(map[string]Command),
		path:         path,
	}
	err := y.readConfig()
	return y, err
}

// Add creates a new command in the config file.
func (y *YAMLDictionary) Add(name string, cmd Command) error {
	// get current state of the config
	newConfig := make(map[string]Command)
	for k, v := range y.translations {
		newConfig[k] = v
	}

	newConfig[name] = cmd
	return y.overwriteDict(newConfig)
}

// Get returns a translation for a provided name.
func (y *YAMLDictionary) Get(name string) (Command, error) {
	cmd, ok := y.translations[name]
	if !ok {
		return cmd, fmt.Errorf("command %q does not exist", name)
	}

	return cmd, nil
}

func (y *YAMLDictionary) overwriteDict(dict map[string]Command) error {
	f, err := os.Create(y.path)
	if err != nil {
		return errorss.Wrapf(err, "couldn't open dict file %q", y.path)
	}
	defer f.Close()

	err = yaml.NewEncoder(f).Encode(dict)
	if err != nil {
		return errorss.Wrap(err, "couldn't write new dict file")
	}
	return nil
}

func (y *YAMLDictionary) readConfig() error {
	rawConf, err := ioutil.ReadFile(y.path)
	if err != nil {
		return errorss.Wrap(err, "couldn't read config file")
	}

	cmds := commands{}
	err = yaml.Unmarshal(rawConf, &cmds)
	if err != nil {
		return errorss.Wrap(err, "couldn't unmarshall config file")
	}

	for _, c := range cmds.Translations {
		y.translations[c.Name] = c
	}
	return nil
}
