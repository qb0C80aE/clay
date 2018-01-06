package logics

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
	"net/url"
	"os"
	"sort"
	"strconv"
	"sync"
)

var commandMap map[int]*models.Command = map[int]*models.Command{}
var commandMapMutex = new(sync.Mutex)

const (
	commandStatusCreated  = "created"
	commandStatusRunning  = "running"
	commandStatusFinished = "finished"
)

type commandLogic struct {
	*BaseLogic
}

func newCommandLogic() *commandLogic {
	logic := &commandLogic{
		BaseLogic: NewBaseLogic(
			models.SharedCommandModel(),
		),
	}
	return logic
}

func (logic *commandLogic) GetSingle(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {

	id, _ := strconv.Atoi(parameters.ByName("id"))

	commandMapMutex.Lock()
	defer commandMapMutex.Unlock()

	command, exists := commandMap[id]

	if !exists {
		return nil, errors.New("record not found")
	}

	return command, nil

}

func (logic *commandLogic) GetMulti(db *gorm.DB, _ gin.Params, _ url.Values, queryFields string) (interface{}, error) {

	commandMapMutex.Lock()
	defer commandMapMutex.Unlock()

	keys := make([]int, 0, len(commandMap))
	for key := range commandMap {
		keys = append(keys, key)
	}

	sort.Ints(keys)

	commands := make([]*models.Command, len(commandMap))

	for i, key := range keys {
		commands[i] = commandMap[key]
	}

	return commands, nil

}

func (logic *commandLogic) Create(db *gorm.DB, _ gin.Params, _ url.Values, data interface{}) (interface{}, error) {

	command := data.(*models.Command)

	if len(command.WorkingDirectory) == 0 {
		var err error
		if command.WorkingDirectory, err = os.Getwd(); err != nil {
			return nil, err
		}
	}

	command.Status = commandStatusCreated

	commandMapMutex.Lock()
	defer commandMapMutex.Unlock()

	id := len(commandMap) + 1
	command.ID = id
	commandMap[id] = command

	return command, nil

}

func (logic *commandLogic) Delete(db *gorm.DB, parameters gin.Params, _ url.Values) error {

	id, _ := strconv.Atoi(parameters.ByName("id"))

	commandMapMutex.Lock()
	defer commandMapMutex.Unlock()

	command, exists := commandMap[id]

	if !exists {
		commandMapMutex.Unlock()
		return errors.New("record not found")
	}

	if command.Status == commandStatusRunning {
		return errors.New("command is running")
	}

	delete(commandMap, id)

	return nil
}

func (logic *commandLogic) Total(_ *gorm.DB, _ interface{}) (int, error) {
	return len(commandMap), nil
}

var uniqueCommandLogic = newCommandLogic()

// UniqueCommandLogic returns the unique command logic instance
func UniqueCommandLogic() extensions.Logic {
	return uniqueCommandLogic
}

func init() {
	extensions.RegisterLogic(models.SharedCommandModel(), UniqueCommandLogic())
}
