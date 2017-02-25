package models

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"reflect"
)

type TemplateExternalParameter struct {
	ID         int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TemplateID int    `json:"template_id" gorm:"index" sql:"type:integer references templates(id) on delete cascade"`
	Name       string `json:"name"`
	Value      string `json:"value"`
}

type Template struct {
	ID                         int                          `json:"id" form:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name                       string                       `json:"name" form:"name"`
	TemplateContent            string                       `json:"template_content" form:"template_content"`
	TemplateExternalParameters []*TemplateExternalParameter `json:"template_external_parameters"`
}

func extractTemplateExternalParametersFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	templateExternalParameters := []*TemplateExternalParameter{}
	if err := db.Select("*").Find(&templateExternalParameters).Error; err != nil {
		return err
	}
	designContent["template_external_parameters"] = templateExternalParameters
	return nil
}

func deleteTemplateExternalParametersFromDesign(db *gorm.DB) error {
	return db.Exec("delete from template_external_parameters;").Error
}

func loadTemplateExternalParametersFromDesign(db *gorm.DB, data interface{}) error {
	container := []*TemplateExternalParameter{}
	design := data.(*Design)
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

func extractTemplatesFromDesign(db *gorm.DB, designContent map[string]interface{}) error {
	templates := []*Template{}
	if err := db.Select("*").Find(&templates).Error; err != nil {
		return err
	}
	designContent["templates"] = templates
	return nil
}

func deleteTemplatesFromDesign(db *gorm.DB) error {
	return db.Exec("delete from templates;").Error
}

func loadTemplatesFromDesign(db *gorm.DB, data interface{}) error {
	container := []*Template{}
	design := data.(*Design)
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

func init() {
	extension.RegisterModelType(reflect.TypeOf(TemplateExternalParameter{}))
	extension.RegisterModelType(reflect.TypeOf(Template{}))
	extension.RegisterDesignExtractor(extractTemplateExternalParametersFromDesign)
	extension.RegisterDesignExtractor(extractTemplatesFromDesign)
	extension.RegisterDesignDeleter(deleteTemplateExternalParametersFromDesign)
	extension.RegisterDesignDeleter(deleteTemplatesFromDesign)
	extension.RegisterDesignLoader(loadTemplateExternalParametersFromDesign)
	extension.RegisterDesignLoader(loadTemplatesFromDesign)
}

var TemplateExternalParameterModel = &TemplateExternalParameter{}
var TemplateModel = &Template{}
