package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
)

// EphemeralScriptRaw is the model class what represents ephemeral script raw
type EphemeralScriptRaw struct {
	Base
	Name string `json:"name" yaml:"name" gorm:"primary_key"`
}

// NewEphemeralScriptRaw creates a ephemeral script raw model instance
func NewEphemeralScriptRaw() *EphemeralScriptRaw {
	return &EphemeralScriptRaw{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralScriptRaw) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralScriptRaw) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	ephemeralScript, exists := nameEphemeralScriptMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	return ephemeralScript.ScriptContent, nil
}

func init() {
	extension.RegisterModel(NewEphemeralScriptRaw())
}
