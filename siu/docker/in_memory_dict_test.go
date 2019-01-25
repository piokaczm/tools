package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		testName    string
		commandName string
		translation dockerCommand
		expectedErr string
	}{
		{
			testName:    "with lg command",
			commandName: "lg",
			translation: dockerCommand{name: "docker logs -f %s", allowMultipleContainers: false},
		},
		{
			testName:    "with stop command",
			commandName: "stop",
			translation: dockerCommand{name: "docker stop %s", allowMultipleContainers: true},
		},
		{
			testName:    "with exec command",
			commandName: "exec",
			translation: dockerCommand{name: "docker exec -ti %s", allowMultipleContainers: false},
		},
		{
			testName:    "with restart command",
			commandName: "restart",
			translation: dockerCommand{name: "docker restart %s", allowMultipleContainers: true},
		},
		{
			testName:    "with sh command",
			commandName: "sh",
			translation: dockerCommand{name: "docker exec -ti %s /bin/sh", allowMultipleContainers: false},
		},
		{
			testName:    "with bash command",
			commandName: "bash",
			translation: dockerCommand{name: "docker exec -ti %s /bin/bash", allowMultipleContainers: false},
		},
		{
			testName:    "with rspec command",
			commandName: "rspec",
			translation: dockerCommand{name: "docker exec -ti %s rspec", allowMultipleContainers: false},
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
