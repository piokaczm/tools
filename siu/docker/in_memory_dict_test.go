package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryGet(t *testing.T) {
	testCases := []struct {
		testName    string
		commandName string
		translation Command
		expectedErr string
	}{
		{
			testName:    "with lg command",
			commandName: "lg",
			translation: Command{Translation: "docker logs -f %s", AllowMultipleContainers: false},
		},
		{
			testName:    "with stop command",
			commandName: "stop",
			translation: Command{Translation: "docker stop %s", AllowMultipleContainers: true},
		},
		{
			testName:    "with exec command",
			commandName: "exec",
			translation: Command{Translation: "docker exec -ti %s", AllowMultipleContainers: false},
		},
		{
			testName:    "with restart command",
			commandName: "restart",
			translation: Command{Translation: "docker restart %s", AllowMultipleContainers: true},
		},
		{
			testName:    "with sh command",
			commandName: "sh",
			translation: Command{Translation: "docker exec -ti %s /bin/sh", AllowMultipleContainers: false},
		},
		{
			testName:    "with bash command",
			commandName: "bash",
			translation: Command{Translation: "docker exec -ti %s /bin/bash", AllowMultipleContainers: false},
		},
		{
			testName:    "with rspec command",
			commandName: "rspec",
			translation: Command{Translation: "docker exec -ti %s rspec", AllowMultipleContainers: false},
		},
		{
			testName:    "with not existing command",
			commandName: "unavailable",
			expectedErr: "command \"unavailable\" does not exist",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			dict := NewInMemDictionary()
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
