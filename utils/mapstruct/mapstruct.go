package mapstruct

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

// MapToStruct maps the map instance into the struct instance using JSON marshal logic
func MapToStruct(m []interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return err
	}
	return nil
}

// SliceToInterfaceSlice creates a []interface{} from a slice of concrete type
func SliceToInterfaceSlice(slice interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, errors.New("given argument is neither a slice nor an array")
	}

	size := v.Len()
	result := make([]interface{}, size)
	for i := 0; i < size; i++ {
		item := v.Index(i)

		itemForField := item
		for (itemForField.Kind() == reflect.Ptr) || (itemForField.Kind() == reflect.Interface) {
			itemForField = itemForField.Elem()
		}
		if !itemForField.IsValid() {
			return nil, fmt.Errorf("the actual item indexed %d in given slice is not valid", i)
		}
		if itemForField.Kind() != reflect.Struct {
			return nil, errors.New("given slice isn't a slice of struct")
		}

		if !item.CanInterface() {
			return nil, fmt.Errorf("the original item indexed %d in given slice cannot interface", i)
		}
		result[i] = item.Interface()
	}
	return result, nil
}

// SliceToInterfaceMap creates a map[interface{}]interface{} from a slice of a struct with key defined in the struct
// It doesn't modify the raw type of elements like strpping pointer or something
func SliceToInterfaceMap(structSlice interface{}, keyFieldName string) (map[interface{}]interface{}, error) {
	v := reflect.ValueOf(structSlice)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, errors.New("given argument is neither a slice nor an array")
	}
	size := v.Len()
	result := make(map[interface{}]interface{}, size)
	for i := 0; i < size; i++ {
		item := v.Index(i)

		itemForField := item
		for (itemForField.Kind() == reflect.Ptr) || (itemForField.Kind() == reflect.Interface) {
			itemForField = itemForField.Elem()
		}
		if !itemForField.IsValid() {
			return nil, fmt.Errorf("the actual item indexed %d in given slice is not valid", i)
		}
		if itemForField.Kind() != reflect.Struct {
			return nil, errors.New("given slice isn't a slice of struct")
		}
		field := itemForField.FieldByName(keyFieldName)
		if !field.IsValid() {
			return nil, fmt.Errorf("the actual struct in given slice doesn't have a field named '%s'", keyFieldName)
		}

		if !item.CanInterface() {
			return nil, fmt.Errorf("the original item indexed %d in given slice cannot interface", i)
		}
		result[field.Interface()] = item.Interface()
	}
	return result, nil
}

// SliceToInterfaceSliceMap creates a map[interface{}][]interface{} from a slice of a struct with key defined in the struct
// It doesn't modify the raw type of elements like strpping pointer or something
func SliceToInterfaceSliceMap(structSlice interface{}, keyFieldName string) (map[interface{}][]interface{}, error) {
	v := reflect.ValueOf(structSlice)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return nil, errors.New("given argument is neither a slice nor an array")
	}
	size := v.Len()
	result := map[interface{}][]interface{}{}
	for i := 0; i < size; i++ {
		item := v.Index(i)

		itemForField := item
		for (itemForField.Kind() == reflect.Ptr) || (itemForField.Kind() == reflect.Interface) {
			itemForField = itemForField.Elem()
		}
		if !itemForField.IsValid() {
			return nil, fmt.Errorf("the actual item indexed %d in given slice is not valid", i)
		}
		if itemForField.Kind() != reflect.Struct {
			return nil, errors.New("given slice isn't a slice of struct")
		}
		field := itemForField.FieldByName(keyFieldName)
		if !field.IsValid() {
			return nil, fmt.Errorf("the actual struct in given slice doesn't have a field named '%s'", keyFieldName)
		}

		if !item.CanInterface() {
			return nil, fmt.Errorf("the original item indexed %d in given slice cannot interface", i)
		}

		itemSlice, exists := result[field.Interface()]
		if exists {
			result[field.Interface()] = append(itemSlice, item.Interface())
		} else {
			result[field.Interface()] = []interface{}{item.Interface()}
		}
	}
	return result, nil
}
