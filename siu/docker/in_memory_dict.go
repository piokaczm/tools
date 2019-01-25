package docker

import "fmt"

var (
	translations = map[string]dockerCommand{
		"lg":      dockerCommand{name: "docker logs -f %s", allowMultipleContainers: false},
		"restart": dockerCommand{name: "docker restart %s", allowMultipleContainers: true},
		"stop":    dockerCommand{name: "docker stop %s", allowMultipleContainers: true},
		"exec":    dockerCommand{name: "docker exec -ti %s", allowMultipleContainers: false},
		"sh":      dockerCommand{name: "docker exec -ti %s /bin/sh", allowMultipleContainers: false},
		"bash":    dockerCommand{name: "docker exec -ti %s /bin/bash", allowMultipleContainers: false},
		"rspec":   dockerCommand{name: "docker exec -ti %s rspec", allowMultipleContainers: false},
	}
)

type InMemDictionary struct {
	translations map[string]dockerCommand
}

func NewInMemDictionary() *InMemDictionary {
	return &InMemDictionary{
		translations: translations,
	}
}

// Add adds a docker command equal to `dockerCmd` under `name` key.
// This command is not used yet.
func (i *InMemDictionary) Add(name, dockerCmd string, forMultipleContainers bool) error {
	if _, ok := i.translations[name]; ok {
		return fmt.Errorf("command %q already exists", name)
	}

	i.translations[name] = dockerCommand{name: dockerCmd, allowMultipleContainers: forMultipleContainers}
	return nil
}

// Get returns a translation for a provided name.
func (i *InMemDictionary) Get(name string) (dockerCommand, error) {
	cmd, ok := i.translations[name]
	if !ok {
		return cmd, fmt.Errorf("command %q does not exist", name)
	}

	return cmd, nil
}
