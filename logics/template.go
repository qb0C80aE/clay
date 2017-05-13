package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/utils/mapstruct"
	"net/url"
	"strconv"
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

func (logic *templateLogic) GetSingle(db *gorm.DB, id string, _ url.Values, queryFields string) (interface{}, error) {

	template := &models.Template{}

	if err := db.Select(queryFields).First(template, id).Error; err != nil {
		return nil, err
	}

	return template, nil

}

func (logic *templateLogic) GetMulti(db *gorm.DB, _ url.Values, queryFields string) (interface{}, error) {
	templates := []*models.Template{}

	if err := db.Select(queryFields).Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil

}

func (logic *templateLogic) Create(db *gorm.DB, _ url.Values, data interface{}) (interface{}, error) {
	template := data.(*models.Template)

	if err := db.Create(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func (logic *templateLogic) Update(db *gorm.DB, id string, _ url.Values, data interface{}) (interface{}, error) {
	template := data.(*models.Template)
	template.ID, _ = strconv.Atoi(id)

	if err := db.Save(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func (logic *templateLogic) Delete(db *gorm.DB, id string, _ url.Values) error {

	template := &models.Template{}

	if err := db.First(&template, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&template).Error; err != nil {
		return err
	}

	return nil

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
			template.TemplatePersistentParameters = nil
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
}
