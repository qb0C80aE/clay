package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"net/url"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"sync"
)

var commandIDCommandMap = map[int]*Command{}
var commandIDCommandMapMutex = new(sync.Mutex)

const (
	commandStatusCreated  = "created"
	commandStatusRunning  = "running"
	commandStatusFinished = "finished"
)

// Command is the model class what represents command to execute something
type Command struct {
	*Base            `json:"base,omitempty"`
	ID               int       `json:"id" form:"id"`
	Description      string    `json:"description" form:"description"`
	WorkingDirectory string    `json:"working_directory" form:"working_directory"`
	CommandLine      string    `json:"command_line" form:"command_line"`
	Status           string    `json:"status" form:"status"`
	StdOut           string    `json:"std_out" form:"std_out"`
	StdErr           string    `json:"std_err" form:"std_err"`
	ExitCode         int       `json:"exit_code" form:"exit_code"`
	Killed           bool      `json:"killed" form:"killed"`
	StartedAt        string    `json:"started_at" form:"started_at"`
	FinishedAt       string    `json:"finished_at" form:"finished_at"`
	Cmd              *exec.Cmd `json:"cmd,omitempty"`
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *Command) GetSingle(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	id, _ := strconv.Atoi(parameters.ByName("id"))

	commandIDCommandMapMutex.Lock()
	defer commandIDCommandMapMutex.Unlock()

	result, exists := commandIDCommandMap[id]

	if !exists {
		return nil, errors.New("record not found")
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *Command) GetMulti(db *gorm.DB, _ gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	commandIDCommandMapMutex.Lock()
	defer commandIDCommandMapMutex.Unlock()

	keyList := make([]int, 0, len(commandIDCommandMap))
	for key := range commandIDCommandMap {
		keyList = append(keyList, key)
	}

	sort.Ints(keyList)

	commandList := make([]*Command, len(commandIDCommandMap))

	for i, key := range keyList {
		commandList[i] = commandIDCommandMap[key]
	}

	return commandList, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *Command) Create(db *gorm.DB, _ gin.Params, _ url.Values, input extension.Model) (interface{}, error) {
	command := input.(*Command)

	if len(command.WorkingDirectory) == 0 {
		var err error
		if command.WorkingDirectory, err = os.Getwd(); err != nil {
			return nil, err
		}
	}

	command.Status = commandStatusCreated

	commandIDCommandMapMutex.Lock()
	defer commandIDCommandMapMutex.Unlock()

	id := len(commandIDCommandMap) + 1
	command.ID = id
	commandIDCommandMap[id] = command

	return command, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *Command) Delete(db *gorm.DB, parameters gin.Params, _ url.Values) error {
	id, _ := strconv.Atoi(parameters.ByName("id"))

	commandIDCommandMapMutex.Lock()
	defer commandIDCommandMapMutex.Unlock()

	result, exists := commandIDCommandMap[id]

	if !exists {
		commandIDCommandMapMutex.Unlock()
		return errors.New("record not found")
	}

	if result.Status == commandStatusRunning {
		return errors.New("command is running")
	}

	delete(commandIDCommandMap, id)

	return nil
}

// GetTotal returns the count of for multi resource
func (receiver *Command) GetTotal(_ *gorm.DB) (int, error) {
	return len(commandIDCommandMap), nil
}

// NewCommand creates a command model instance
func NewCommand() *Command {
	return CreateModel(&Command{}).(*Command)
}

func init() {
	extension.RegisterModel(NewCommand(), false)
}
