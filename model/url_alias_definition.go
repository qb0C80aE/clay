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

var nameURLAliasDefinitionMap = map[string]*URLAliasDefinition{}
var nameURLAliasDefinitionMapMutex = new(sync.Mutex)

// URLAliasDefinition is the model class what represents url alias definitions
type URLAliasDefinition struct {
	Base
	Name  string `json:"name" clay:"key_parameter" validate:"required"`
	From  string `json:"from" validate:"required"`
	To    string `json:"to" validate:"required"`
	Query string `json:"query"`
}

// NewURLAliasDefinition creates a template raw model instance
func NewURLAliasDefinition() *URLAliasDefinition {
	return &URLAliasDefinition{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *URLAliasDefinition) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *URLAliasDefinition) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	result, exists := nameURLAliasDefinitionMap[parameters.ByName("name")]

	if !exists {
		logging.Logger().Debug("record not found")
		return nil, errors.New("record not found")
	}

	return result, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *URLAliasDefinition) GetMulti(_ extension.Model, _ *gorm.DB, _ gin.Params, _ url.Values, _ string) (interface{}, error) {
	nameURLAliasDefinitionMapMutex.Lock()
	defer nameURLAliasDefinitionMapMutex.Unlock()

	keyList := make([]string, 0, len(nameURLAliasDefinitionMap))
	for key := range nameURLAliasDefinitionMap {
		keyList = append(keyList, key)
	}

	sort.Strings(keyList)

	urlAliasDefinitionList := make([]*URLAliasDefinition, len(nameURLAliasDefinitionMap))

	for i, key := range keyList {
		urlAliasDefinitionList[i] = nameURLAliasDefinitionMap[key]
	}

	return urlAliasDefinitionList, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *URLAliasDefinition) Create(_ extension.Model, _ *gorm.DB, _ gin.Params, _ url.Values, inputContainer interface{}) (interface{}, error) {
	nameURLAliasDefinitionMapMutex.Lock()
	defer nameURLAliasDefinitionMapMutex.Unlock()

	urlAliasDefinition := NewURLAliasDefinition()
	if err := mapstruct.RemapToStruct(inputContainer, urlAliasDefinition); err != nil {
		return nil, err
	}

	if _, exists := nameURLAliasDefinitionMap[urlAliasDefinition.Name]; exists {
		logging.Logger().Debugf("alias %s already exists", urlAliasDefinition.Name)
		return nil, fmt.Errorf("alias %s already exists", urlAliasDefinition.Name)
	}

	nameURLAliasDefinitionMap[urlAliasDefinition.Name] = urlAliasDefinition

	return urlAliasDefinition, nil
}

// GetTotal returns the count of for multi resource
func (receiver *URLAliasDefinition) GetTotal(_ extension.Model, _ *gorm.DB) (int, error) {
	return len(nameURLAliasDefinitionMap), nil
}

func init() {
	extension.RegisterModel(NewURLAliasDefinition())
}
