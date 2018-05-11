package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
)

// EphemeralBinaryObjectRaw is the model class what represents ephemeral binary object raw
type EphemeralBinaryObjectRaw struct {
	Base
	Name string `json:"name" yaml:"name" gorm:"primary_key"`
}

// NewEphemeralBinaryObjectRaw creates a ephemeral binary object raw model instance
func NewEphemeralBinaryObjectRaw() *EphemeralBinaryObjectRaw {
	return &EphemeralBinaryObjectRaw{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralBinaryObjectRaw) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralBinaryObjectRaw) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	ephemeralBinaryObject, exists := nameEphemeralBinaryObjectMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	return ephemeralBinaryObject.Content, nil
}

func init() {
	extension.RegisterModel(NewEphemeralBinaryObjectRaw())
}
