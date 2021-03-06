package docker

import (
	"os/exec"
	"strings"

	errorss "github.com/pkg/errors"
)

type grepLister struct{}

func (g grepLister) List(containersNames []string) ([]string, error) {
	dockerOut, err := getGreppedDockerOut(containersNames)
	if err != nil {
		return nil, err
	}

	return splitOutput(dockerOut), nil
}

func getGreppedDockerOut(containersNames []string) (string, error) {
	ps := exec.Command("docker", "ps", "--format", `{{.ID}} {{.Names}}`)
	grep := exec.Command("grep", strings.Join(containersNames, "\\|"))

	pipe, err := ps.StdoutPipe()
	if err != nil {
		return "", errorss.Wrap(err, "docker: couldn't create pipe for getting container names")
	}
	grep.Stdin = pipe

	err = ps.Start()
	if err != nil {
		return "", errorss.Wrap(err, "docker: couldn't start pipe for getting container names")
	}

	out, err := grep.Output()
	if err != nil && !strings.Contains("exit status 1", err.Error()) {
		return "", errorss.Wrap(err, "docker: couldn't filter container names for standard Lister")
	}

	return string(out), nil
}

func splitOutput(out string) []string {
	lines := strings.Split(out, "\n")
	return lines[:len(lines)-1] // get rid of empty line
}
