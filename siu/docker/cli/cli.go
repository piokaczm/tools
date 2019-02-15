package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/piokaczm/tools/siu/docker"
	errorss "github.com/pkg/errors"
)

const (
	configDir      = "%s/.siu" // TODO: allow passing different config using flags
	configFileName = "config.yaml"
	configTemplate = `translations:
# translations syntax is like that (remember to use spaces, not tabs!):
# - name: "exec"
#   translation: "docker exec -ti %s"
#   multipleContainers: false`
	helpText = `use one of the commands:
	- init : creates a config file in ~/.siu/config.yaml
	- edit : open config file using nano
	- help : show description of available commands
	- <cmd_name> <container_name> <args> : tries to get the command translation from the config file and build docker command`
)

var (
	// ErrNotEnoughArgs is returned when user provided not enough arguments to execute any operation.
	ErrNotEnoughArgs = errors.New("not enough arguments provided")
)

// Invoke is an entry point for docker translations CLI. It handles basic checks and delegates
// commands to proper execution paths.
func Invoke(args []string) error {
	if len(args) == 0 {
		return ErrNotEnoughArgs
	}

	switch args[0] { // base cmd
	case "init":
		if err := initConfig(); err != nil {
			return err
		}
	case "help", "--help", "-h": // TODO: properly handle flags
		fmt.Println(helpText)
	default:
		err := translate(args)
		if _, ok := err.(docker.CommandError); ok {
			err = fmt.Errorf("%s, see 'siu -help' for help", err)
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func initConfig() error {
	dir := fmt.Sprintf(configDir, os.Getenv("HOME"))
	err := os.MkdirAll(dir, os.ModePerm) // TODO: support non osx systems
	if err != nil {
		return errorss.Wrap(err, "couldn't create config dir")
	}

	f, err := os.Create(dir + "/" + configFileName)
	if err != nil {
		return errorss.Wrap(err, "couldn't create config file")
	}
	defer f.Close()

	_, err = f.Write([]byte(configTemplate))
	return errorss.Wrap(err, "couldn't write to config file")
}

func translate(args []string) error {
	if len(args) < 2 {
		return ErrNotEnoughArgs // min 2 args for translation (cmd, container)
	}

	t, err := docker.New(args[1:])
	if err != nil {
		return err
	}

	v, err := t.Translate(args[0])
	if err != nil {
		return err
	}

	fmt.Println(v)
	return nil
}
