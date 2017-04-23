package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"strconv"
)

type templateExternalParameterLogic struct {
	*BaseLogic
}

func newTemplateExternalParameterLogic() *templateExternalParameterLogic {
	logic := &templateExternalParameterLogic{
		BaseLogic: &BaseLogic{},
	}
	return logic
}

func (logic *templateExternalParameterLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	templateExternalParameter := &models.TemplateExternalParameter{}

	if err := db.Select(queryFields).First(templateExternalParameter, id).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func (logic *templateExternalParameterLogic) GetMulti(db *gorm.DB, queryFields string) (interface{}, error) {

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

func (logic *templateExternalParameterLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {

	templateExternalParameter := data.(*models.TemplateExternalParameter)

	if err := db.Create(templateExternalParameter).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func (logic *templateExternalParameterLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	templateExternalParameter := data.(*models.TemplateExternalParameter)
	templateExternalParameter.ID, _ = strconv.Atoi(id)

	if err := db.Save(&templateExternalParameter).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func (logic *templateExternalParameterLogic) Delete(db *gorm.DB, id string) error {

	templateExternalParameter := &models.TemplateExternalParameter{}

	if err := db.First(&templateExternalParameter, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&templateExternalParameter).Error; err != nil {
		return err
	}

	return nil

}

func (logic *templateExternalParameterLogic) ExtractFromDesign(db *gorm.DB) (string, interface{}, error) {
	templateExternalParameters := []*models.TemplateExternalParameter{}
	if err := db.Select("*").Find(&templateExternalParameters).Error; err != nil {
		return "", nil, err
	}
	return "template_external_parameters", templateExternalParameters, nil
}

func (logic *templateExternalParameterLogic) DeleteFromDesign(db *gorm.DB) error {
	return db.Exec("delete from template_external_parameters;").Error
}

func (logic *templateExternalParameterLogic) LoadToDesign(db *gorm.DB, data interface{}) error {
	container := []*models.TemplateExternalParameter{}
	design := data.(*models.Design)
	if value, exists := design.Content["template_external_parameters"]; exists {
		if err := mapstruct.MapToStruct(value.([]interface{}), &container); err != nil {
			return err
		}
		for _, templateExternalParameter := range container {
			if err := db.Create(templateExternalParameter).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

var uniqueTemplateExternalParameterLogic = newTemplateExternalParameterLogic()

// UniqueTemplateExternalParameterLogic returns the unique template external parameter logic instance
func UniqueTemplateExternalParameterLogic() extensions.Logic {
	return uniqueTemplateExternalParameterLogic
}

func init() {
	extensions.RegisterDesignAccessor(uniqueTemplateExternalParameterLogic)
	extensions.RegisterTemplateParameterGenerator("TemplateExternalParameter", uniqueTemplateExternalParameterLogic)
}
