package conversion

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qb0C80aE/clay/logging"
	"gopkg.in/yaml.v2"
	"reflect"
	"strconv"
)

var utility = &Utility{}

// Utility handles conversion operation
type Utility struct {
}

// GetUtility returns the instance of utility
func GetUtility() *Utility {
	return utility
}

// Int converts interface{} into interface{} what has a int value
func (receiver *Utility) Int(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return int(value.(int)), nil
	case int8:
		return int(value.(int8)), nil
	case int16:
		return int(value.(int16)), nil
	case int32:
		return int(value.(int32)), nil
	case int64:
		return int(value.(int64)), nil
	case uint:
		return int(value.(uint)), nil
	case uint8:
		return int(value.(uint8)), nil
	case uint16:
		return int(value.(uint16)), nil
	case uint32:
		return int(value.(uint32)), nil
	case uint64:
		return int(value.(uint64)), nil
	case float32:
		return int(value.(float32)), nil
	case float64:
		return int(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseInt(string(value.([]byte)), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int(val), nil
	case string:
		val, err := strconv.ParseInt(value.(string), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Int8 converts interface{} into interface{} what has a int8 value
func (receiver *Utility) Int8(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return int8(value.(int)), nil
	case int8:
		return int8(value.(int8)), nil
	case int16:
		return int8(value.(int16)), nil
	case int32:
		return int8(value.(int32)), nil
	case int64:
		return int8(value.(int64)), nil
	case uint:
		return int8(value.(uint)), nil
	case uint8:
		return int8(value.(uint8)), nil
	case uint16:
		return int8(value.(uint16)), nil
	case uint32:
		return int8(value.(uint32)), nil
	case uint64:
		return int8(value.(uint64)), nil
	case float32:
		return int8(value.(float32)), nil
	case float64:
		return int8(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseInt(string(value.([]byte)), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int8(val), nil
	case string:
		val, err := strconv.ParseInt(value.(string), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int8(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Int16 converts interface{} into interface{} what has a int16 value
func (receiver *Utility) Int16(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return int16(value.(int)), nil
	case int8:
		return int16(value.(int8)), nil
	case int16:
		return int16(value.(int16)), nil
	case int32:
		return int16(value.(int32)), nil
	case int64:
		return int16(value.(int64)), nil
	case uint:
		return int16(value.(uint)), nil
	case uint8:
		return int16(value.(uint8)), nil
	case uint16:
		return int16(value.(uint16)), nil
	case uint32:
		return int16(value.(uint32)), nil
	case uint64:
		return int16(value.(uint64)), nil
	case float32:
		return int16(value.(float32)), nil
	case float64:
		return int16(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseInt(string(value.([]byte)), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int16(val), nil
	case string:
		val, err := strconv.ParseInt(value.(string), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int16(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Int32 converts interface{} into interface{} what has a int32 value
func (receiver *Utility) Int32(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return int32(value.(int)), nil
	case int8:
		return int32(value.(int8)), nil
	case int16:
		return int32(value.(int16)), nil
	case int32:
		return int32(value.(int32)), nil
	case int64:
		return int32(value.(int64)), nil
	case uint:
		return int32(value.(uint)), nil
	case uint8:
		return int32(value.(uint8)), nil
	case uint16:
		return int32(value.(uint16)), nil
	case uint32:
		return int32(value.(uint32)), nil
	case uint64:
		return int32(value.(uint64)), nil
	case float32:
		return int32(value.(float32)), nil
	case float64:
		return int32(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseInt(string(value.([]byte)), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int32(val), nil
	case string:
		val, err := strconv.ParseInt(value.(string), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int32(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Int64 converts interface{} into interface{} what has a int64 value
func (receiver *Utility) Int64(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return int64(value.(int)), nil
	case int8:
		return int64(value.(int8)), nil
	case int16:
		return int64(value.(int16)), nil
	case int32:
		return int64(value.(int32)), nil
	case int64:
		return int64(value.(int64)), nil
	case uint:
		return int64(value.(uint)), nil
	case uint8:
		return int64(value.(uint8)), nil
	case uint16:
		return int64(value.(uint16)), nil
	case uint32:
		return int64(value.(uint32)), nil
	case uint64:
		return int64(value.(uint64)), nil
	case float32:
		return int64(value.(float32)), nil
	case float64:
		return int64(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseInt(string(value.([]byte)), 10, 64)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int64(val), nil
	case string:
		val, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return int64(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Uint converts interface{} into interface{} what has a uint value
func (receiver *Utility) Uint(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return uint(value.(int)), nil
	case int8:
		return uint(value.(int8)), nil
	case int16:
		return uint(value.(int16)), nil
	case int32:
		return uint(value.(int32)), nil
	case int64:
		return uint(value.(int64)), nil
	case uint:
		return uint(value.(uint)), nil
	case uint8:
		return uint(value.(uint8)), nil
	case uint16:
		return uint(value.(uint16)), nil
	case uint32:
		return uint(value.(uint32)), nil
	case uint64:
		return uint(value.(uint64)), nil
	case float32:
		return uint(value.(float32)), nil
	case float64:
		return uint(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseUint(string(value.([]byte)), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint(val), nil
	case string:
		val, err := strconv.ParseUint(value.(string), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Uint8 converts interface{} into interface{} what has a uint8 value
func (receiver *Utility) Uint8(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return uint8(value.(int)), nil
	case int8:
		return uint8(value.(int8)), nil
	case int16:
		return uint8(value.(int16)), nil
	case int32:
		return uint8(value.(int32)), nil
	case int64:
		return uint8(value.(int64)), nil
	case uint:
		return uint8(value.(uint)), nil
	case uint8:
		return uint8(value.(uint8)), nil
	case uint16:
		return uint8(value.(uint16)), nil
	case uint32:
		return uint8(value.(uint32)), nil
	case uint64:
		return uint8(value.(uint64)), nil
	case float32:
		return uint8(value.(float32)), nil
	case float64:
		return uint8(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseUint(string(value.([]byte)), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint8(val), nil
	case string:
		val, err := strconv.ParseUint(value.(string), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint8(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Uint16 converts interface{} into interface{} what has a uint16 value
func (receiver *Utility) Uint16(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return uint16(value.(int)), nil
	case int8:
		return uint16(value.(int8)), nil
	case int16:
		return uint16(value.(int16)), nil
	case int32:
		return uint16(value.(int32)), nil
	case int64:
		return uint16(value.(int64)), nil
	case uint:
		return uint16(value.(uint)), nil
	case uint8:
		return uint16(value.(uint8)), nil
	case uint16:
		return uint16(value.(uint16)), nil
	case uint32:
		return uint16(value.(uint32)), nil
	case uint64:
		return uint16(value.(uint64)), nil
	case float32:
		return uint16(value.(float32)), nil
	case float64:
		return uint16(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseUint(string(value.([]byte)), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint16(val), nil
	case string:
		val, err := strconv.ParseUint(value.(string), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint16(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Uint32 converts interface{} into interface{} what has a uint32 value
func (receiver *Utility) Uint32(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return uint32(value.(int)), nil
	case int8:
		return uint32(value.(int8)), nil
	case int16:
		return uint32(value.(int16)), nil
	case int32:
		return uint32(value.(int32)), nil
	case int64:
		return uint32(value.(int64)), nil
	case uint:
		return uint32(value.(uint)), nil
	case uint8:
		return uint32(value.(uint8)), nil
	case uint16:
		return uint32(value.(uint16)), nil
	case uint32:
		return uint32(value.(uint32)), nil
	case uint64:
		return uint32(value.(uint64)), nil
	case float32:
		return uint32(value.(float32)), nil
	case float64:
		return uint32(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseUint(string(value.([]byte)), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint32(val), nil
	case string:
		val, err := strconv.ParseUint(value.(string), 10, 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint32(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Uint64 converts interface{} into interface{} what has a uint64 value
func (receiver *Utility) Uint64(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return uint64(value.(int)), nil
	case int8:
		return uint64(value.(int8)), nil
	case int16:
		return uint64(value.(int16)), nil
	case int32:
		return uint64(value.(int32)), nil
	case int64:
		return uint64(value.(int64)), nil
	case uint:
		return uint64(value.(uint)), nil
	case uint8:
		return uint64(value.(uint8)), nil
	case uint16:
		return uint64(value.(uint16)), nil
	case uint32:
		return uint64(value.(uint32)), nil
	case uint64:
		return uint64(value.(uint64)), nil
	case float32:
		return uint64(value.(float32)), nil
	case float64:
		return uint64(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseUint(string(value.([]byte)), 10, 64)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint64(val), nil
	case string:
		val, err := strconv.ParseUint(value.(string), 10, 64)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return uint64(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Float32 converts interface{} into interface{} what has a float32 value
func (receiver *Utility) Float32(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return float32(value.(int)), nil
	case int8:
		return float32(value.(int8)), nil
	case int16:
		return float32(value.(int16)), nil
	case int32:
		return float32(value.(int32)), nil
	case int64:
		return float32(value.(int64)), nil
	case uint:
		return float32(value.(uint)), nil
	case uint8:
		return float32(value.(uint8)), nil
	case uint16:
		return float32(value.(uint16)), nil
	case uint32:
		return float32(value.(uint32)), nil
	case uint64:
		return float32(value.(uint64)), nil
	case float32:
		return float32(value.(float32)), nil
	case float64:
		return float32(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseFloat(string(value.([]byte)), 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return float32(val), nil
	case string:
		val, err := strconv.ParseFloat(value.(string), 32)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return float32(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// Float64 converts interface{} into interface{} what has a float64 value
func (receiver *Utility) Float64(value interface{}) (interface{}, error) {
	switch value.(type) {
	case int:
		return float64(value.(int)), nil
	case int8:
		return float64(value.(int8)), nil
	case int16:
		return float64(value.(int16)), nil
	case int32:
		return float64(value.(int32)), nil
	case int64:
		return float64(value.(int64)), nil
	case uint:
		return float64(value.(uint)), nil
	case uint8:
		return float64(value.(uint8)), nil
	case uint16:
		return float64(value.(uint16)), nil
	case uint32:
		return float64(value.(uint32)), nil
	case uint64:
		return float64(value.(uint64)), nil
	case float32:
		return float64(value.(float32)), nil
	case float64:
		return float64(value.(float64)), nil
	case []byte:
		val, err := strconv.ParseFloat(string(value.([]byte)), 64)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return float64(val), nil
	case string:
		val, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return float64(val), nil
	default:
		logging.Logger().Debug("value is not int, float, or number-strting")
		return nil, errors.New("value is not int, float, or number-strting")
	}
}

// String converts interface{} into interface{} what has a string value
func (receiver *Utility) String(value interface{}) interface{} {
	switch value.(type) {
	case int:
		return strconv.FormatInt(int64(value.(int)), 10)
	case int8:
		return strconv.FormatInt(int64(value.(int8)), 10)
	case int16:
		return strconv.FormatInt(int64(value.(int16)), 10)
	case int32:
		return strconv.FormatInt(int64(value.(int32)), 10)
	case int64:
		return strconv.FormatInt(int64(value.(int64)), 10)
	case uint:
		return strconv.FormatUint(uint64(value.(uint)), 10)
	case uint8:
		return strconv.FormatUint(uint64(value.(uint8)), 10)
	case uint16:
		return strconv.FormatUint(uint64(value.(uint16)), 10)
	case uint32:
		return strconv.FormatUint(uint64(value.(uint32)), 10)
	case uint64:
		return strconv.FormatUint(uint64(value.(uint64)), 10)
	case float32:
		return strconv.FormatFloat(float64(value.(float32)), 'E', -1, 32)
	case float64:
		return strconv.FormatFloat(float64(value.(float64)), 'E', -1, 64)
	case bool:
		return strconv.FormatBool(value.(bool))
	case []byte:
		return string(value.([]byte))
	default:
		return fmt.Sprintf("%v", value)
	}
}

// Boolean converts interface{} into interface{} what has a boolean value
func (receiver *Utility) Boolean(value interface{}) (interface{}, error) {
	switch value.(type) {
	case []byte:
		val, err := strconv.ParseBool(string(value.([]byte)))
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return val, nil
	case string:
		val, err := strconv.ParseBool(value.(string))
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return val, nil
	default:
		logging.Logger().Debug("value is not boolean-string")
		return nil, errors.New("value is not boolean-string")
	}
}

// Bytes converts string interface{} into interface{} what has []byte value
func (receiver *Utility) Bytes(value interface{}) (interface{}, error) {
	switch value.(type) {
	case string:
		return []byte(value.(string)), nil
	default:
		logging.Logger().Debug("value is not string")
		return nil, errors.New("value is not string")
	}
}

func (receiver *Utility) convertToStringKeyMap(dataInterface interface{}) interface{} {
	reflectValue := reflect.ValueOf(dataInterface)
	for (reflectValue.Kind() == reflect.Ptr) || (reflectValue.Kind() == reflect.Interface) {
		reflectValue = reflectValue.Elem()
	}

	switch reflectValue.Kind() {
	case reflect.Map:
		result := map[string]interface{}{}
		mapKeys := reflectValue.MapKeys()

		for _, key := range mapKeys {
			itemValue := reflectValue.MapIndex(key)
			itemInterface := itemValue.Interface()
			item := receiver.convertToStringKeyMap(itemInterface)
			result[fmt.Sprintf("%v", key)] = item
		}

		return result
	default:
		return dataInterface
	}
}

// JSONMarshal generates JSON text from given object
func (receiver *Utility) JSONMarshal(dataInterface interface{}, indentInterface interface{}) (interface{}, error) {
	indent := indentInterface.(string)

	switch len(indent) {
	case 0:
		result, err := json.Marshal(receiver.convertToStringKeyMap(dataInterface))
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return string(result), nil
	default:
		result, err := json.MarshalIndent(receiver.convertToStringKeyMap(dataInterface), "", indent)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
		return string(result), nil
	}
}

// YAMLMarshal generates YAML text from given object
func (receiver *Utility) YAMLMarshal(dataInterface interface{}) (interface{}, error) {
	result, err := yaml.Marshal(receiver.convertToStringKeyMap(dataInterface))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	return string(result), nil
}
