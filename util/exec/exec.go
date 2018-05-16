package exec

import (
	"os/exec"
)

var utility = &Utility{}

// Utility handles command execution operation
type Utility struct {
}

// GetUtility returns the instance of utility
func GetUtility() *Utility {
	return utility
}

// Command creates exec.Cmd instance
func (receiver *Utility) Command(commandLine []string) *exec.Cmd {
	return exec.Command(commandLine[0], commandLine[1:]...)
}
