package extension

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
	"reflect"
	"strings"
)

var modelList = []Model{}
var modelToBeMigratedList = []Model{}
var resourceNameModelMap = map[string]Model{}
var typeNameResourceNameMap = map[string]string{}
var typeNameModelKeyMap = map[string]ModelKey{}

// Model is the interface what contains the actual model data and handle it
// * NewModelContainer creates a container, which means a model without Base
// * IsContainer tells if the instance is container, that is the model instance without Base-ActualModel relationship
// * GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
// * GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
// * Create corresponds HTTP POST message and handles a request for multi resource to create a new information
// * Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
// * Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
// * Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
// * GetOptions corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
// * GetTotal returns the count of for multi resource
type Model interface {
	New() Model
	IsContainer() bool
	GetSingle(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryString string) (interface{}, error)
	GetMulti(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryString string) (interface{}, error)
	Create(db *gorm.DB, parameters gin.Params, urlValues url.Values, input Model) (interface{}, error)
	Update(db *gorm.DB, parameters gin.Params, urlValues url.Values, input Model) (interface{}, error)
	Delete(db *gorm.DB, parameters gin.Params, urlValues url.Values) error
	Patch(db *gorm.DB, parameters gin.Params, urlValues url.Values, input Model) (interface{}, error)
	GetOptions(db *gorm.DB, parameters gin.Params, urlValues url.Values) error
	GetTotal(db *gorm.DB) (int, error)
}

// ModelKey is the type that has keys used in clay for "to_be_deleted" (delete specific children in update) logic or any other key replacement logic.
type ModelKey struct {
	KeyParameter string
	KeyField     string
}

// GetActualType returns the type of given object stripping pointer and interface
func GetActualType(object interface{}) reflect.Type {
	objectType := reflect.TypeOf(object)
	for (objectType.Kind() == reflect.Ptr) || (objectType.Kind() == reflect.Interface) {
		objectType = objectType.Elem()
	}
	return objectType
}

// GetActualValue returns the value of given object stripping pointer and interface
func GetActualValue(object interface{}) reflect.Value {
	objectValue := reflect.ValueOf(object)
	for (objectValue.Kind() == reflect.Ptr) || (objectValue.Kind() == reflect.Interface) {
		objectValue = objectValue.Elem()
	}
	return objectValue
}

// RegisterModel registers a model to migrate automatically, and to generate new instances in processing requests
func RegisterModel(model Model, autoMigration bool) {
	if model.IsContainer() {
		panic(errors.New("the given model is a container"))
	}

	modelList = append(modelList, model)
	if autoMigration {
		modelToBeMigratedList = append(modelToBeMigratedList, model)
	}

	modelType := GetActualType(model)
	modelTypeName := modelType.String()
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag, ok := field.Tag.Lookup("clay")
		if ok {
			tagStatementList := strings.Split(tag, ";")
			for _, tagStatement := range tagStatementList {
				tagKeyValue := strings.Split(tagStatement, ":")
				switch tagKeyValue[0] {
				case "key_parameter":
					typeNameModelKeyMap[modelTypeName] = ModelKey{
						KeyParameter: tagKeyValue[1],
						KeyField:     field.Name,
					}
					break
				}
			}
		}
		if _, exists := typeNameModelKeyMap[modelTypeName]; exists {
			break
		}
	}

	if _, exists := typeNameModelKeyMap[modelTypeName]; !exists {
		typeNameModelKeyMap[modelTypeName] = ModelKey{
			KeyParameter: "id",
			KeyField:     "ID",
		}
	}
}

// GetRegisteredModelList returns the registered models
func GetRegisteredModelList() []Model {
	result := make([]Model, len(modelList))
	for i, model := range modelList {
		result[i] = model
	}
	return result
}

// GetRegisteredModelKey returns the registered model key
func GetRegisteredModelKey(model Model) (ModelKey, error) {
	modelType := GetActualType(model)
	modelTypeName := modelType.String()
	result, exists := typeNameModelKeyMap[modelTypeName]
	if !exists {
		logging.Logger().Debugf("the model key of '%s' has not been registered yet", modelTypeName)
		return ModelKey{}, fmt.Errorf("the model key of '%s' has not been registered yet", modelTypeName)
	}
	return result, nil
}

// GetRegisteredModelToBeMigratedList returns the registered models to be migrated
func GetRegisteredModelToBeMigratedList() []interface{} {
	result := make([]interface{}, len(modelToBeMigratedList))
	for i, model := range modelToBeMigratedList {
		result[i] = model
	}
	return result
}

// AssociateResourceNameWithModel registers a name of given model
func AssociateResourceNameWithModel(resourceName string, model Model) {
	if model.IsContainer() {
		panic(errors.New("the given model is a container"))
	}

	modelType := GetActualType(model)
	modelTypeName := modelType.String()
	typeNameResourceNameMap[modelTypeName] = resourceName
	resourceNameModelMap[resourceName] = model
}

// GetAssociateResourceNameWithModel returns the registered resource name of given model
func GetAssociateResourceNameWithModel(model Model) string {
	modelType := GetActualType(model)
	return typeNameResourceNameMap[modelType.String()]
}

// CreateModelByResourceName creates a model instance using given resource name and the registered model related to the resource name
func CreateModelByResourceName(resourceName string) (Model, error) {
	model, exists := resourceNameModelMap[resourceName]
	if !exists {
		logging.Logger().Debugf("the type which associated with the resource name '%s' has not been registered yet", resourceName)
		return nil, fmt.Errorf("the type named associated with the resource name '%s' has not been registered yet", resourceName)
	}
	return model.New(), nil
}
