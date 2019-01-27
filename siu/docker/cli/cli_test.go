package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvoke(t *testing.T) {
	testCases := []struct {
		name        string
		args        []string
		expectedErr error
	}{
		{
			name:        "with zero arguments",
			expectedErr: ErrNotEnoughArgs,
		},
		{
			name:        "with not enough arguments for running translation",
			args:        []string{"exec"},
			expectedErr: ErrNotEnoughArgs,
		},
	}

	// for now just test error handling
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Invoke(tc.args)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
