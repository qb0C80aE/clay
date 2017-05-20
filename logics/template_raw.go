package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
	"net/url"
)

type templateRawLogic struct {
	*BaseLogic
}

func newTemplateRawLogic() *templateRawLogic {
	logic := &templateRawLogic{
		BaseLogic: &BaseLogic{},
	}
	return logic
}

func (logic *templateRawLogic) GetSingle(db *gorm.DB, id string, parameters url.Values, queryFields string) (interface{}, error) {
	template := &models.Template{}

	if err := db.Select(queryFields).First(template, id).Error; err != nil {
		return nil, err
	}

	return template.TemplateContent, nil
}

var uniqueTemplateRawLogic = newTemplateRawLogic()

// UniqueTemplateRawLogic returns the unique template logic instance
func UniqueTemplateRawLogic() extensions.Logic {
	return uniqueTemplateRawLogic
}

func init() {
}
