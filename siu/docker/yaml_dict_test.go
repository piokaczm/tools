package docker

import (
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
