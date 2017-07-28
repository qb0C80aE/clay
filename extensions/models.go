package extensions

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
)

var models = []interface{}{}
var typeMap = map[string]reflect.Type{}
var modelMap = map[string]interface{}{}
var resourceNameMap = map[reflect.Type]string{}
var initialDataLoaders = []InitialDataLoader{}

// InitialDataLoader is the interface what setups the initial data
// * SetupInitialData setups the initial data
type InitialDataLoader interface {
	SetupInitialData(db *gorm.DB) error
}

// ModelType returns the type of given model stripping pointer and interface
func ModelType(model interface{}) reflect.Type {
	reflectType := reflect.TypeOf(model)
	for (reflectType.Kind() == reflect.Ptr) || (reflectType.Kind() == reflect.Interface) {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

// RegisterModel registers a model to migrate automatically, and to generate new instances in processing requests
func RegisterModel(model interface{}) {
	reflectType := ModelType(model)
	models = append(models, reflect.New(reflectType).Elem().Interface())
}

// RegisteredModels returns the registered models
func RegisteredModels() []interface{} {
	result := make([]interface{}, len(models))
	for i, model := range models {
		result[i] = model
	}
	return result
}

// RegisterResourceName registers a name of given model
func RegisterResourceName(model interface{}, name string) {
	reflectType := ModelType(model)
	resourceNameMap[reflectType] = name
	typeMap[name] = reflectType
	modelMap[name] = reflect.New(reflectType).Elem().Interface()
}

// RegisteredResourceName returns the registered resource name of given model
func RegisteredResourceName(model interface{}) string {
	reflectType := ModelType(model)
	return resourceNameMap[reflectType]
}

// CreateModel creates a model instance using given name and the registered model related to the name
func CreateModel(name string) (interface{}, error) {
	reflectType, exists := typeMap[name]
	if !exists {
		return nil, fmt.Errorf("the type named '%s' has not been registered yet", name)
	}
	return reflect.New(reflectType).Elem().Addr().Interface(), nil
}

// RegisterInitialDataLoader registers an initial data loader
func RegisterInitialDataLoader(initialDataLoader InitialDataLoader) {
	initialDataLoaders = append(initialDataLoaders, initialDataLoader)
}

// RegisteredInitialDataLoaders returns the registered initial data loaders
func RegisteredInitialDataLoaders() []InitialDataLoader {
	result := make([]InitialDataLoader, len(initialDataLoaders))
	for i, initialDataLoader := range initialDataLoaders {
		result[i] = initialDataLoader
	}
	return result
}
