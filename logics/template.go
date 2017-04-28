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
	"strconv"
	tplpkg "text/template"
)

type templateLogic struct {
	*BaseLogic
}

func newTemplateLogic() *templateLogic {
	logic := &templateLogic{
		BaseLogic: &BaseLogic{},
	}
	return logic
}

// GenerateTemplate generates text data based on registered templates
func GenerateTemplate(db *gorm.DB, id string, templateInternalParameterMap map[string]interface{}) (interface{}, error) {
	templateParameter := map[string]interface{}{}

	template := &models.Template{}
	template.ID, _ = strconv.Atoi(id)

	if err := db.Preload("TemplateExternalParameters").Select("*").First(template, template.ID).Error; err != nil {
		return nil, err
	}

	for key, value := range templateInternalParameterMap {
		templateParameter[key] = value
	}
	for _, templateExternalParameter := range template.TemplateExternalParameters {
		templateParameter[templateExternalParameter.Name] = templateExternalParameter
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

func (logic *templateLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	template := &models.Template{}

	if err := db.Select(queryFields).First(template, id).Error; err != nil {
		return nil, err
	}

	return template, nil

}

func (logic *templateLogic) GetMulti(db *gorm.DB, queryFields string) (interface{}, error) {
	templates := []*models.Template{}

	if err := db.Select(queryFields).Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil

}

func (logic *templateLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {
	template := data.(*models.Template)

	if err := db.Create(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func (logic *templateLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {
	template := data.(*models.Template)
	template.ID, _ = strconv.Atoi(id)

	if err := db.Save(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func (logic *templateLogic) Delete(db *gorm.DB, id string) error {

	template := &models.Template{}

	if err := db.First(&template, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&template).Error; err != nil {
		return err
	}

	return nil

}

// Patch generates text data based on registered templates
func (logic *templateLogic) Patch(db *gorm.DB, id string) (interface{}, error) {
	return GenerateTemplate(db, id, nil)
}

func (logic *templateLogic) ExtractFromDesign(db *gorm.DB) (string, interface{}, error) {
	templates := []*models.Template{}
	if err := db.Select("*").Find(&templates).Error; err != nil {
		return "", nil, err
	}
	return extensions.RegisteredResourceName(models.SharedTemplateModel()), templates, nil
}

func (logic *templateLogic) DeleteFromDesign(db *gorm.DB) error {
	return db.Delete(models.SharedTemplateModel()).Error
}

func (logic *templateLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*models.Template{}
	design := data.(*models.Design)
	if value, exists := design.Content[extensions.RegisteredResourceName(models.SharedTemplateModel())]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		for _, template := range container {
			template.TemplateExternalParameters = nil
			if err := db.Create(template).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

var uniqueTemplateLogic = newTemplateLogic()

// UniqueTemplateLogic returns the unique template logic instance
func UniqueTemplateLogic() extensions.Logic {
	return uniqueTemplateLogic
}

func init() {
	extensions.RegisterDesignAccessor(uniqueTemplateLogic)
	extensions.RegisterTemplateParameterGenerator(models.SharedTemplateModel(), uniqueTemplateLogic)

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
			result, err := templateParameterGenerator.GetSingle(db, strconv.Itoa(id), queryFields)
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
			result, err := templateParameterGenerator.GetMulti(db, queryFields)
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
