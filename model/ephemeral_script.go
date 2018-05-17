package model

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
	"github.com/robertkrimen/otto"
	"net/url"
	"sort"
	"sync"
)

var nameEphemeralScriptMap = map[string]*EphemeralScript{}
var nameEphemeralScriptMapMutex = new(sync.Mutex)

const (
	ephemeralScriptStatusCreated  = "created"
	ephemeralScriptStatusRunning  = "running"
	ephemeralScriptStatusFinished = "finished"
)

// EphemeralScript is the model class what represents ephemeralScript to execute something
type EphemeralScript struct {
	Base
	Name          string `json:"name" yaml:"name" form:"name" gorm:"primary_key" validate:"required"`
	Description   string `json:"description" yaml:"description" form:"description"`
	ScriptContent string `json:"script_content" yaml:"script_content" form:"script_content" sql:"type:text"`
	Status        string `json:"status" yaml:"status"`
	ReturnValue   string `json:"return_value" yaml:"return_value" sql:"type:text"`
	Error         string `json:"error" yaml:"error" sql:"type:text"`
	StartedAt     string `json:"started_at" yaml:"started_at"`
	FinishedAt    string `json:"finished_at" yaml:"finished_at"`
	vm            *otto.Otto
}

// NewEphemeralScript creates a ephemeralScript model instance
func NewEphemeralScript() *EphemeralScript {
	return &EphemeralScript{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralScript) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralScript) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	result, exists := nameEphemeralScriptMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *EphemeralScript) GetMulti(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	nameEphemeralScriptMapMutex.Lock()
	defer nameEphemeralScriptMapMutex.Unlock()

	keyList := make([]string, 0, len(nameEphemeralScriptMap))
	for key := range nameEphemeralScriptMap {
		keyList = append(keyList, key)
	}

	sort.Strings(keyList)

	ephemeralScriptList := make([]*EphemeralScript, len(nameEphemeralScriptMap))

	for i, key := range keyList {
		ephemeralScriptList[i] = nameEphemeralScriptMap[key]
	}

	return ephemeralScriptList, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *EphemeralScript) Create(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	ephemeralScript := NewEphemeralScript()
	if err := mapstructutilpkg.GetUtility().MapToStruct(inputContainer, ephemeralScript); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	nameEphemeralScriptMapMutex.Lock()
	defer nameEphemeralScriptMapMutex.Unlock()

	if _, exists := nameEphemeralScriptMap[ephemeralScript.Name]; exists {
		logging.Logger().Debugf("ephemeralScript %s already exists", ephemeralScript.Name)
		return nil, fmt.Errorf("ephemeralScript %s already exists", ephemeralScript.Name)
	}

	ephemeralScript.Status = ephemeralScriptStatusCreated
	ephemeralScript.ReturnValue = ""
	ephemeralScript.Error = ""
	ephemeralScript.StartedAt = ""
	ephemeralScript.FinishedAt = ""
	ephemeralScript.vm = otto.New()
	ephemeralScript.vm.Interrupt = make(chan func(), 1)

	nameEphemeralScriptMap[ephemeralScript.Name] = ephemeralScript

	return ephemeralScript, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (receiver *EphemeralScript) Update(model extension.Model, _ *gorm.DB, parameters gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	nameEphemeralScriptMapMutex.Lock()
	defer nameEphemeralScriptMapMutex.Unlock()

	ephemeralScript, exists := nameEphemeralScriptMap[name]

	if !exists {
		ephemeralScript = NewEphemeralScript()
		if err := mapstructutilpkg.GetUtility().MapToStruct(inputContainer, ephemeralScript); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		ephemeralScript.Status = ephemeralScriptStatusCreated
		ephemeralScript.ReturnValue = ""
		ephemeralScript.Error = ""
		ephemeralScript.StartedAt = ""
		ephemeralScript.FinishedAt = ""
		ephemeralScript.vm = otto.New()
		ephemeralScript.vm.Interrupt = make(chan func(), 1)
		nameEphemeralScriptMap[name] = ephemeralScript
	} else {
		if ephemeralScript.Status == ephemeralScriptStatusRunning {
			logging.Logger().Debug(errors.New("ephemeralScript is running"))
			return nil, errors.New("ephemeralScript is running")
		}

		newEephemeralScript := NewEphemeralScript()
		if err := mapstructutilpkg.GetUtility().MapToStruct(inputContainer, newEephemeralScript); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}

		// overwrite properties except for otto vm
		ephemeralScript.Status = ephemeralScriptStatusCreated
		ephemeralScript.Description = newEephemeralScript.Description
		ephemeralScript.ScriptContent = newEephemeralScript.ScriptContent
		ephemeralScript.ReturnValue = ""
		ephemeralScript.Error = ""
		ephemeralScript.StartedAt = ""
		ephemeralScript.FinishedAt = ""
	}

	return ephemeralScript, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *EphemeralScript) Delete(model extension.Model, _ *gorm.DB, parameters gin.Params, _ url.Values) error {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	nameEphemeralScriptMapMutex.Lock()
	defer nameEphemeralScriptMapMutex.Unlock()

	ephemeralScript, exists := nameEphemeralScriptMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return errors.New("record not found")
	}

	if ephemeralScript.Status == ephemeralScriptStatusRunning {
		logging.Logger().Debug(errors.New("ephemeralScript is running"))
		return errors.New("ephemeralScript is running")
	}

	close(ephemeralScript.vm.Interrupt)

	delete(nameEphemeralScriptMap, name)

	return nil
}

// GetCount returns the record count under current db conditions
func (receiver *EphemeralScript) GetCount(_ extension.Model, _ *gorm.DB) (int, error) {
	return len(nameEphemeralScriptMap), nil
}

func init() {
	extension.RegisterModel(NewEphemeralScript())
}
