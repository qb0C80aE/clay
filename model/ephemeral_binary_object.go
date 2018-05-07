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

var nameEphemeralBinaryObjectMap = map[string]*EphemeralBinaryObject{}
var nameEphemeralBinaryObjectMapMutex = new(sync.Mutex)

// EphemeralBinaryObject is the model class what represents ephemeral binary object data which can exist at runtime only
type EphemeralBinaryObject struct {
	Base
	Name        string `json:"name" form:"name" gorm:"primary_key" clay:"key_parameter" validate:"required"`
	Content     []byte `json:"content" form:"content"`
	Description string `json:"description" form:"description" sql:"type:text"`
}

// NewEphemeralBinaryObject creates a ephemeralBinaryObject model instance
func NewEphemeralBinaryObject() *EphemeralBinaryObject {
	return &EphemeralBinaryObject{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralBinaryObject) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralBinaryObject) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	result, exists := nameEphemeralBinaryObjectMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	snipped := &EphemeralBinaryObject{
		Name:        result.Name,
		Description: result.Description,
		Content:     []byte{},
	}

	return snipped, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *EphemeralBinaryObject) GetMulti(_ extension.Model, db *gorm.DB, _ gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	nameEphemeralBinaryObjectMapMutex.Lock()
	defer nameEphemeralBinaryObjectMapMutex.Unlock()

	keyList := make([]string, 0, len(nameEphemeralBinaryObjectMap))
	for key := range nameEphemeralBinaryObjectMap {
		keyList = append(keyList, key)
	}

	sort.Strings(keyList)

	ephemeralBinaryObjectList := make([]*EphemeralBinaryObject, len(nameEphemeralBinaryObjectMap))

	for i, key := range keyList {
		item := nameEphemeralBinaryObjectMap[key]
		snipped := &EphemeralBinaryObject{
			Name:        item.Name,
			Description: item.Description,
			Content:     []byte{},
		}
		ephemeralBinaryObjectList[i] = snipped
	}

	return ephemeralBinaryObjectList, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *EphemeralBinaryObject) Create(_ extension.Model, _ *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	ephemeralBinaryObject := NewEphemeralBinaryObject()

	if err := mapstruct.RemapToStruct(inputContainer, ephemeralBinaryObject); err != nil {
		return nil, err
	}

	nameEphemeralBinaryObjectMapMutex.Lock()
	defer nameEphemeralBinaryObjectMapMutex.Unlock()

	if _, exists := nameEphemeralBinaryObjectMap[ephemeralBinaryObject.Name]; exists {
		logging.Logger().Debugf("ephemeral binary object %s already exists", ephemeralBinaryObject.Name)
		return nil, fmt.Errorf("ephemeral binary object %s already exists", ephemeralBinaryObject.Name)
	}

	nameEphemeralBinaryObjectMap[ephemeralBinaryObject.Name] = ephemeralBinaryObject

	snipped := &EphemeralBinaryObject{
		Name:        ephemeralBinaryObject.Name,
		Description: ephemeralBinaryObject.Description,
		Content:     []byte{},
	}

	return snipped, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (receiver *EphemeralBinaryObject) Update(model extension.Model, _ *gorm.DB, parameters gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	ephemeralBinaryObject := NewEphemeralBinaryObject()
	if err := mapstruct.RemapToStruct(inputContainer, ephemeralBinaryObject); err != nil {
		logging.Logger().Debug(err)
		return nil, err
	}

	ephemeralBinaryObject.Name = name

	nameEphemeralBinaryObjectMapMutex.Lock()
	defer nameEphemeralBinaryObjectMapMutex.Unlock()

	nameEphemeralBinaryObjectMap[ephemeralBinaryObject.Name] = ephemeralBinaryObject

	snipped := &EphemeralBinaryObject{
		Name:        ephemeralBinaryObject.Name,
		Description: ephemeralBinaryObject.Description,
		Content:     []byte{},
	}

	return snipped, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *EphemeralBinaryObject) Delete(model extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values) error {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	nameEphemeralBinaryObjectMapMutex.Lock()
	defer nameEphemeralBinaryObjectMapMutex.Unlock()

	_, exists := nameEphemeralBinaryObjectMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return errors.New("record not found")
	}

	delete(nameEphemeralBinaryObjectMap, name)

	return nil
}

// GetTotal returns the count of for multi resource
func (receiver *EphemeralBinaryObject) GetTotal(_ extension.Model, _ *gorm.DB) (int, error) {
	return len(nameEphemeralBinaryObjectMap), nil
}

func init() {
	extension.RegisterModel(NewEphemeralBinaryObject())
}
