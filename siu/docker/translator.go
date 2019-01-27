package docker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	errorss "github.com/pkg/errors"
)

var (
	// ErrEmptyCommand is returned when a caller didn't specify the command to translate.
	ErrEmptyCommand = errors.New("docker: provided command to translate is empty")
	// ErrEmptyArgs is returned when a caller didn't specify arguments for the command (container names/dynamic params).
	ErrEmptyArgs = errors.New("docker: no containers names/command arguments provided")
	// ErrWrongUserInput is returned when a user provided a response to the interactive prompt which is not valid.
	ErrWrongUserInput = errors.New("wrong input, non existing option was chosen")
)

// Lister is an interface describing an entity that can list docker images and possibly filter them
// basing on provided container names.
type Lister interface {
	List(containerNames []string) (outputLines []string, err error)
}

// Dictionary describes an entity able to translate an arbitrary command to a docker command.
type Dictionary interface {
	Get(word string) (translation Command, err error)
}

// Translator is a struct keeping raw arguments for processing during translation process.
// When it comes to how it interprets command the rule of thumb is that if command can run
// on multiple containers at once, it won't allow passing any dynamic paramaters for the command.
// If it's allowed to run on a single instance only, all additional arguments will be passed as
// dynamic parameters to the command.
type Translator struct {
	Dictionary
	Lister
	args              []string
	inputSource       *os.File
	outputDestination io.Writer
}

// Command represents a user defined translation of an arbitrary command to docker command.
// It also describes if provided translation can be executed on multiple containers at once.
type Command struct {
	Name                    string `yaml:"name"`
	Translation             string `yaml:"translation"`
	AllowMultipleContainers bool   `yaml:"multipleContainers"`
}

// New builds a Translator which based on passed command and arguments can invoke
// the command on several docker instance or build a compound command for a single instance.
func New(args []string) (*Translator, error) {
	// TODO: move args to Translate() as well
	if len(args) == 0 {
		return nil, ErrEmptyArgs
	}

	// TODO: pass config dir from cli.
	dict, err := NewYAMLDictionary(fmt.Sprintf("%s/.siu/config.yaml", os.Getenv("HOME")))
	if err != nil {
		return nil, err
	}

	return &Translator{
		Dictionary:        dict,
		Lister:            grepLister{},
		args:              args,
		inputSource:       os.Stdin,
		outputDestination: os.Stdout,
	}, nil
}

// Translate lists docker containers and greps the output according to passed params and based on that
// builds a final docker command.
func (t *Translator) Translate(command string) (string, error) {
	if command == "" {
		return "", ErrEmptyCommand
	}

	translation, err := t.Get(command)
	if err != nil {
		return "", err
	}

	containers, commandArgs := t.splitArguments(translation)
	ids, err := t.findContainersIDs(containers)
	if err != nil {
		return "", err
	}

	return t.buildFinalCommand(translation, commandArgs, ids), nil
}

func (t *Translator) buildFinalCommand(c Command, commandArgs, containersIDs []string) string {
	var ids string
	if c.AllowMultipleContainers {
		ids = strings.Join(containersIDs, " ")
	} else {
		ids = containersIDs[0]
	}

	translation := strings.Replace(c.Translation, "%s", "%[1]s", 5) // handle chained commands
	translationSlice := append([]string{translation}, commandArgs...)
	translationWithArgs := strings.Join(translationSlice, " ")
	return fmt.Sprintf(translationWithArgs, ids)
}

func (t *Translator) splitArguments(c Command) (containerNames []string, commandArgs []string) {
	if c.AllowMultipleContainers {
		containerNames = t.args
		return
	}

	containerNames = t.args[0:1]
	commandArgs = t.args[1:]
	return
}

func (t *Translator) findContainersIDs(names []string) ([]string, error) {
	lines, err := t.List(names)
	if err != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("there are no running containers with name matching any of %v", names)
	}

	return t.extractContinersIDsFromLines(lines)
}

func (t *Translator) extractContinersIDsFromLines(rawOutputLines []string) ([]string, error) {
	if len(rawOutputLines) > 1 {
		var chosenIDs []string

		lineIdxs, err := t.chooseContainerIdx(rawOutputLines)
		if err != nil {
			return nil, err
		}

		err = validateUserInput(lineIdxs, len(rawOutputLines))
		if err != nil {
			return nil, err
		}

		for _, idx := range lineIdxs {
			currentLine := rawOutputLines[idx]
			chosenIDs = append(chosenIDs, strings.Split(currentLine, " ")[0])
		}

		return chosenIDs, nil
	}

	return []string{strings.Split(rawOutputLines[0], " ")[0]}, nil
}

func validateUserInput(chosenIdxs []int, dockerOutputLength int) error {
	if len(chosenIdxs) == 1 && chosenIdxs[0] >= dockerOutputLength {
		return ErrWrongUserInput
	}

	if len(chosenIdxs) > dockerOutputLength {
		return ErrWrongUserInput
	}
	return nil
}

func (t *Translator) chooseContainerIdx(options []string) ([]int, error) {
	input, err := t.promptForInput(options)
	if err != nil {
		return nil, err
	}

	// normalize input -> get rid of commas etc
	cleanInput := strings.Replace(input, "\n", "", -1)
	cleanInput = strings.Replace(cleanInput, ", ", " ", -1)
	cleanInput = strings.Replace(cleanInput, ",", " ", -1)

	stringIDs := strings.Split(cleanInput, " ")
	ids := make([]int, len(stringIDs))
	for i, token := range stringIDs {
		val, err := strconv.Atoi(token)
		if err != nil {
			return nil, errorss.Wrap(err, "docker: couldn't convert string to integer")
		}

		ids[i] = val
	}

	return ids, nil
}

func (t *Translator) promptForInput(options []string) (string, error) {
	reader := bufio.NewReader(t.inputSource)

	prompt := fmt.Sprint("There are several containers matching provided name, which one do you want to use?\n")
	for i, option := range options {
		prompt += fmt.Sprintf("%d) %q\n", i, option)
	}
	prompt += ">"

	_, err := fmt.Fprint(t.outputDestination, prompt)
	if err != nil {
		return "", errorss.Wrap(err, "docker: couldn't print prompt to provided io.Writer")
	}
	return reader.ReadString('\n')
}
