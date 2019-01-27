package docker

import "fmt"

var (
	translations = map[string]Command{
		"lg":      Command{Translation: "docker logs -f %s", AllowMultipleContainers: false},
		"restart": Command{Translation: "docker restart %s", AllowMultipleContainers: true},
		"stop":    Command{Translation: "docker stop %s", AllowMultipleContainers: true},
		"exec":    Command{Translation: "docker exec -ti %s", AllowMultipleContainers: false},
		"sh":      Command{Translation: "docker exec -ti %s /bin/sh", AllowMultipleContainers: false},
		"bash":    Command{Translation: "docker exec -ti %s /bin/bash", AllowMultipleContainers: false},
		"rspec":   Command{Translation: "docker exec -ti %s rspec", AllowMultipleContainers: false},
	}
)

type InMemDictionary struct {
	translations map[string]Command
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

	i.translations[name] = Command{Translation: dockerCmd, AllowMultipleContainers: forMultipleContainers}
	return nil
}

// Get returns a translation for a provided name.
func (i *InMemDictionary) Get(name string) (Command, error) {
	cmd, ok := i.translations[name]
	if !ok {
		return cmd, fmt.Errorf("command %q does not exist", name)
	}

	return cmd, nil
}
