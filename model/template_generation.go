package model

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/helper"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/util/conversion"
	"github.com/qb0C80aE/clay/util/mapstruct"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	tplpkg "text/template"
)

// TemplateGeneration is the model class what represents template generation
type TemplateGeneration struct {
	*Base `json:"base,omitempty"`
}

// NewTemplateGeneration creates a template generation model instance
func NewTemplateGeneration() *TemplateGeneration {
	return CreateModel(&TemplateGeneration{}).(*TemplateGeneration)
}

// GenerateTemplate generates text data based on registered templates
func (receiver *TemplateGeneration) GenerateTemplate(db *gorm.DB, id string, templateArgumentParameterMap map[interface{}]interface{}) (interface{}, error) {
	templateArgumentMap := map[interface{}]*TemplateArgument{}
	templateParameterMap := map[interface{}]interface{}{}

	template := NewTemplate().NewModelContainer().(*Template)
	template.ID, _ = strconv.Atoi(id)

	// GenerateTemplate reset db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db = db.New()

	if err := db.Preload("TemplateArguments").Select("*").First(template, template.ID).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	for _, templateArgument := range template.TemplateArguments {
		var err error = nil
		templateArgumentMap[templateArgument.Name] = templateArgument
		switch templateArgument.Type {
		case TemplateArgumentTypeInt:
			templateParameterMap[templateArgument.Name], err = conversion.ToInt64Interface(templateArgument.DefaultValue)
		case TemplateArgumentTypeFloat:
			templateParameterMap[templateArgument.Name], err = conversion.ToFloat64Interface(templateArgument.DefaultValue)
		case TemplateArgumentTypeBool:
			templateParameterMap[templateArgument.Name], err = conversion.ToBooleanInterface(templateArgument.DefaultValue)
		case TemplateArgumentTypeString:
			templateParameterMap[templateArgument.Name] = templateArgument.DefaultValue
		default:
			err = fmt.Errorf("invalid type: %v", templateArgument.Type)
		}

		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}
	}

	for key, value := range templateArgumentParameterMap {
		templateArgument, exists := templateArgumentMap[key]
		if !exists {
			logging.Logger().Debugf("the argument related to a parameter %s does not exist", key)
			return nil, fmt.Errorf("the argument related to a parameter %s does not exist", key)
		}

		valueType := reflect.TypeOf(value)
		switch templateArgument.Type {
		case TemplateArgumentTypeInt:
			switch valueType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				templateParameterMap[key] = int64(value.(int))
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				templateParameterMap[key] = int64(value.(uint))
			case reflect.String:
				var err error
				templateParameterMap[key], err = strconv.ParseInt(value.(string), 10, 64)
				if err != nil {
					logging.Logger().Debug(err.Error())
					return nil, fmt.Errorf("parameter type mistmatch, %s should be int, or integer-formatted string, but value is %v", key, value)
				}
			default:
				return nil, fmt.Errorf("parameter type mistmatch, %s should be int, or integer-formatted string, but value is %v", key, value)
			}
		case TemplateArgumentTypeFloat:
			switch valueType.Kind() {
			case reflect.Float32, reflect.Float64:
				templateParameterMap[key] = float64(value.(float64))
			case reflect.String:
				var err error
				templateParameterMap[key], err = strconv.ParseFloat(value.(string), 64)
				if err != nil {
					logging.Logger().Debug(err.Error())
					return nil, fmt.Errorf("parameter type mistmatch, %s should be float, or float-formatted string, but value is %v", key, value)
				}
			default:
				return nil, fmt.Errorf("parameter type mistmatch, %s should be float, or float-formatted string, but value is %v", key, value)
			}
		case TemplateArgumentTypeBool:
			switch valueType.Kind() {
			case reflect.Bool:
				templateParameterMap[key] = value.(bool)
			case reflect.String:
				var err error
				templateParameterMap[key], err = strconv.ParseBool(value.(string))
				if err != nil {
					logging.Logger().Debug(err.Error())
					return nil, fmt.Errorf("parameter type mistmatch, %s should be bool, or bool-formatted string, but value is %v", key, value)
				}
			default:
				return nil, fmt.Errorf("parameter type mistmatch, %s should be bool, or bool-formatted string, but value is %v", key, value)
			}
		case TemplateArgumentTypeString:
			templateParameterMap[key] = fmt.Sprintf("%v", value)
		}
	}

	templateParameterMap["ModelStore"] = db

	tpl := tplpkg.New("template")
	templateFuncMaps := extension.GetRegisteredTemplateFuncMapList()
	for _, templateFuncMap := range templateFuncMaps {
		tpl = tpl.Funcs(templateFuncMap)
	}
	tpl, err := tpl.Parse(template.TemplateContent)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	var doc bytes.Buffer
	if err := tpl.Execute(&doc, templateParameterMap); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	result := doc.String()

	return result, nil
}

// GetSingle generates text data based on registered templates
func (receiver *TemplateGeneration) GetSingle(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	templateArgumentParameterMap := make(map[interface{}]interface{}, len(urlValues))
	for key := range urlValues {
		templateArgumentParameterMap[key] = urlValues.Get(key)
	}
	return receiver.GenerateTemplate(db, parameters.ByName("id"), templateArgumentParameterMap)
}

func init() {
	funcMap := tplpkg.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
		"mod": func(a, b int) int { return a % b },
		"int": func(value interface{}) (interface{}, error) {
			return conversion.ToIntInterface(value)
		},
		"int8": func(value interface{}) (interface{}, error) {
			return conversion.ToInt8Interface(value)
		},
		"int16": func(value interface{}) (interface{}, error) {
			return conversion.ToInt16Interface(value)
		},
		"int32": func(value interface{}) (interface{}, error) {
			return conversion.ToInt32Interface(value)
		},
		"int64": func(value interface{}) (interface{}, error) {
			return conversion.ToInt64Interface(value)
		},
		"uint": func(value interface{}) (interface{}, error) {
			return conversion.ToUintInterface(value)
		},
		"uint8": func(value interface{}) (interface{}, error) {
			return conversion.ToUint8Interface(value)
		},
		"uint16": func(value interface{}) (interface{}, error) {
			return conversion.ToUint16Interface(value)
		},
		"uint32": func(value interface{}) (interface{}, error) {
			return conversion.ToUint32Interface(value)
		},
		"uint64": func(value interface{}) (interface{}, error) {
			return conversion.ToUint64Interface(value)
		},
		"float32": func(value interface{}) (interface{}, error) {
			return conversion.ToFloat32Interface(value)
		},
		"float64": func(value interface{}) (interface{}, error) {
			return conversion.ToFloat64Interface(value)
		},
		"string": func(value interface{}) interface{} {
			return conversion.ToStringInterface(value)
		},
		"boolean": func(value interface{}) (interface{}, error) {
			return conversion.ToBooleanInterface(value)
		},
		"split": func(value interface{}, separator string) interface{} {
			data := fmt.Sprintf("%v", value)
			return strings.Split(data, separator)
		},
		"join": func(slice interface{}, separator string) (interface{}, error) {
			interfaceSlice, err := mapstruct.SliceToInterfaceSlice(slice)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			stringSlice := make([]string, len(interfaceSlice))

			for index, item := range interfaceSlice {
				stringSlice[index] = fmt.Sprintf("%v", item)
			}

			return strings.Join(stringSlice, separator), nil
		},
		"slice": func(items ...interface{}) interface{} {
			slice := []interface{}{}
			return append(slice, items...)
		},
		"subslice": func(sliceInterface interface{}, begin int, end int) (interface{}, error) {
			slice, err := mapstruct.SliceToInterfaceSlice(sliceInterface)
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
		},
		"append": func(sliceInterface interface{}, item ...interface{}) (interface{}, error) {
			slice, err := mapstruct.SliceToInterfaceSlice(sliceInterface)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return append(slice, item...), nil
		},
		"concatenate": func(sliceInterface1 interface{}, sliceInterface2 interface{}) (interface{}, error) {
			slice1, err := mapstruct.SliceToInterfaceSlice(sliceInterface1)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			slice2, err := mapstruct.SliceToInterfaceSlice(sliceInterface2)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return append(slice1, slice2...), nil
		},
		"fieldslice": func(slice interface{}, fieldName string) ([]interface{}, error) {
			return mapstruct.StructSliceToFieldValueInterfaceSlice(slice, fieldName)
		},
		"sort": func(slice interface{}, order string) ([]interface{}, error) {
			return mapstruct.SortSlice(slice, order)
		},
		"map": func(pairs ...interface{}) (map[interface{}]interface{}, error) {
			if len(pairs)%2 == 1 {
				logging.Logger().Debug("numebr of arguments must be even")
				return nil, fmt.Errorf("numebr of arguments must be even")
			}
			m := make(map[interface{}]interface{}, len(pairs)/2)
			for i := 0; i < len(pairs); i += 2 {
				m[pairs[i]] = pairs[i+1]
			}
			return m, nil
		},
		"exists": func(target map[interface{}]interface{}, key interface{}) bool {
			_, exists := target[key]
			return exists
		},
		"put": func(target map[interface{}]interface{}, key interface{}, value interface{}) map[interface{}]interface{} {
			target[key] = value
			return target
		},
		"get": func(target map[interface{}]interface{}, key interface{}) interface{} {
			return target[key]
		},
		"delete": func(target map[interface{}]interface{}, key interface{}) map[interface{}]interface{} {
			delete(target, key)
			return target
		},
		"merge": func(source, destination map[interface{}]interface{}) map[interface{}]interface{} {
			for key, value := range source {
				destination[key] = value
			}
			return destination
		},
		"keys": func(target map[interface{}]interface{}) ([]interface{}, error) {
			return mapstruct.MapToKeySlice(target)
		},
		"hash": func(slice interface{}, keyField string) (map[interface{}]interface{}, error) {
			return mapstruct.StructSliceToInterfaceMap(slice, keyField)
		},
		"slicemap": func(slice interface{}, keyField string) (map[interface{}]interface{}, error) {
			sliceMap, err := mapstruct.StructSliceToInterfaceSliceMap(slice, keyField)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			result := make(map[interface{}]interface{}, len(sliceMap))
			for key, value := range sliceMap {
				result[key] = value
			}
			return result, nil
		},
		"sequence": func(begin, end int) interface{} {
			count := end - begin + 1
			result := make([]int, count)
			for i, j := 0, begin; i < count; i, j = i+1, j+1 {
				result[i] = j
			}
			return result
		},
		"single": func(dbObject interface{}, pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
			path := pathInterface.(string)
			controller, err := extension.GetAssociatedControllerWithPath(path)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}

			pathElements := strings.Split(strings.Trim(path, "/"), "/")
			resourceName := pathElements[0]
			singleURL := controller.GetResourceSingleURL()
			routeElements := strings.Split(strings.Trim(singleURL, "/"), "/")

			parameters := gin.Params{}
			for index, routeElement := range routeElements {
				if routeElement[:1] == ":" {
					parameter := gin.Param{
						Key:   routeElement[1:],
						Value: pathElements[index],
					}
					parameters = append(parameters, parameter)
				}
			}

			query := queryInterface.(string)
			URL := "/"
			if query != "" {
				URL = "/?" + query
			}
			requestForParameter, err := http.NewRequest(
				http.MethodGet,
				URL,
				nil,
			)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			modelContainer, err := extension.CreateModelContainerByResourceName(resourceName)
			model := CreateModel(modelContainer)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			urlQuery := requestForParameter.URL.Query()
			parameter, err := dbpkg.NewParameter(urlQuery)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			db := dbObject.(*gorm.DB)
			db, err = parameter.Paginate(db)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			db = parameter.SetPreloads(db)
			db = parameter.FilterFields(db)
			fields := helper.ParseFields(parameter.DefaultQuery(urlQuery, "fields", "*"))
			queryFields := helper.QueryFields(model, fields)

			result, err := model.ExecuteActualGetSingle(db, parameters, requestForParameter.URL.Query(), queryFields)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return result, nil
		},
		"multi": func(dbObject interface{}, pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
			path := pathInterface.(string)
			controller, err := extension.GetAssociatedControllerWithPath(path)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}

			pathElements := strings.Split(strings.Trim(path, "/"), "/")
			resourceName := pathElements[0]
			multiURL := controller.GetResourceMultiURL()
			routeElements := strings.Split(strings.Trim(multiURL, "/"), "/")

			parameters := gin.Params{}
			for index, routeElement := range routeElements {
				if routeElement[:1] == ":" {
					parameter := gin.Param{
						Key:   routeElement[1:],
						Value: pathElements[index],
					}
					parameters = append(parameters, parameter)
				}
			}

			query := queryInterface.(string)
			URL := "/"
			if query != "" {
				URL = "/?" + query
			}
			requestForParameter, err := http.NewRequest(
				http.MethodGet,
				URL,
				nil,
			)
			modelContainer, err := extension.CreateModelContainerByResourceName(resourceName)
			model := CreateModel(modelContainer)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			urlQuery := requestForParameter.URL.Query()
			parameter, err := dbpkg.NewParameter(urlQuery)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			db := dbObject.(*gorm.DB)
			db, err = parameter.Paginate(db)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			db = parameter.SetPreloads(db)
			db = parameter.SortRecords(db)
			db = parameter.FilterFields(db)
			fields := helper.ParseFields(parameter.DefaultQuery(urlQuery, "fields", "*"))
			queryFields := helper.QueryFields(model, fields)
			result, err := model.ExecuteActualGetMulti(db, parameters, requestForParameter.URL.Query(), queryFields)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return result, nil
		},
		"first": func(dbObject interface{}, pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
			path := pathInterface.(string)
			controller, err := extension.GetAssociatedControllerWithPath(path)
			if err != nil {
				return nil, err
			}

			pathElements := strings.Split(strings.Trim(path, "/"), "/")
			resourceName := pathElements[0]
			multiURL := controller.GetResourceMultiURL()
			routeElements := strings.Split(strings.Trim(multiURL, "/"), "/")

			parameters := gin.Params{}
			for index, routeElement := range routeElements {
				if routeElement[:1] == ":" {
					parameter := gin.Param{
						Key:   routeElement[1:],
						Value: pathElements[index],
					}
					parameters = append(parameters, parameter)
				}
			}

			query := queryInterface.(string)
			URL := "/"
			if query != "" {
				URL = "/?" + query
			}
			requestForParameter, err := http.NewRequest(
				http.MethodGet,
				URL,
				nil,
			)
			modelContainer, err := extension.CreateModelContainerByResourceName(resourceName)
			model := CreateModel(modelContainer)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			urlQuery := requestForParameter.URL.Query()
			parameter, err := dbpkg.NewParameter(urlQuery)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			db := dbObject.(*gorm.DB)
			db, err = parameter.Paginate(db)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			db = parameter.SetPreloads(db)
			db = parameter.SortRecords(db)
			db = parameter.FilterFields(db)
			fields := helper.ParseFields(parameter.DefaultQuery(urlQuery, "fields", "*"))
			queryFields := helper.QueryFields(model, fields)
			result, err := model.ExecuteActualGetMulti(db, parameters, requestForParameter.URL.Query(), queryFields)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}

			resultValue := reflect.ValueOf(result)
			if resultValue.Len() == 0 {
				logging.Logger().Debug("no record selected")
				return nil, errors.New("no record selected")
			}

			return resultValue.Index(0).Interface(), nil
		},
		"total": func(dbObject interface{}, pathInterface interface{}) (interface{}, error) {
			path := pathInterface.(string)
			pathElements := strings.Split(strings.Trim(path, "/"), "/")
			resourceName := pathElements[0]

			modelContainer, err := extension.CreateModelContainerByResourceName(resourceName)
			model := CreateModel(modelContainer)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			db := dbObject.(*gorm.DB)
			total, err := model.ExecuteActualGetTotal(db)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return total, nil
		},
		"include": func(dbObject interface{}, templateName string, templateArgumentParameterMap map[interface{}]interface{}) (interface{}, error) {
			db := dbObject.(*gorm.DB)
			template := NewTemplate()
			if err := db.Select("*").First(template, "name = ?", templateName).Error; err != nil {
				logging.Logger().Debug(err.Error())
				return nil, fmt.Errorf("template %s does not exist", templateName)
			}
			result, err := NewTemplateGeneration().GenerateTemplate(db, strconv.Itoa(template.ID), templateArgumentParameterMap)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return result, nil
		},
	}
	extension.RegisterTemplateFuncMap(funcMap)
}

func init() {
	extension.RegisterModel(NewTemplateGeneration(), false)
}
