package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
)

// EphemeralTemplateRaw is the model class what represents ephemeral template raw
type EphemeralTemplateRaw struct {
	Base
	Name string `json:"name" yaml:"name" clay:"key_parameter"`
}

// NewEphemeralTemplateRaw creates a ephemeral template raw model instance
func NewEphemeralTemplateRaw() *EphemeralTemplateRaw {
	return &EphemeralTemplateRaw{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralTemplateRaw) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralTemplateRaw) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	ephemeralTemplate, exists := nameEphemeralTemplateMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	return ephemeralTemplate.TemplateContent, nil
}

func init() {
	extension.RegisterModel(NewEphemeralTemplateRaw())
}
