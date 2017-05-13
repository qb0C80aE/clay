package logics

import (
	"bytes"
	"fmt"
	"github.com/jinzhu/gorm"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/helper"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"net/http"
	"net/url"
	"strconv"
	tplpkg "text/template"
)

type templateGenerationLogic struct {
	*BaseLogic
}

func newTemplateGenerationLogic() *templateGenerationLogic {
	logic := &templateGenerationLogic{
		BaseLogic: &BaseLogic{},
	}
	return logic
}

// GenerateTemplate generates text data based on registered templates
func GenerateTemplate(db *gorm.DB, id string, templateVolatileParameterMap map[string]interface{}) (interface{}, error) {
	templateParameter := map[string]interface{}{}

	template := &models.Template{}
	template.ID, _ = strconv.Atoi(id)

	if err := db.Preload("TemplatePersistentParameters").Select("*").First(template, template.ID).Error; err != nil {
		return nil, err
	}

	for key, value := range templateVolatileParameterMap {
		templateParameter[key] = value
	}
	for _, templatePersistentParameter := range template.TemplatePersistentParameters {
		templateParameter[templatePersistentParameter.Name] = templatePersistentParameter
	}
	templateParameter["ModelStore"] = db

	tpl := tplpkg.New("template")
	templateFuncMaps := extensions.RegisteredTemplateFuncMaps()
	for _, templateFuncMap := range templateFuncMaps {
		tpl = tpl.Funcs(templateFuncMap)
	}
	tpl, err := tpl.Parse(template.TemplateContent)
	if err != nil {
		return nil, err
	}

	var doc bytes.Buffer
	if err := tpl.Execute(&doc, templateParameter); err != nil {
		return nil, err
	}

	result := doc.String()

	return result, nil
}

// GetSingle generates text data based on registered templates
func (logic *templateGenerationLogic) GetSingle(db *gorm.DB, id string, parameters url.Values, queryFields string) (interface{}, error) {
	templateVolatileParameterMap := make(map[string]interface{}, len(parameters))
	for key, value := range parameters {
		templateVolatileParameterMap[key] = value
	}
	return GenerateTemplate(db, id, templateVolatileParameterMap)
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
		"slice": func(items ...interface{}) []interface{} {
			slice := []interface{}{}
			return append(slice, items...)
		},
		"subslice": func(slice []interface{}, begin int, end int) []interface{} {
			if begin < 0 {
				if end < 0 {
					return slice[:]
				}
				return slice[:end]
			}
			if end < 0 {
				return slice[begin:]
			}
			return slice[begin:end]
		},
		"append": func(slice []interface{}, item ...interface{}) []interface{} {
			return append(slice, item...)
		},
		"concatenate": func(slice1 []interface{}, slice2 []interface{}) []interface{} {
			return append(slice1, slice2...)
		},
		"map": func(pairs ...interface{}) (map[interface{}]interface{}, error) {
			if len(pairs)%2 == 1 {
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
		"hash": func(slice interface{}, keyField string) (map[interface{}]interface{}, error) {
			return mapstruct.SliceToInterfaceMap(slice, keyField)
		},
		"slicemap": func(slice interface{}, keyField string) (map[interface{}]interface{}, error) {
			sliceMap, err := mapstruct.SliceToInterfaceSliceMap(slice, keyField)
			if err != nil {
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
		"single": func(dbObject interface{}, nameInterface interface{}, idInterface interface{}, queryInterface interface{}) (interface{}, error) {
			name := nameInterface.(string)
			var id int
			id, ok := idInterface.(int)
			if !ok {
				id = int(idInterface.(int64))
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
				return nil, err
			}
			model, err := extensions.CreateModel(name)
			if err != nil {
				return nil, err
			}
			urlQuery := requestForParameter.URL.Query()
			parameter, err := dbpkg.NewParameter(urlQuery, model)
			if err != nil {
				return nil, err
			}
			templateParameterGenerator, err := extensions.RegisteredTemplateParameterGenerator(model)
			if err != nil {
				return nil, err
			}
			db := dbObject.(*gorm.DB)
			db, err = parameter.Paginate(db)
			if err != nil {
				return nil, err
			}
			db = parameter.SetPreloads(db)
			db = parameter.FilterFields(db)
			fields := helper.ParseFields(parameter.DefaultQuery(urlQuery, "fields", "*"))
			queryFields := helper.QueryFields(model, fields)
			result, err := templateParameterGenerator.GetSingle(db, strconv.Itoa(id), requestForParameter.URL.Query(), queryFields)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
		"multi": func(dbObject interface{}, nameInterface interface{}, queryInterface interface{}) (interface{}, error) {
			name := nameInterface.(string)
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
			model, err := extensions.CreateModel(name)
			if err != nil {
				return nil, err
			}
			urlQuery := requestForParameter.URL.Query()
			parameter, err := dbpkg.NewParameter(urlQuery, model)
			if err != nil {
				return nil, err
			}
			templateParameterGenerator, err := extensions.RegisteredTemplateParameterGenerator(model)
			if err != nil {
				return nil, err
			}
			db := dbObject.(*gorm.DB)
			db, err = parameter.Paginate(db)
			if err != nil {
				return nil, err
			}
			db = parameter.SetPreloads(db)
			db = parameter.SortRecords(db)
			db = parameter.FilterFields(db)
			fields := helper.ParseFields(parameter.DefaultQuery(urlQuery, "fields", "*"))
			queryFields := helper.QueryFields(model, fields)
			result, err := templateParameterGenerator.GetMulti(db, requestForParameter.URL.Query(), queryFields)
			if err != nil {
				return nil, err
			}
			return result, nil
		},
		"total": func(dbObject interface{}, nameInterface interface{}) (interface{}, error) {
			name := nameInterface.(string)
			model, err := extensions.CreateModel(name)
			if err != nil {
				return nil, err
			}
			db := dbObject.(*gorm.DB)
			var total = 0
			db.Model(model).Count(&total)
			return total, nil
		},
	}
	extensions.RegisterTemplateFuncMap(funcMap)
}
