package docker

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	ErrEmptyCommand       = errors.New("docker: provided command to translate is empty")
	ErrEmptyContainerName = errors.New("docker: provided container name is empty")

	translations = map[string]dockerCommand{
		"sh":      dockerCommand{name: "docker exec -ti %s /bin/sh", allowMultipleContainers: false},
		"bash":    dockerCommand{name: "docker exec -ti %s /bin/bash", allowMultipleContainers: false},
		"lg":      dockerCommand{name: "docker logs -f %s", allowMultipleContainers: false},
		"restart": dockerCommand{name: "docker restart ", allowMultipleContainers: true},
		"stop":    dockerCommand{name: "docker stop ", allowMultipleContainers: true},
	}
)

type dockerCommand struct {
	name                    string
	allowMultipleContainers bool
}

type Translator struct {
	command   string
	container string
}

func New(command, container string) (*Translator, error) {
	if command == "" {
		return nil, ErrEmptyCommand
	}

	if container == "" {
		return nil, ErrEmptyContainerName
	}

	return &Translator{
		command,
		container,
	}, nil
}

func (t *Translator) Translate() (string, error) {
	ids, err := findContainerID(t.container)
	if err != nil {
		return "", err
	}

	translation, ok := translations[t.command]
	if !ok {
		return "", fmt.Errorf("translation for command %q not found", t.command)
	}

	if !translation.allowMultipleContainers && len(ids) > 1 {
		return "", fmt.Errorf("command %q is not supporting running on multiple containers", t.command)
	}

	var command string
	if translation.allowMultipleContainers {
		ids := strings.Join(ids, " ")
		command = translation.name + ids
	} else {
		command = fmt.Sprintf(translation.name, ids[0])
	}

	return command, nil
}

func findContainerID(name string) ([]string, error) {
	out, err := getDockerOutput(name)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	lines = lines[:len(lines)-1] // get rid of empty line
	if len(lines) == 0 {
		return nil, fmt.Errorf("there are no running containers with name %q", name)
	}

	// TODO: split it for god's sake
	var chosenIDs []string
	if len(lines) > 1 {
		lineIdx, err := chooseContainerIdx(lines)
		if err != nil {
			return nil, err
		}

		if len(lineIdx) > len(lines) {
			return nil, errors.New("wrong input, non existing option was chosen")
		}

		for _, idx := range lineIdx {
			currentLine := lines[idx]
			chosenIDs = append(chosenIDs, strings.Split(currentLine, " ")[0])
		}
	} else {
		chosenIDs = []string{strings.Split(lines[0], " ")[0]}
	}

	return chosenIDs, nil
}

func chooseContainerIdx(options []string) ([]int, error) {
	reader := bufio.NewReader(os.Stdin)

	prompt := fmt.Sprint("There are several containers matching provided name, which one do you want to use?\n")
	for i, option := range options {
		prompt += fmt.Sprintf("%d) %q\n", i, option)
	}
	prompt += ">"

	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	cleanInput := strings.Replace(input, "\n", "", -1)
	cleanInput = strings.Replace(cleanInput, ", ", " ", -1)
	cleanInput = strings.Replace(cleanInput, ",", " ", -1)

	stringIDs := strings.Split(cleanInput, " ")
	ids := make([]int, len(stringIDs))
	for i, token := range stringIDs {
		val, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		}

		ids[i] = val
	}

	return ids, nil
}

func getDockerOutput(name string) (string, error) {
	ps := exec.Command("docker", "ps", "--format", `{{.ID}} {{.Names}}`)
	grep := exec.Command("grep", name)

	pipe, err := ps.StdoutPipe()
	if err != nil {
		return "", nil
	}
	grep.Stdin = pipe

	err = ps.Start()
	if err != nil {
		return "", nil
	}

	out, err := grep.Output()
	if err != nil {
		return "", nil
	}

	return string(out), nil
}
