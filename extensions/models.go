package extensions

import (
	"fmt"
	"reflect"
)

var typeMap = map[string]reflect.Type{}
var modelMap = map[string]interface{}{}

// RegisterModel registers a model to migrate automatically, and to generate new instances in processing requests
func RegisterModel(name string, model interface{}) {
	reflectType := reflect.TypeOf(model)
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}
	typeMap[name] = reflectType
	modelMap[name] = reflect.New(reflectType).Elem().Interface()
}

// RegisteredModels returns the registered models
func RegisteredModels() []interface{} {
	result := []interface{}{}
	for _, model := range modelMap {
		result = append(result, model)
	}
	return result
}

// CreateModel creates a model instance using given name and the registered model related to the name
func CreateModel(name string) (interface{}, error) {
	reflectType, exists := typeMap[name]
	if !exists {
		return nil, fmt.Errorf("the type named '%s' has not been registered yet", name)
	}
	return reflect.New(reflectType).Elem().Addr().Interface(), nil
}
