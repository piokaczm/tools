package docker

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testLister struct {
	output []string
	err    error
}

func (t testLister) List([]string) ([]string, error) {
	return t.output, t.err
}

func getTestLister(out []string, err error) testLister {
	return testLister{
		output: out,
		err:    err,
	}
}

func prepareUserInput(input string) (*os.File, error) {
	f, err := ioutil.TempFile("", "example_imput")
	if err != nil {
		return nil, err
	}

	_, err = f.Write([]byte(input))
	if err != nil {
		return nil, err
	}
	_, err = f.Seek(0, 0)
	return f, err
}

func TestTranslate(t *testing.T) {
	testCases := []struct {
		name           string
		dockerOutput   []string
		dockerErr      error
		command        string
		args           []string
		userInput      string
		expectedOutput string
		expectedErr    string
	}{
		// different input/output combinations
		{
			name:           "with single container command with no extra args and full container name",
			command:        "lg",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom"},
			expectedOutput: "docker logs -f f4321fb02881",
		},
		{
			name:           "with single container command with extra args and full container name",
			command:        "exec",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom", "ls", "-al"},
			expectedOutput: "docker exec -ti f4321fb02881 ls -al",
		},
		{
			name:           "with single container command with no extra args and partial container name",
			command:        "lg",
			dockerOutput:   []string{"f4321fb02881 my_mom", "f4321fb02882 your_mom"},
			args:           []string{"mom"},
			userInput:      "0\n",
			expectedOutput: "docker logs -f f4321fb02881",
		},
		{
			name:           "with single container command with extra args and partial container name",
			command:        "exec",
			dockerOutput:   []string{"f4321fb02881 my_mom", "f4321fb02882 your_mom"},
			args:           []string{"my_mom", "ls", "-al"},
			userInput:      "0\n",
			expectedOutput: "docker exec -ti f4321fb02881 ls -al",
		},
		{
			name:    "with multiple containers command and comma and white space separated input",
			command: "stop",
			dockerOutput: []string{
				"f4321fb02880 mommy",
				"f4321fb02881 my_mom",
				"f4321fb02882 your_mom",
				"f4321fb02883 your_dad",
			},
			args:           []string{"mom", "dad"},
			userInput:      "0, 2\n",
			expectedOutput: "docker stop f4321fb02880 f4321fb02882",
		},
		{
			name:    "with multiple containers command and comma separated input",
			command: "stop",
			dockerOutput: []string{
				"f4321fb02880 mommy",
				"f4321fb02881 my_mom",
				"f4321fb02882 your_mom",
				"f4321fb02883 your_dad",
			},
			args:           []string{"mom", "dad"},
			userInput:      "0,2\n",
			expectedOutput: "docker stop f4321fb02880 f4321fb02882",
		},
		{
			name:    "with multiple containers command and white space separated input",
			command: "stop",
			dockerOutput: []string{
				"f4321fb02880 mommy",
				"f4321fb02881 my_mom",
				"f4321fb02882 your_mom",
				"f4321fb02883 your_dad",
			},
			args:           []string{"mom", "dad"},
			userInput:      "0 2\n",
			expectedOutput: "docker stop f4321fb02880 f4321fb02882",
		},
		// specific commands translations
		{
			name:           "with lg command",
			command:        "lg",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom"},
			expectedOutput: "docker logs -f f4321fb02881",
		},
		{
			name:           "with stop command",
			command:        "stop",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom"},
			expectedOutput: "docker stop f4321fb02881",
		},
		{
			name:           "with exec command",
			command:        "exec",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom", "ls"},
			expectedOutput: "docker exec -ti f4321fb02881 ls",
		},
		{
			name:           "with restart command",
			command:        "restart",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom"},
			expectedOutput: "docker restart f4321fb02881",
		},
		{
			name:           "with sh command",
			command:        "sh",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom"},
			expectedOutput: "docker exec -ti f4321fb02881 /bin/sh",
		},
		{
			name:           "with bash command",
			command:        "bash",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom"},
			expectedOutput: "docker exec -ti f4321fb02881 /bin/bash",
		},
		{
			name:           "with rspec command",
			command:        "rspec",
			dockerOutput:   []string{"f4321fb02881 my_mom"},
			args:           []string{"my_mom"},
			expectedOutput: "docker exec -ti f4321fb02881 rspec",
		},
		// errors tests
		{
			name:        "with not existing command",
			command:     "dumb",
			args:        []string{"my_mom"},
			expectedErr: "translation for command \"dumb\" not found",
		},
		{
			name:         "with no matching containers",
			command:      "restart",
			args:         []string{"my_mom", "your_mum"},
			dockerOutput: []string{},
			expectedErr:  "there are no running containers with name matching any of [my_mom your_mum]",
		},
		{
			name:         "with wrong user input",
			command:      "exec",
			dockerOutput: []string{"f4321fb02881 my_mom", "f4321fb02882 your_mom"},
			args:         []string{"my_mom", "ls", "-al"},
			userInput:    "3\n",
			expectedErr:  "wrong input, non existing option was chosen",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			translator, err := New(testCase.command, testCase.args)
			if err != nil {
				t.Errorf("expected no errors, got %s", err)
			}
			translator.outputDestination = ioutil.Discard
			translator.Lister = getTestLister(testCase.dockerOutput, testCase.dockerErr)

			tempFile, err := prepareUserInput(testCase.userInput)
			if err != nil {
				t.Errorf("expected no errors, got %s", err)
			}
			defer os.Remove(tempFile.Name())
			translator.inputSource = tempFile

			output, err := translator.Translate()
			if testCase.expectedErr != "" {
				assert.EqualError(t, err, testCase.expectedErr)
				return
			}

			if err != nil {
				t.Errorf("expected no errors, got %s", err)
			}
			assert.Equal(t, testCase.expectedOutput, output)
		})
	}
}
