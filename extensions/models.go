package extensions

import "reflect"

var typeMap = map[string]reflect.Type{}
var modelMap = map[reflect.Type]interface{}{}

func RegisterModelType(model interface{}) {
	reflectType := reflect.TypeOf(model)
	typeMap[reflectType.String()] = reflectType
	modelMap[reflectType] = reflect.New(reflectType).Elem().Interface()
}

func GetModels() []interface{} {
	result := []interface{}{}
	for _, model := range modelMap {
		result = append(result, model)
	}
	return result
}
