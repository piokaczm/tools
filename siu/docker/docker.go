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

	translations = map[string]string{
		"sh":   "docker exec -ti %s /bin/sh",
		"bash": "docker exec -ti %s /bin/bash",
		"lg":   "docker logs -f %s",
	}
)

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
	id, err := findContainerID(t.container)
	if err != nil {
		return "", err
	}

	translation, ok := translations[t.command]
	if !ok {
		return "", fmt.Errorf("translation for command %q not found", t.command)
	}

	return fmt.Sprintf(translation, id), nil
}

func findContainerID(name string) (string, error) {
	out, err := getDockerOutput(name)
	if err != nil {
		return "", nil
	}

	lines := strings.Split(out, "\n")
	lines = lines[:len(lines)-1] // get rid of empty line
	if len(lines) == 0 {
		return "", fmt.Errorf("there are no running containers with name %q", name)
	}

	chosenLine := lines[0]
	if len(lines) > 1 {
		lineIdx, err := chooseContainerIdx(lines)
		if err != nil {
			return "", err
		}

		if lineIdx >= len(lines) {
			return "", errors.New("wrong input, non existing option was chosen")
		}
		chosenLine = lines[lineIdx]
	}

	return strings.Split(chosenLine, " ")[0], nil
}

func chooseContainerIdx(options []string) (int, error) {
	reader := bufio.NewReader(os.Stdin)

	prompt := fmt.Sprint("There are several containers matching provided name, which one do you want to use?\n")
	for i, option := range options {
		prompt += fmt.Sprintf("%d) %q\n", i, option)
	}
	prompt += ">"

	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, err
	}

	return strconv.Atoi(strings.Replace(input, "\n", "", -1))
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
