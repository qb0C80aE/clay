package logics

import (
	"bytes"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/mapstruct"
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

	result := make([]interface{}, len(templates))
	for i, data := range templates {
		result[i] = data
	}

	return result, nil

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
	templateParameter := map[string]interface{}{}

	templateParameterGenerators := extensions.RegisteredTemplateParameterGenerators()
	for _, generator := range templateParameterGenerators {
		key, value, err := generator.GenerateTemplateParameter(db)
		if err != nil {
			return nil, err
		}
		templateParameter[key] = value
	}

	template := &models.Template{}
	template.ID, _ = strconv.Atoi(id)

	if err := db.Preload("TemplateExternalParameters").Select("*").First(template, template.ID).Error; err != nil {
		return nil, err
	}

	templateExternalParameterMap := make(map[string]string)
	for _, templateExternalParameter := range template.TemplateExternalParameters {
		templateExternalParameterMap[templateExternalParameter.Name] = templateExternalParameter.Value
	}

	templateParameter["TemplateExternalParameters"] = templateExternalParameterMap

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

func (logic *templateLogic) ExtractFromDesign(db *gorm.DB) (string, interface{}, error) {
	templates := []*models.Template{}
	if err := db.Select("*").Find(&templates).Error; err != nil {
		return "", nil, err
	}
	return "templates", templates, nil
}

func (logic *templateLogic) DeleteFromDesign(db *gorm.DB) error {
	return db.Exec("delete from templates;").Error
}

func (logic *templateLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*models.Template{}
	design := data.(*models.Design)
	if value, exists := design.Content["templates"]; exists {
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
}
