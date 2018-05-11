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

var nameEphemeralTemplateMap = map[string]*EphemeralTemplate{}
var nameEphemeralTemplateMapMutex = new(sync.Mutex)

// EphemeralTemplate is the model class what represents ephemeral template data which can exist at runtime only
type EphemeralTemplate struct {
	Base
	Name            string `json:"name" yaml:"name" form:"name" gorm:"primary_key" validate:"required"`
	TemplateContent string `json:"template_content" yaml:"template_content" form:"template_content" sql:"type:text"`
	Description     string `json:"description" yaml:"description" form:"description" sql:"type:text"`
}

// NewEphemeralTemplate creates a ephemeralTemplate model instance
func NewEphemeralTemplate() *EphemeralTemplate {
	return &EphemeralTemplate{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralTemplate) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralTemplate) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	result, exists := nameEphemeralTemplateMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *EphemeralTemplate) GetMulti(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	nameEphemeralTemplateMapMutex.Lock()
	defer nameEphemeralTemplateMapMutex.Unlock()

	keyList := make([]string, 0, len(nameEphemeralTemplateMap))
	for key := range nameEphemeralTemplateMap {
		keyList = append(keyList, key)
	}

	sort.Strings(keyList)

	ephemeralTemplateList := make([]*EphemeralTemplate, len(nameEphemeralTemplateMap))

	for i, key := range keyList {
		ephemeralTemplateList[i] = nameEphemeralTemplateMap[key]
	}

	return ephemeralTemplateList, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *EphemeralTemplate) Create(_ extension.Model, _ *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	ephemeralTemplate := NewEphemeralTemplate()
	if err := mapstruct.RemapToStruct(inputContainer, ephemeralTemplate); err != nil {
		return nil, err
	}

	nameEphemeralTemplateMapMutex.Lock()
	defer nameEphemeralTemplateMapMutex.Unlock()

	if _, exists := nameEphemeralTemplateMap[ephemeralTemplate.Name]; exists {
		logging.Logger().Debugf("ephemeral template %s already exists", ephemeralTemplate.Name)
		return nil, fmt.Errorf("ephemeral template %s already exists", ephemeralTemplate.Name)
	}

	nameEphemeralTemplateMap[ephemeralTemplate.Name] = ephemeralTemplate

	return ephemeralTemplate, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (receiver *EphemeralTemplate) Update(model extension.Model, _ *gorm.DB, parameters gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	ephemeralTemplate := NewEphemeralTemplate()
	if err := mapstruct.RemapToStruct(inputContainer, ephemeralTemplate); err != nil {
		logging.Logger().Debug(err)
		return nil, err
	}

	ephemeralTemplate.Name = name

	nameEphemeralTemplateMapMutex.Lock()
	defer nameEphemeralTemplateMapMutex.Unlock()

	nameEphemeralTemplateMap[ephemeralTemplate.Name] = ephemeralTemplate

	return ephemeralTemplate, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *EphemeralTemplate) Delete(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values) error {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	nameEphemeralTemplateMapMutex.Lock()
	defer nameEphemeralTemplateMapMutex.Unlock()

	_, exists := nameEphemeralTemplateMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return errors.New("record not found")
	}

	delete(nameEphemeralTemplateMap, name)

	return nil
}

// GetCount returns the record count under current db conditions
func (receiver *EphemeralTemplate) GetCount(_ extension.Model, _ *gorm.DB) (int, error) {
	return len(nameEphemeralTemplateMap), nil
}

func init() {
	extension.RegisterModel(NewEphemeralTemplate())
}
