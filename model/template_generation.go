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
	"github.com/qb0C80aE/clay/util/network"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	tplpkg "text/template"
)

var parameterRegexp = regexp.MustCompile("p\\[(.+)\\]")

type templateParameter struct {
	ModelStore *modelStore
	Core       *coreUtil
	Network    *networkUtil
	Parameter  map[interface{}]interface{}
	Query      url.Values
}

type modelStore struct {
	db *gorm.DB
}

type coreUtil struct {
}

type networkUtil struct {
}

// TemplateGeneration is the model class what represents template generation
type TemplateGeneration struct {
	Base
}

// NewTemplateGeneration creates a template generation model instance
func NewTemplateGeneration() *TemplateGeneration {
	return &TemplateGeneration{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *TemplateGeneration) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GenerateTemplate generates text data based on registered templates
// parameters include either id or name
// actual parameters for template arguments must be included in urlValues as shaped like q[...]=...
func (receiver *TemplateGeneration) GenerateTemplate(db *gorm.DB, parameters gin.Params, urlValues url.Values) (interface{}, error) {
	templateArgumentMap := map[interface{}]*TemplateArgument{}
	templateParameterMap := map[interface{}]interface{}{}

	templateArgumentParameterMap := map[interface{}]interface{}{}
	for key := range urlValues {
		subMatch := parameterRegexp.FindStringSubmatch(key)
		if len(subMatch) == 2 {
			templateArgumentParameterMap[subMatch[1]] = urlValues.Get(key)
		}
	}

	templateModel := NewTemplate()
	templateModelAsContainer := NewTemplate()

	// GenerateTemplate resets db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db = db.New()

	newURLValues := url.Values{}
	if len(urlValues.Get("key_parameter")) > 0 {
		newURLValues.Set("key_parameter", urlValues.Get("key_parameter"))
	}
	newURLValues.Set("preloads", "template_arguments")

	dbParameter, err := dbpkg.NewParameter(newURLValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	db = dbParameter.SetPreloads(db)

	container, err := templateModel.GetSingle(templateModel, db, parameters, newURLValues, "*")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	if err := mapstruct.RemapToStruct(container, templateModelAsContainer); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	for _, templateArgument := range templateModelAsContainer.TemplateArguments {
		var err error
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

	tpl := tplpkg.New("template")

	templateFuncMaps := extension.GetRegisteredTemplateFuncMapList()
	for _, templateFuncMap := range templateFuncMaps {
		tpl = tpl.Funcs(templateFuncMap)
	}

	tpl, err = tpl.Parse(templateModelAsContainer.TemplateContent)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	templateParameter := &templateParameter{
		ModelStore: &modelStore{
			db: db,
		},
		Core:      &coreUtil{},
		Network:   &networkUtil{},
		Parameter: templateParameterMap,
		Query:     urlValues,
	}

	var doc bytes.Buffer
	if err := tpl.Execute(&doc, templateParameter); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	result := doc.String()

	return result, nil
}

// GetSingle generates text data based on registered templates
// parameters must be given as p[...]=...
func (receiver *TemplateGeneration) GetSingle(_ extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, _ string) (interface{}, error) {
	return receiver.GenerateTemplate(db, parameters, urlValues)
}

func (receiver *coreUtil) Add(a, b int) int {
	return a + b
}

func (receiver *coreUtil) Sub(a, b int) int {
	return a - b
}

func (receiver *coreUtil) Mul(a, b int) int {
	return a * b
}

func (receiver *coreUtil) Div(a, b int) int {
	return a / b
}

func (receiver *coreUtil) Mod(a, b int) int {
	return a % b
}

func (receiver *coreUtil) Int(value interface{}) (interface{}, error) {
	return conversion.ToIntInterface(value)
}

func (receiver *coreUtil) Int8(value interface{}) (interface{}, error) {
	return conversion.ToInt8Interface(value)
}

func (receiver *coreUtil) Int16(value interface{}) (interface{}, error) {
	return conversion.ToInt16Interface(value)
}

func (receiver *coreUtil) Int32(value interface{}) (interface{}, error) {
	return conversion.ToInt32Interface(value)
}

func (receiver *coreUtil) Int64(value interface{}) (interface{}, error) {
	return conversion.ToInt64Interface(value)
}

func (receiver *coreUtil) Uint(value interface{}) (interface{}, error) {
	return conversion.ToUintInterface(value)
}

func (receiver *coreUtil) Uint8(value interface{}) (interface{}, error) {
	return conversion.ToUint8Interface(value)
}

func (receiver *coreUtil) Uint16(value interface{}) (interface{}, error) {
	return conversion.ToUint16Interface(value)
}

func (receiver *coreUtil) Uint32(value interface{}) (interface{}, error) {
	return conversion.ToUint32Interface(value)
}

func (receiver *coreUtil) Uint64(value interface{}) (interface{}, error) {
	return conversion.ToUint64Interface(value)
}

func (receiver *coreUtil) Float32(value interface{}) (interface{}, error) {
	return conversion.ToFloat32Interface(value)
}

func (receiver *coreUtil) Float64(value interface{}) (interface{}, error) {
	return conversion.ToFloat64Interface(value)
}

func (receiver *coreUtil) String(value interface{}) interface{} {
	return conversion.ToStringInterface(value)
}

func (receiver *coreUtil) Boolean(value interface{}) (interface{}, error) {
	return conversion.ToBooleanInterface(value)
}

func (receiver *coreUtil) Split(value interface{}, separator string) interface{} {
	data := fmt.Sprintf("%v", value)
	return strings.Split(data, separator)
}

func (receiver *coreUtil) Join(slice interface{}, separator string) (interface{}, error) {
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
}

func (receiver *coreUtil) Slice(items ...interface{}) interface{} {
	slice := []interface{}{}
	return append(slice, items...)
}

func (receiver *coreUtil) SubSlice(sliceInterface interface{}, begin int, end int) (interface{}, error) {
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
}

func (receiver *coreUtil) Append(sliceInterface interface{}, item ...interface{}) (interface{}, error) {
	slice, err := mapstruct.SliceToInterfaceSlice(sliceInterface)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	return append(slice, item...), nil
}

func (receiver *coreUtil) Concatenate(sliceInterface1 interface{}, sliceInterface2 interface{}) (interface{}, error) {
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
}

func (receiver *coreUtil) FieldSlice(slice interface{}, fieldName string) ([]interface{}, error) {
	return mapstruct.StructSliceToFieldValueInterfaceSlice(slice, fieldName)
}

func (receiver *coreUtil) Sort(slice interface{}, order string) ([]interface{}, error) {
	return mapstruct.SortSlice(slice, order)
}

func (receiver *coreUtil) Map(pairs ...interface{}) (map[interface{}]interface{}, error) {
	if len(pairs)%2 == 1 {
		logging.Logger().Debug("numebr of arguments must be even")
		return nil, fmt.Errorf("numebr of arguments must be even")
	}
	m := make(map[interface{}]interface{}, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m, nil
}

func (receiver *coreUtil) Exists(target map[interface{}]interface{}, key interface{}) bool {
	_, exists := target[key]
	return exists
}

func (receiver *coreUtil) Put(target map[interface{}]interface{}, key interface{}, value interface{}) map[interface{}]interface{} {
	target[key] = value
	return target
}

func (receiver *coreUtil) Get(target map[interface{}]interface{}, key interface{}) interface{} {
	return target[key]
}

func (receiver *coreUtil) Delete(target map[interface{}]interface{}, key interface{}) map[interface{}]interface{} {
	delete(target, key)
	return target
}

func (receiver *coreUtil) Merge(source, destination map[interface{}]interface{}) map[interface{}]interface{} {
	for key, value := range source {
		destination[key] = value
	}
	return destination
}

func (receiver *coreUtil) Keys(target map[interface{}]interface{}) ([]interface{}, error) {
	return mapstruct.MapToKeySlice(target)
}

func (receiver *coreUtil) Hash(slice interface{}, keyField string) (map[interface{}]interface{}, error) {
	return mapstruct.StructSliceToInterfaceMap(slice, keyField)
}

func (receiver *coreUtil) SliceMap(slice interface{}, keyField string) (map[interface{}]interface{}, error) {
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
}

func (receiver *coreUtil) Sequence(begin, end int) interface{} {
	count := end - begin + 1
	result := make([]int, count)
	for i, j := 0, begin; i < count; i, j = i+1, j+1 {
		result[i] = j
	}
	return result
}

func (receiver *modelStore) Single(pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
	path := pathInterface.(string)

	controller, err := extension.GetAssociatedControllerWithPath(path)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	singleURL, err := controller.GetResourceSingleURL()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	parameters, err := extension.CreateParametersFromPathAntRoute(path, singleURL)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	urlValues, err := url.ParseQuery(queryInterface.(string))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	model := controller.GetModel()

	parameter, err := dbpkg.NewParameter(urlValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// single resets db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db := receiver.db.New()
	db, err = parameter.Paginate(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	db = parameter.SetPreloads(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(parameter.DefaultQuery(urlValues, "fields", "*"))
	queryFields := helper.QueryFields(model, fields)

	result, err := model.GetSingle(model, db, parameters, urlValues, queryFields)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	return result, nil
}

func (receiver *modelStore) Multi(pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
	path := pathInterface.(string)

	controller, err := extension.GetAssociatedControllerWithPath(path)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	multiURL, err := controller.GetResourceMultiURL()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	parameters, err := extension.CreateParametersFromPathAntRoute(path, multiURL)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	urlValues, err := url.ParseQuery(queryInterface.(string))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	model := controller.GetModel()

	parameter, err := dbpkg.NewParameter(urlValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// multi resets db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db := receiver.db.New()
	db, err = parameter.Paginate(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(parameter.DefaultQuery(urlValues, "fields", "*"))
	queryFields := helper.QueryFields(model, fields)
	result, err := model.GetMulti(model, db, parameters, urlValues, queryFields)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// reset all conditions in order to get the total number of records
	db = db.New()
	total, err := model.GetCount(model, db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// reset conditions except for limit and offset in order to get the record count before limitation
	db = db.New()
	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	countBeforePagination, err := model.GetCount(model, db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	type multiResult struct {
		Records               interface{}
		Total                 interface{}
		CountBeforePagination interface{}
	}

	multiResultObject := &multiResult{
		Records: result,
		Total:   total,
		CountBeforePagination: countBeforePagination,
	}

	return multiResultObject, nil
}

func (receiver *modelStore) First(pathInterface interface{}, queryInterface interface{}) (interface{}, error) {
	path := pathInterface.(string)

	controller, err := extension.GetAssociatedControllerWithPath(path)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	multiURL, err := controller.GetResourceMultiURL()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	parameters, err := extension.CreateParametersFromPathAntRoute(path, multiURL)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	urlValues, err := url.ParseQuery(queryInterface.(string))
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	model := controller.GetModel()

	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	parameter, err := dbpkg.NewParameter(urlValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// first resets db conditions like preloads, so you should use this method in GetSingle or GetMulti only,
	// and note that all conditions go away after this method.
	db := receiver.db.New()
	db, err = parameter.Paginate(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}
	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(parameter.DefaultQuery(urlValues, "fields", "*"))
	queryFields := helper.QueryFields(model, fields)
	result, err := model.GetMulti(model, db, parameters, urlValues, queryFields)
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
}

func (receiver *networkUtil) ParseCIDR(cidr string) (*network.Ipv4Address, error) {
	return network.ParseCIDR(cidr)
}

func init() {
	extension.RegisterModel(NewTemplateGeneration())
}
