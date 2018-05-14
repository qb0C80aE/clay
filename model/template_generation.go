package model

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	collectionutilpkg "github.com/qb0C80aE/clay/util/collection"
	conversionutilpkg "github.com/qb0C80aE/clay/util/conversion"
	loggingutilpkg "github.com/qb0C80aE/clay/util/logging"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
	modelstorepkg "github.com/qb0C80aE/clay/util/modelstore"
	networkutilpkg "github.com/qb0C80aE/clay/util/network"
	stringutilpkg "github.com/qb0C80aE/clay/util/string"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	tplpkg "text/template"
)

var parameterRegexp = regexp.MustCompile("p\\[(.+)\\]")

type templateParameter struct {
	ModelStore         *modelstorepkg.ModelStore
	Collection         *collectionutilpkg.Utility
	Conversion         *conversionutilpkg.Utility
	MapStruct          *mapstructutilpkg.Utility
	Network            *networkutilpkg.Utility
	String             *stringutilpkg.Utility
	Parameter          map[interface{}]interface{}
	Query              url.Values
	ProgramInformation extension.ProgramInformation
	Logging            *loggingutilpkg.Utility
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

	if err := mapstructutilpkg.GetUtility().RemapToStruct(container, templateModelAsContainer); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	for _, templateArgument := range templateModelAsContainer.TemplateArguments {
		var err error
		templateArgumentMap[templateArgument.Name] = templateArgument
		switch templateArgument.Type {
		case TemplateArgumentTypeInt:
			templateParameterMap[templateArgument.Name], err = conversionutilpkg.GetUtility().Int64(templateArgument.DefaultValue)
		case TemplateArgumentTypeFloat:
			templateParameterMap[templateArgument.Name], err = conversionutilpkg.GetUtility().Float64(templateArgument.DefaultValue)
		case TemplateArgumentTypeBool:
			templateParameterMap[templateArgument.Name], err = conversionutilpkg.GetUtility().Boolean(templateArgument.DefaultValue)
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
		ModelStore:         modelstorepkg.NewModelStore(db),
		Collection:         collectionutilpkg.GetUtility(),
		Conversion:         conversionutilpkg.GetUtility(),
		MapStruct:          mapstructutilpkg.GetUtility(),
		Network:            networkutilpkg.GetUtility(),
		String:             stringutilpkg.GetUtility(),
		Parameter:          templateParameterMap,
		Query:              urlValues,
		ProgramInformation: extension.GetRegisteredProgramInformation(),
		Logging:            loggingutilpkg.GetUtility(),
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

func init() {
	extension.RegisterModel(NewTemplateGeneration())

	funcMap := tplpkg.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
		"mod": func(a, b int) int { return a % b },
	}
	extension.RegisterTemplateFuncMap(funcMap)
}
