package docker

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYAMLGet(t *testing.T) {
	testConfigPath := "testdata/config.yaml"

	testCases := []struct {
		testName    string
		commandName string
		translation Command
		expectedErr string
	}{
		{
			testName:    "with stop command",
			commandName: "stop",
			translation: Command{Name: "stop", Translation: "docker stop %s", AllowMultipleContainers: true},
		},
		{
			testName:    "with exec command",
			commandName: "exec",
			translation: Command{Name: "exec", Translation: "docker exec -ti %s", AllowMultipleContainers: false},
		},
		{
			testName:    "with not existing command",
			commandName: "unavailable",
			expectedErr: "command \"unavailable\" does not exist",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			dict, err := NewYAMLDictionary(testConfigPath)
			if err != nil {
				t.Error("couldnt create test yaml dictionary")
			}

			res, err := dict.Get(tc.commandName)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tc.translation, res)
		})
	}
}

func TestYAMLAdd(t *testing.T) {
	testConfigPath := "testdata/config.yaml"

	testCases := []struct {
		testName    string
		commandName string
		translation Command
		expectedErr string
	}{
		{
			testName:    "adding existing command",
			commandName: "stop",
			translation: Command{Name: "stop", Translation: "docker stop %s", AllowMultipleContainers: true},
		},
		{
			testName:    "adding new command",
			commandName: "test",
			translation: Command{Name: "test", Translation: "docker test %s", AllowMultipleContainers: false},
		},
	}

	orgFile := cacheFile(testConfigPath)
	defer restoreFile(testConfigPath, orgFile)

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			dict, err := NewYAMLDictionary(testConfigPath)
			if err != nil {
				t.Error("couldnt create test yaml dictionary")
			}

			err = dict.Add(tc.commandName, tc.translation)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			assert.Nil(t, err)
			assertTranslationAdded(t, testConfigPath, tc.translation)
		})
	}
}

func assertTranslationAdded(t *testing.T, path string, cmd Command) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	var name, translation bool

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if !name {
			name = strings.Contains(line, cmd.Name)
			continue
		}

		if name && translation && strings.Contains(line, strconv.FormatBool(cmd.AllowMultipleContainers)) {
			count++
			continue
		}

		if name && strings.Contains(line, cmd.Translation) {
			translation = true
			continue
		}
	}

	if !(name && translation) {
		t.Errorf("command %v not found in the config file", cmd)
	}

	if count != 1 {
		t.Errorf("there should be one instance of command %v, found %d times", cmd, count)
	}
}

func containsTranslation(line string, cmd Command) bool {
	return strings.Contains(line, cmd.Name) && strings.Contains(line, cmd.Translation) &&
		strings.Contains(line, strconv.FormatBool(cmd.AllowMultipleContainers))
}

func cacheFile(path string) []byte {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	return d
}

func restoreFile(path string, data []byte) {
	err := os.Remove(path)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}

	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
