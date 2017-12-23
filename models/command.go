package models

import (
	"github.com/qb0C80aE/clay/extensions"
	"os/exec"
)

// Command is the model class what represents command to execute something
type Command struct {
	ID               int       `json:"id" form:"id"`
	Description      string    `json:"description" form:"description"`
	WorkingDirectory string    `json:"working_directory" form:"working_directory"`
	CommandLine      string    `json:"command_line" form:"command_line"`
	Status           string    `json:"status" form:"status"`
	StdOut           string    `json:"std_out" form:"std_out"`
	StdErr           string    `json:"std_err" form:"std_err"`
	ExitCode         int       `json:"exit_code" form:"exit_code"`
	Killed           bool      `json:"killed" form:"killed"`
	Cmd              *exec.Cmd `json:"cmd,omitempty"`
}

// NewCommandModel creates a Command model instance
func NewCommandModel() *Command {
	return &Command{}
}

var sharedCommandModel = NewCommandModel()

// SharedCommandModel returns the command model instance used as a model prototype and type analysis
func SharedCommandModel() *Command {
	return sharedCommandModel
}

func init() {
	extensions.RegisterModel(sharedCommandModel, false)
}
