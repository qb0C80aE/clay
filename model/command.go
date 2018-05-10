package model

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/util/mapstruct"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"sync"
)

var commandNameCommandMap = map[string]*Command{}
var commandNameCommandMapMutex = new(sync.Mutex)

const (
	commandStatusCreated  = "created"
	commandStatusRunning  = "running"
	commandStatusFinished = "finished"
)

// Command is the model class what represents command to execute something
type Command struct {
	Base
	Name             string    `json:"name" yaml:"name" form:"name" clay:"key_parameter" validate:"required"`
	Description      string    `json:"description" yaml:"description" form:"description"`
	WorkingDirectory string    `json:"working_directory" yaml:"working_directory" form:"working_directory"`
	CommandLine      string    `json:"command_line" yaml:"command_line" form:"command_line" validate:"required"`
	Status           string    `json:"status" yaml:"status" form:"status"`
	StdOut           string    `json:"std_out" yaml:"std_out" form:"std_out"`
	StdErr           string    `json:"std_err" yaml:"std_err" form:"std_err"`
	ExitCode         int       `json:"exit_code" yaml:"exit_code" form:"exit_code"`
	Killed           bool      `json:"killed" yaml:"killed" form:"killed"`
	StartedAt        string    `json:"started_at" yaml:"started_at" form:"started_at"`
	FinishedAt       string    `json:"finished_at" yaml:"finished_at" form:"finished_at"`
	Cmd              *exec.Cmd `json:"cmd,omitempty" yaml:"cmd,omitempty"`
}

// NewCommand creates a command model instance
func NewCommand() *Command {
	return &Command{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *Command) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *Command) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	result, exists := commandNameCommandMap[name]

	if !exists {
		return nil, errors.New("record not found")
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *Command) GetMulti(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	commandNameCommandMapMutex.Lock()
	defer commandNameCommandMapMutex.Unlock()

	keyList := make([]string, 0, len(commandNameCommandMap))
	for key := range commandNameCommandMap {
		keyList = append(keyList, key)
	}

	sort.Strings(keyList)

	commandList := make([]*Command, len(commandNameCommandMap))

	for i, key := range keyList {
		commandList[i] = commandNameCommandMap[key]
	}

	return commandList, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *Command) Create(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	command := NewCommand()
	if err := mapstruct.RemapToStruct(inputContainer, command); err != nil {
		return nil, err
	}

	if len(command.WorkingDirectory) == 0 {
		var err error
		if command.WorkingDirectory, err = os.Getwd(); err != nil {
			return nil, err
		}
	}

	command.Status = commandStatusCreated

	commandNameCommandMapMutex.Lock()
	defer commandNameCommandMapMutex.Unlock()

	if _, exists := commandNameCommandMap[command.Name]; exists {
		logging.Logger().Debugf("command %s already exists", command.Name)
		return nil, fmt.Errorf("command %s already exists", command.Name)
	}

	commandNameCommandMap[command.Name] = command

	return command, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *Command) Delete(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values) error {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	commandNameCommandMapMutex.Lock()
	defer commandNameCommandMapMutex.Unlock()

	result, exists := commandNameCommandMap[name]

	if !exists {
		commandNameCommandMapMutex.Unlock()
		return errors.New("record not found")
	}

	if result.Status == commandStatusRunning {
		return errors.New("command is running")
	}

	delete(commandNameCommandMap, name)

	return nil
}

// GetCount returns the record count under current db conditions
func (receiver *Command) GetCount(_ extension.Model, _ *gorm.DB) (int, error) {
	return len(commandNameCommandMap), nil
}

func init() {
	extension.RegisterModel(NewCommand())
}
