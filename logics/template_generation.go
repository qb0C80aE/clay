package logics

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/helper"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/conversion"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	tplpkg "text/template"
)

type templateGenerationLogic struct {
	*BaseLogic
}

func newTemplateGenerationLogic() *templateGenerationLogic {
	logic := &templateGenerationLogic{
		BaseLogic: NewBaseLogic(
			models.SharedTemplateModel(),
		),
	}
	return logic
}

// GenerateTemplate generates text data based on registered templates
func GenerateTemplate(db *gorm.DB, id string, templateArgumentParameterMap map[interface{}]interface{}) (interface{}, error) {
	templateArgumentMap := map[interface{}]*models.TemplateArgument{}
	templateParameterMap := map[interface{}]interface{}{}

	template := &models.Template{}
	template.ID, _ = strconv.Atoi(id)

	// GenerateTemplate reset db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db = db.New()

	if err := db.Preload("TemplateArguments").Select("*").First(template, template.ID).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	for _, templateArgument := range template.TemplateArguments {
		templateArgumentMap[templateArgument.Name] = templateArgument
		switch templateArgument.Type {
		case models.TemplateArgumentTypeInt:
			templateParameterMap[templateArgument.Name] = templateArgument.DefaultValueInt.Int64
		case models.TemplateArgumentTypeFloat:
			templateParameterMap[templateArgument.Name] = templateArgument.DefaultValueFloat.Float64
		case models.TemplateArgumentTypeBool:
			templateParameterMap[templateArgument.Name] = templateArgument.DefaultValueBool.Bool
		case models.TemplateArgumentTypeString:
			templateParameterMap[templateArgument.Name] = templateArgument.DefaultValueString.String
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
		case models.TemplateArgumentTypeInt:
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
		case models.TemplateArgumentTypeFloat:
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
		case models.TemplateArgumentTypeBool:
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
		case models.TemplateArgumentTypeString:
			templateParameterMap[key] = fmt.Sprintf("%v", value)
		}
	}

	templateParameterMap["ModelStore"] = db

	tpl := tplpkg.New("template")
	templateFuncMaps := extensions.RegisteredTemplateFuncMaps()
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
func (logic *templateGenerationLogic) GetSingle(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	templateArgumentParameterMap := make(map[interface{}]interface{}, len(urlValues))
	for key := range urlValues {
		templateArgumentParameterMap[key] = urlValues.Get(key)
	}
	return GenerateTemplate(db, parameters.ByName("id"), templateArgumentParameterMap)
}

var uniqueTemplateGenerationLogic = newTemplateGenerationLogic()

// UniqueTemplateGenerationLogic returns the unique template logic instance
func UniqueTemplateGenerationLogic() extensions.Logic {
	return uniqueTemplateGenerationLogic
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
			controller, err := extensions.AssociatedControllerWithPath(path)
			if err != nil {
				return nil, err
			}

			pathElements := strings.Split(strings.Trim(path, "/"), "/")
			resourceName := pathElements[0]
			singleURL := controller.ResourceSingleURL()
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
			model, err := extensions.CreateModel(resourceName)
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
			logic, err := extensions.RegisteredLogic(model)
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

			result, err := logic.GetSingle(db, parameters, requestForParameter.URL.Query(), queryFields)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return result, nil
		},
		"multi": func(dbObject interface{}, pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
			path := pathInterface.(string)
			controller, err := extensions.AssociatedControllerWithPath(path)
			if err != nil {
				return nil, err
			}

			pathElements := strings.Split(strings.Trim(path, "/"), "/")
			resourceName := pathElements[0]
			multiURL := controller.ResourceMultiURL()
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
			model, err := extensions.CreateModel(resourceName)
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
			logic, err := extensions.RegisteredLogic(model)
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
			result, err := logic.GetMulti(db, parameters, requestForParameter.URL.Query(), queryFields)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return result, nil
		},
		"total": func(dbObject interface{}, pathInterface interface{}) (interface{}, error) {
			path := pathInterface.(string)
			pathElements := strings.Split(strings.Trim(path, "/"), "/")
			resourceName := pathElements[0]

			model, err := extensions.CreateModel(resourceName)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			db := dbObject.(*gorm.DB)
			var total = 0
			db.Model(model).Count(&total)
			return total, nil
		},
		"include": func(dbObject interface{}, templateName string, templateArgumentParameterMap map[interface{}]interface{}) (interface{}, error) {
			db := dbObject.(*gorm.DB)
			template := models.NewTemplateModel()
			if err := db.Select("*").First(template, "name = ?", templateName).Error; err != nil {
				logging.Logger().Debug(err.Error())
				return nil, fmt.Errorf("template %s does not exist", templateName)
			}
			result, err := GenerateTemplate(db, strconv.Itoa(template.ID), templateArgumentParameterMap)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			return result, nil
		},
	}
	extensions.RegisterTemplateFuncMap(funcMap)
}
