package collection

import (
	"fmt"
	"github.com/qb0C80aE/clay/logging"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
)

var utility = &Utility{}

// Utility handles collection operation
type Utility struct {
}

// GetUtility returns the instance of utility
func GetUtility() *Utility {
	return utility
}

// MapInterface is the alias of map[interface{}]interface{} what have operation methods
type MapInterface map[interface{}]interface{}

// Slice returns a slice object which consists of given items
func (receiver *Utility) Slice(items ...interface{}) interface{} {
	slice := []interface{}{}
	return append(slice, items...)
}

// SubSlice returns a sub slice object of given slice
func (receiver *Utility) SubSlice(sliceInterface interface{}, begin int, end int) (interface{}, error) {
	slice, err := mapstructutilpkg.GetUtility().SliceToInterfaceSlice(sliceInterface)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	if begin < 0 {
		if end < 0 {
			return slice[:], nil
		}
		return slice[:end], nil
	}
	if end < 0 {
		return slice[begin:], nil
	}
	return slice[begin:end], nil
}

// Append returns a slice appended given items
func (receiver *Utility) Append(sliceInterface interface{}, item ...interface{}) (interface{}, error) {
	slice, err := mapstructutilpkg.GetUtility().SliceToInterfaceSlice(sliceInterface)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	return append(slice, item...), nil
}

// Concatenate returns a slice concatenated two slices
func (receiver *Utility) Concatenate(sliceInterface1 interface{}, sliceInterface2 interface{}) (interface{}, error) {
	slice1, err := mapstructutilpkg.GetUtility().SliceToInterfaceSlice(sliceInterface1)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	slice2, err := mapstructutilpkg.GetUtility().SliceToInterfaceSlice(sliceInterface2)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	return append(slice1, slice2...), nil
}

// FieldSlice returns a slice of specified field value from given struct slice
func (receiver *Utility) FieldSlice(slice interface{}, fieldName string) (interface{}, error) {
	return mapstructutilpkg.GetUtility().StructSliceToFieldValueInterfaceSlice(slice, fieldName)
}

// Sort sorts the given slice and returns sorted one
func (receiver *Utility) Sort(slice interface{}, order string) (interface{}, error) {
	return mapstructutilpkg.GetUtility().SortSlice(slice, order)
}

// Hash makes a map from given struct slice with the specified field as the key
func (receiver *Utility) Hash(slice interface{}, keyField string) (MapInterface, error) {
	return mapstructutilpkg.GetUtility().StructSliceToInterfaceMap(slice, keyField)
}

// Sequence makes a slice of sequential integer values
func (receiver *Utility) Sequence(begin, end int) interface{} {
	count := end - begin + 1
	result := make([]int, count)
	for i, j := 0, begin; i < count; i, j = i+1, j+1 {
		result[i] = j
	}
	return result
}

// Map makes a map from given pairs
func (receiver *Utility) Map(pairs ...interface{}) (MapInterface, error) {
	if len(pairs)%2 == 1 {
		logging.Logger().Debug("numebr of arguments must be even")
		return nil, fmt.Errorf("numebr of arguments must be even")
	}
	m := make(MapInterface, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m, nil
}

// Exists returns if the key is in the map or not
func (receiver MapInterface) Exists(key interface{}) bool {
	_, exists := receiver[key]
	return exists
}

// Put puts given value into the map
func (receiver MapInterface) Put(key interface{}, value interface{}) MapInterface {
	receiver[key] = value
	return receiver
}

// Get gets the value related to given key from the map
func (receiver MapInterface) Get(key interface{}) interface{} {
	return receiver[key]
}

// Delete removes the specified value related to given key from the map
func (receiver MapInterface) Delete(key interface{}) MapInterface {
	delete(receiver, key)
	return receiver
}

// Merge merges the given map into this map
func (receiver MapInterface) Merge(source MapInterface) MapInterface {
	for key, value := range source {
		receiver[key] = value
	}
	return receiver
}

// Keys returns its keys as an slice
func (receiver MapInterface) Keys() (interface{}, error) {
	return mapstructutilpkg.GetUtility().MapToKeySlice(receiver)
}

// Values returns its values as an slice
func (receiver MapInterface) Values() interface{} {
	result := make([]interface{}, len(receiver))

	i := 0
	for _, value := range receiver {
		result[i] = value
		i++
	}

	return result
}

// Length returns its length
func (receiver MapInterface) Length() interface{} {
	return len(receiver)
}
