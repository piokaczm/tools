package docker

import "fmt"

const (
	emptyCommand     = "empty_command"
	translationError = "translation_error"
)

// CommandError represents an error regarding commands - it can be lack of translation for the command,
// an empty command or anything related.
type CommandError struct {
	err     error
	errType string
}

// NewCommandError returns a new instance od CommandError.
func NewCommandError(err error, errType string) CommandError {
	return CommandError{
		err:     err,
		errType: errType,
	}
}

// Error returns a formatted error message.
func (c CommandError) Error() string {
	return fmt.Sprintf("command error (%s): %s", c.errType, c.err)
}
