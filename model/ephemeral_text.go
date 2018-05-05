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
	"sort"
	"sync"
)

var nameEphemeralTextMap = map[string]*EphemeralText{}
var nameEphemeralTextMapMutex = new(sync.Mutex)

// EphemeralText is the model class what represents ephemeral text data which can exist at runtime only
type EphemeralText struct {
	Base
	Name        string `json:"name" form:"name" gorm:"primary_key" clay:"key_parameter" validate:"required"`
	Content     string `json:"content" form:"content" sql:"type:text"`
	Description string `json:"description" form:"description" sql:"type:text"`
}

// NewEphemeralText creates a ephemeralText model instance
func NewEphemeralText() *EphemeralText {
	return &EphemeralText{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralText) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralText) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	result, exists := nameEphemeralTextMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *EphemeralText) GetMulti(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	nameEphemeralTextMapMutex.Lock()
	defer nameEphemeralTextMapMutex.Unlock()

	keyList := make([]string, 0, len(nameEphemeralTextMap))
	for key := range nameEphemeralTextMap {
		keyList = append(keyList, key)
	}

	sort.Strings(keyList)

	ephemeralTextList := make([]*EphemeralText, len(nameEphemeralTextMap))

	for i, key := range keyList {
		ephemeralTextList[i] = nameEphemeralTextMap[key]
	}

	return ephemeralTextList, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *EphemeralText) Create(_ extension.Model, _ *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	ephemeralText := NewEphemeralText()
	if err := mapstruct.RemapToStruct(inputContainer, ephemeralText); err != nil {
		return nil, err
	}

	nameEphemeralTextMapMutex.Lock()
	defer nameEphemeralTextMapMutex.Unlock()

	if _, exists := nameEphemeralTextMap[ephemeralText.Name]; exists {
		logging.Logger().Debugf("ephemeral text %s already exists", ephemeralText.Name)
		return nil, fmt.Errorf("ephemeral text %s already exists", ephemeralText.Name)
	}

	nameEphemeralTextMap[ephemeralText.Name] = ephemeralText

	return ephemeralText, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (receiver *EphemeralText) Update(model extension.Model, _ *gorm.DB, parameters gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	ephemeralText := NewEphemeralText()
	if err := mapstruct.RemapToStruct(inputContainer, ephemeralText); err != nil {
		logging.Logger().Debug(err)
		return nil, err
	}

	ephemeralText.Name = name

	nameEphemeralTextMapMutex.Lock()
	defer nameEphemeralTextMapMutex.Unlock()

	nameEphemeralTextMap[ephemeralText.Name] = ephemeralText

	return ephemeralText, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *EphemeralText) Delete(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values) error {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	nameEphemeralTextMapMutex.Lock()
	defer nameEphemeralTextMapMutex.Unlock()

	_, exists := nameEphemeralTextMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return errors.New("record not found")
	}

	delete(nameEphemeralTextMap, name)

	return nil
}

// GetTotal returns the count of for multi resource
func (receiver *EphemeralText) GetTotal(_ extension.Model, _ *gorm.DB) (int, error) {
	return len(nameEphemeralTextMap), nil
}

func init() {
	extension.RegisterModel(NewEphemeralText())
}
