package model

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	collectionutilpkg "github.com/qb0C80aE/clay/util/collection"
	conversionutilpkg "github.com/qb0C80aE/clay/util/conversion"
	executilpkg "github.com/qb0C80aE/clay/util/exec"
	ioutilpkg "github.com/qb0C80aE/clay/util/io"
	loggingutilpkg "github.com/qb0C80aE/clay/util/logging"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
	modelstorepkg "github.com/qb0C80aE/clay/util/modelstore"
	networkutilpkg "github.com/qb0C80aE/clay/util/network"
	stringutilpkg "github.com/qb0C80aE/clay/util/string"
	"github.com/robertkrimen/otto"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// EphemeralScriptExecution is the model class what represents ephemeralScript execution
type EphemeralScriptExecution struct {
	Base
	Name string                 `json:"name" yaml:"name" form:"name" gorm:"primary_key"`
	Data map[string]interface{} `json:"data" yaml:"data" form:"data" sql:"-"`
}

// NewEphemeralScriptExecution creates a ephemeralScriptExecution model instance
func NewEphemeralScriptExecution() *EphemeralScriptExecution {
	return &EphemeralScriptExecution{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralScriptExecution) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *EphemeralScriptExecution) Create(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, inputContainer interface{}) (interface{}, error) {
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
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	if ephemeralScript.Status == ephemeralScriptStatusRunning {
		logging.Logger().Debug(errors.New("the ephemeralScript is already running"))
		return nil, errors.New("the ephemeralScript is already running")
	}

	ephemeralScript.Status = ephemeralScriptStatusRunning
	ephemeralScript.ReturnValue = ""
	ephemeralScript.Error = ""
	ephemeralScript.StartedAt = time.Now().String()
	ephemeralScript.FinishedAt = ""
	ephemeralScript.vm.Set("Data", reflect.ValueOf(inputContainer).Elem().FieldByName("Data").Interface())
	ephemeralScript.vm.Set("ModelStore", modelstorepkg.NewModelStore(db))
	ephemeralScript.vm.Set("Query", urlValues)
	ephemeralScript.vm.Set("Exec", executilpkg.GetUtility())
	ephemeralScript.vm.Set("IO", ioutilpkg.GetUtility())
	ephemeralScript.vm.Set("Collection", collectionutilpkg.GetUtility())
	ephemeralScript.vm.Set("Conversion", conversionutilpkg.GetUtility())
	ephemeralScript.vm.Set("MapStruct", mapstructutilpkg.GetUtility())
	ephemeralScript.vm.Set("Network", networkutilpkg.GetUtility())
	ephemeralScript.vm.Set("String", stringutilpkg.GetUtility())
	ephemeralScript.vm.Set("Logging", loggingutilpkg.GetUtility())
	ephemeralScript.vm.Set("error", func(message interface{}) {
		errorMessage := fmt.Sprintf("an error has occured during execution of %s: %s at line %d",
			ephemeralScript.Name,
			message,
			ephemeralScript.vm.Context().Line)
		panic(errorMessage)
	})

	if urlValues.Get("execution_mode") == "sync" {
		executeEphemeralScript(ephemeralScript)
		if len(ephemeralScript.Error) != 0 {
			return nil, errors.New(ephemeralScript.Error)
		}
	} else {
		go executeEphemeralScript(ephemeralScript)
	}

	return ephemeralScript, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *EphemeralScriptExecution) Delete(model extension.Model, _ *gorm.DB, parameters gin.Params, _ url.Values) error {
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

	if ephemeralScript.Status != ephemeralScriptStatusRunning {
		logging.Logger().Debug(errors.New("the ephemeralScript is already running"))
		return errors.New("the ephemeralScript is not running")
	}

	ephemeralScript.vm.Interrupt <- func() {
		panic("script has been killed")
	}

	ephemeralScript.Status = ephemeralScriptStatusFinished
	ephemeralScript.FinishedAt = time.Now().String()

	return nil
}

func executeEphemeralScript(ephemeralScript *EphemeralScript) {
	defer func() {
		if err := recover(); err != nil {
			ephemeralScript.ReturnValue = ""
			ephemeralScript.Error = fmt.Sprintf("%v", err)
			ephemeralScript.FinishedAt = time.Now().String()
			ephemeralScript.Status = ephemeralScriptStatusFinished

			logging.Logger().Debug(fmt.Sprintf("%v", err))
		}
	}()

	scriptContent := fmt.Sprintf(`(function() { %s })();`, ephemeralScript.ScriptContent)
	value, err := ephemeralScript.vm.Run(scriptContent)

	if err != nil {
		logging.Logger().Debugf("an error has occured during execution of %s: ", ephemeralScript.Name)
		if _, ok := err.(*otto.Error); ok {
			err := err.(*otto.Error)
			errStringList := strings.Split(err.String(), "\n")
			for _, errString := range errStringList {
				logging.Logger().Debug(errString)
			}
			ephemeralScript.Error = fmt.Sprintf("%v", err.String())
		} else {
			logging.Logger().Debug(err.Error())
			ephemeralScript.Error = fmt.Sprintf("%v", err.Error())
		}
	} else {
		ephemeralScript.Error = ""
	}
	ephemeralScript.ReturnValue, _ = value.ToString()

	ephemeralScript.FinishedAt = time.Now().String()
	ephemeralScript.Status = ephemeralScriptStatusFinished
}

func init() {
	extension.RegisterModel(NewEphemeralScriptExecution())
}
