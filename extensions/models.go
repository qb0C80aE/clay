package extensions

import "reflect"

var typeMap = map[string]reflect.Type{}
var modelMap = map[reflect.Type]interface{}{}

// RegisterModel registers a model to migrate automatically, and to generate new instances in processing requests
func RegisterModel(model interface{}) {
	reflectType := reflect.TypeOf(model)
	typeMap[reflectType.String()] = reflectType
	modelMap[reflectType] = reflect.New(reflectType).Elem().Interface()
}

// RegisteredModels returns the registered models
func RegisteredModels() []interface{} {
	result := []interface{}{}
	for _, model := range modelMap {
		result = append(result, model)
	}
	return result
}
