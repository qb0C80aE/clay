package mapstruct

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qb0C80aE/clay/logging"
	"reflect"
	"sort"
)

// MapToStruct maps the map instance into the struct instance using JSON marshal logic
func MapToStruct(m []interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}
	return nil
}

// SliceToInterfaceSlice creates a []interface{} from a slice of concrete type
// It doesn't modify the raw type of elements like strpping pointer or something
func SliceToInterfaceSlice(slice interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(slice)

	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		logging.Logger().Debug("given argument is neither a slice nor an array")
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
			logging.Logger().Debugf("the actual item indexed %d in given slice is not valid", i)
			return nil, fmt.Errorf("the actual item indexed %d in given slice is not valid", i)
		}
		// This method does not require the field name of struct
		//if itemForField.Kind() != reflect.Struct {
		//	return nil, errors.New("given slice isn't a slice of struct")
		//}

		if !item.CanInterface() {
			logging.Logger().Debugf("the original item indexed %d in given slice cannot interface", i)
			return nil, fmt.Errorf("the original item indexed %d in given slice cannot interface", i)
		}
		result[i] = item.Interface()
	}
	return result, nil
}

// StructSliceToInterfaceMap creates a map[interface{}]interface{} from a slice of a struct with key defined in the struct
// It doesn't modify the raw type of elements like strpping pointer or something
func StructSliceToInterfaceMap(structSlice interface{}, keyFieldName string) (map[interface{}]interface{}, error) {
	structSliceValue := reflect.ValueOf(structSlice)

	if structSliceValue.Kind() != reflect.Slice && structSliceValue.Kind() != reflect.Array {
		logging.Logger().Debug("given argument is neither a slice nor an array")
		return nil, errors.New("given argument is neither a slice nor an array")
	}

	size := structSliceValue.Len()
	result := make(map[interface{}]interface{}, size)
	for i := 0; i < size; i++ {
		sliceElementValue := structSliceValue.Index(i)

		sliceElementValueForField := sliceElementValue
		for (sliceElementValueForField.Kind() == reflect.Ptr) || (sliceElementValueForField.Kind() == reflect.Interface) {
			sliceElementValueForField = sliceElementValueForField.Elem()
		}
		if !sliceElementValueForField.IsValid() {
			logging.Logger().Debugf("the actual item indexed %d in given slice is not valid", i)
			return nil, fmt.Errorf("the actual item indexed %d in given slice is not valid", i)
		}
		// This method requires the field name of struct
		if sliceElementValueForField.Kind() != reflect.Struct {
			logging.Logger().Debug("given slice isn't a slice of struct")
			return nil, errors.New("given slice isn't a slice of struct")
		}
		structElementFieldValue := sliceElementValueForField.FieldByName(keyFieldName)
		if !structElementFieldValue.IsValid() {
			logging.Logger().Debugf("the actual struct in given slice doesn't have a field named '%s'", keyFieldName)
			return nil, fmt.Errorf("the actual struct in given slice doesn't have a field named '%s'", keyFieldName)
		}

		if !sliceElementValue.CanInterface() {
			logging.Logger().Debugf("the original item indexed %d in given slice cannot interface", i)
			return nil, fmt.Errorf("the original item indexed %d in given slice cannot interface", i)
		}
		result[structElementFieldValue.Interface()] = sliceElementValue.Interface()
	}
	return result, nil
}

// StructSliceToInterfaceSliceMap creates a map[interface{}][]interface{} from a slice of a struct with key defined in the struct
// It doesn't modify the raw type of elements like strpping pointer or something
func StructSliceToInterfaceSliceMap(structSlice interface{}, keyFieldName string) (map[interface{}][]interface{}, error) {
	structSliceValue := reflect.ValueOf(structSlice)

	if structSliceValue.Kind() != reflect.Slice && structSliceValue.Kind() != reflect.Array {
		logging.Logger().Debug("given argument is neither a slice nor an array")
		return nil, errors.New("given argument is neither a slice nor an array")
	}

	size := structSliceValue.Len()
	result := map[interface{}][]interface{}{}
	for i := 0; i < size; i++ {
		sliceElementValue := structSliceValue.Index(i)

		sliceElementValueForField := sliceElementValue
		for (sliceElementValueForField.Kind() == reflect.Ptr) || (sliceElementValueForField.Kind() == reflect.Interface) {
			sliceElementValueForField = sliceElementValueForField.Elem()
		}
		if !sliceElementValueForField.IsValid() {
			logging.Logger().Debugf("the actual item indexed %d in given slice is not valid", i)
			return nil, fmt.Errorf("the actual item indexed %d in given slice is not valid", i)
		}
		// This method requires the field name of struct
		if sliceElementValueForField.Kind() != reflect.Struct {
			logging.Logger().Debug("given slice isn't a slice of struct")
			return nil, errors.New("given slice isn't a slice of struct")
		}
		structElementFieldValue := sliceElementValueForField.FieldByName(keyFieldName)
		if !structElementFieldValue.IsValid() {
			logging.Logger().Debugf("the actual struct in given slice doesn't have a field named '%s'", keyFieldName)
			return nil, fmt.Errorf("the actual struct in given slice doesn't have a field named '%s'", keyFieldName)
		}

		if !sliceElementValue.CanInterface() {
			logging.Logger().Debugf("the original item indexed %d in given slice cannot interface", i)
			return nil, fmt.Errorf("the original item indexed %d in given slice cannot interface", i)
		}

		itemSlice, exists := result[structElementFieldValue.Interface()]
		if exists {
			result[structElementFieldValue.Interface()] = append(itemSlice, sliceElementValue.Interface())
		} else {
			result[structElementFieldValue.Interface()] = []interface{}{sliceElementValue.Interface()}
		}
	}
	return result, nil
}

// StructSliceToFieldValueInterfaceSlice creates a []interface{} from a slice of a struct with field defined in the struct
// It doesn't modify the raw type of elements like strpping pointer or something
func StructSliceToFieldValueInterfaceSlice(structSlice interface{}, fieldName string) ([]interface{}, error) {
	structSliceValue := reflect.ValueOf(structSlice)

	if structSliceValue.Kind() != reflect.Slice && structSliceValue.Kind() != reflect.Array {
		logging.Logger().Debug("given argument is neither a slice nor an array")
		return nil, errors.New("given argument is neither a slice nor an array")
	}

	size := structSliceValue.Len()
	result := make([]interface{}, size)
	for i := 0; i < size; i++ {
		sliceElementValue := structSliceValue.Index(i)

		sliceElementValueForField := sliceElementValue
		for (sliceElementValueForField.Kind() == reflect.Ptr) || (sliceElementValueForField.Kind() == reflect.Interface) {
			sliceElementValueForField = sliceElementValueForField.Elem()
		}
		if !sliceElementValueForField.IsValid() {
			logging.Logger().Debugf("the actual item indexed %d in given slice is not valid", i)
			return nil, fmt.Errorf("the actual item indexed %d in given slice is not valid", i)
		}
		// This method requires the field name of struct
		if sliceElementValueForField.Kind() != reflect.Struct {
			logging.Logger().Debug("given slice isn't a slice of struct")
			return nil, errors.New("given slice isn't a slice of struct")
		}
		structElementFieldValue := sliceElementValueForField.FieldByName(fieldName)
		if !structElementFieldValue.IsValid() {
			logging.Logger().Debugf("the actual struct in given slice doesn't have a field named '%s'", fieldName)
			return nil, fmt.Errorf("the actual struct in given slice doesn't have a field named '%s'", fieldName)
		}
		if !structElementFieldValue.CanInterface() {
			logging.Logger().Debugf("the original field of item indexed %d in given slice cannot interface", i)
			return nil, fmt.Errorf("the original field of item indexed %d in given slice cannot interface", i)
		}
		result[i] = structElementFieldValue.Interface()
	}
	return result, nil
}

// MapToKeySlice creates a []interface{} from keys of map
// It doesn't modify the raw type of elements like strpping pointer or something
func MapToKeySlice(mapInterface interface{}) ([]interface{}, error) {
	mapInterfaceValue := reflect.ValueOf(mapInterface)

	if mapInterfaceValue.Kind() != reflect.Map {
		logging.Logger().Debug("given argument is not a map")
		return nil, errors.New("given argument is not a map")
	}

	keys := mapInterfaceValue.MapKeys()

	size := len(keys)
	result := make([]interface{}, size)

	for i := 0; i < size; i++ {
		key := keys[i]
		for (key.Kind() == reflect.Ptr) || (key.Kind() == reflect.Interface) {
			key = key.Elem()
		}
		if !key.IsValid() {
			logging.Logger().Debug("the key in given map is not valid")
			return nil, fmt.Errorf("the key in given map is not valid")
		}
		if !key.CanInterface() {
			logging.Logger().Debug("the key in given map cannot interface")
			return nil, fmt.Errorf("the key in given map cannot interface")
		}
		result[i] = keys[i].Interface()
	}

	return result, nil
}

// SortSlice sorts a slice consists of digit type or string type
// It doesn't modify the raw type of elements like strpping pointer or something
// The elements in the given slice must be the same type of digit primitive or string
// The order must be "asc" or "desc"
func SortSlice(slice interface{}, order string) ([]interface{}, error) {
	sliceValue := reflect.ValueOf(slice)

	if sliceValue.Kind() != reflect.Slice && sliceValue.Kind() != reflect.Array {
		logging.Logger().Debug("given argument is neither a slice nor an array")
		return nil, errors.New("given argument is neither a slice nor an array")
	}

	size := sliceValue.Len()
	result := make([]interface{}, size)
	if size == 0 {
		return result, nil
	}

	var kindExpected reflect.Kind
	for i := 0; i < size; i++ {
		sliceElementValue := sliceValue.Index(i)
		for (sliceElementValue.Kind() == reflect.Ptr) || (sliceElementValue.Kind() == reflect.Interface) {
			sliceElementValue = sliceElementValue.Elem()
		}
		if !sliceElementValue.IsValid() {
			logging.Logger().Debugf("the original item indexed %d in given slice is invalid", i)
			return nil, fmt.Errorf("the original item indexed %d in given slice is invalid", i)
		}
		if !sliceElementValue.CanInterface() {
			logging.Logger().Debugf("the original item indexed %d in given slice cannot interface", i)
			return nil, fmt.Errorf("the original item indexed %d in given slice cannot interface", i)
		}

		if i == 0 {
			kindExpected = sliceElementValue.Kind()
		} else {
			if kindExpected != sliceElementValue.Kind() {
				logging.Logger().Debugf("the kind was expected as %s, but another kind %s has been found in the slice index %d", kindExpected, sliceElementValue.Kind(), i)
				return nil, fmt.Errorf("the kind was expected as %s, but another kind %s has been found in the slice index %d", kindExpected, sliceElementValue.Kind(), i)
			}
		}

		result[i] = sliceElementValue.Interface()
	}

	sort.Slice(result, func(i, j int) bool {
		if order == "desc" {
			i, j = j, i
		} else {
			if order != "asc" {
				logging.Logger().Debug("the order must be asc or desc")
				panic("the order must be asc or desc")
			}
		}
		leftInterface := result[i]
		rightInterface := result[j]
		switch kindExpected {
		case reflect.Int:
			left := leftInterface.(int)
			right := rightInterface.(int)
			return left < right
		case reflect.Int8:
			left := leftInterface.(int8)
			right := rightInterface.(int8)
			return left < right
		case reflect.Int16:
			left := leftInterface.(int16)
			right := rightInterface.(int16)
			return left < right
		case reflect.Int32:
			left := leftInterface.(int32)
			right := rightInterface.(int32)
			return left < right
		case reflect.Int64:
			left := leftInterface.(int64)
			right := rightInterface.(int64)
			return left < right
		case reflect.Uint:
			left := leftInterface.(uint)
			right := rightInterface.(uint)
			return left < right
		case reflect.Uint8:
			left := leftInterface.(uint8)
			right := rightInterface.(uint8)
			return left < right
		case reflect.Uint16:
			left := leftInterface.(uint16)
			right := rightInterface.(uint16)
			return left < right
		case reflect.Uint32:
			left := leftInterface.(uint32)
			right := rightInterface.(uint32)
			return left < right
		case reflect.Uint64:
			left := leftInterface.(uint64)
			right := rightInterface.(uint64)
			return left < right
		case reflect.Float32:
			left := leftInterface.(float32)
			right := rightInterface.(float32)
			return left < right
		case reflect.Float64:
			left := leftInterface.(float64)
			right := rightInterface.(float64)
			return left < right
		case reflect.String:
			left := leftInterface.(string)
			right := rightInterface.(string)
			return left < right
		default:
			logging.Logger().Debug("the value in the given slice must be comparable type like digit or string")
			panic("the value in the given slice must be comparable type like digit or string")
		}
	})
	return result, nil
}
