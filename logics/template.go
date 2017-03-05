package logics

import (
	"bytes"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"strconv"
	tplpkg "text/template"
)

type TemplateExternalParameterLogic struct {
}

type TemplateLogic struct {
}

func (_ *TemplateExternalParameterLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	templateExternalParameter := &models.TemplateExternalParameter{}

	if err := db.Select(queryFields).First(templateExternalParameter, id).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func (_ *TemplateExternalParameterLogic) GetMulti(db *gorm.DB, queryFields string) ([]interface{}, error) {

	templateExternalParameters := []*models.TemplateExternalParameter{}

	if err := db.Select(queryFields).Find(&templateExternalParameters).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(templateExternalParameters))
	for i, data := range templateExternalParameters {
		result[i] = data
	}

	return result, nil

}

func (_ *TemplateExternalParameterLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {

	templateExternalParameter := data.(*models.TemplateExternalParameter)

	if err := db.Create(templateExternalParameter).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func (_ *TemplateExternalParameterLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	templateExternalParameter := data.(*models.TemplateExternalParameter)
	templateExternalParameter.ID, _ = strconv.Atoi(id)

	if err := db.Save(&templateExternalParameter).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func (_ *TemplateExternalParameterLogic) Delete(db *gorm.DB, id string) error {

	templateExternalParameter := &models.TemplateExternalParameter{}

	if err := db.First(&templateExternalParameter, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&templateExternalParameter).Error; err != nil {
		return err
	}

	return nil

}

func (_ *TemplateExternalParameterLogic) Patch(_ *gorm.DB, _ string, _ string) (interface{}, error) {
	return nil, nil
}

func (_ *TemplateExternalParameterLogic) Options(db *gorm.DB) error {
	return nil
}

func (_ *TemplateLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	template := &models.Template{}

	if err := db.Select(queryFields).First(template, id).Error; err != nil {
		return nil, err
	}

	return template, nil

}

func (_ *TemplateLogic) GetMulti(db *gorm.DB, queryFields string) ([]interface{}, error) {

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

func (_ *TemplateLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {
	template := data.(*models.Template)

	if err := db.Create(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func (_ *TemplateLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {
	template := data.(*models.Template)
	template.ID, _ = strconv.Atoi(id)

	if err := db.Save(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func (_ *TemplateLogic) Delete(db *gorm.DB, id string) error {

	template := &models.Template{}

	if err := db.First(&template, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&template).Error; err != nil {
		return err
	}

	return nil

}

func (_ *TemplateLogic) Patch(db *gorm.DB, id string, _ string) (interface{}, error) {
	design := &models.Design{
		Content: map[string]interface{}{},
	}

	designAccessors := extension.GetDesignAccessos()
	for _, accessor := range designAccessors {
		if err := accessor.ExtractFromDesign(db, design.Content); err != nil {
			return nil, err
		}
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

	design.Content["TemplateExternalParameters"] = templateExternalParameterMap

	tpl, _ := tplpkg.New("template").Parse(template.TemplateContent)

	var doc bytes.Buffer
	tpl.Execute(&doc, design)
	result := doc.String()

	return result, nil
}

func (_ *TemplateLogic) Options(db *gorm.DB) error {
	return nil
}

func (_ *TemplateExternalParameterLogic) ExtractFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	templateExternalParameters := []*models.TemplateExternalParameter{}
	if err := db.Select("*").Find(&templateExternalParameters).Error; err != nil {
		return err
	}
	designContent["template_external_parameters"] = templateExternalParameters
	return nil
}

func (_ *TemplateExternalParameterLogic) DeleteFromDesign(db *gorm.DB) error {
	return db.Exec("delete from template_external_parameters;").Error
}

func (_ *TemplateExternalParameterLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*models.TemplateExternalParameter{}
	design := data.(*models.Design)
	if value, exists := design.Content["template_external_parameters"]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		for _, template := range container {
			if err := db.Create(template).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (_ *TemplateLogic) ExtractFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	templates := []*models.Template{}
	if err := db.Select("*").Find(&templates).Error; err != nil {
		return err
	}
	designContent["templates"] = templates
	return nil
}

func (_ *TemplateLogic) DeleteFromDesign(db *gorm.DB) error {
	return db.Exec("delete from templates;").Error
}

func (_ *TemplateLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*models.Template{}
	design := data.(*models.Design)
	if value, exists := design.Content["templates"]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		for _, template := range container {
			if err := db.Create(template).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

var TemplateExternalParameterLogicInstance = &TemplateExternalParameterLogic{}
var TemplateLogicInstance = &TemplateLogic{}

func init() {
	extension.RegisterDesignAccessor(TemplateExternalParameterLogicInstance)
	extension.RegisterDesignAccessor(TemplateLogicInstance)
}
