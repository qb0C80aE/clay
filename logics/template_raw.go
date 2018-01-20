package logics

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/models"
	"net/url"
)

type templateRawLogic struct {
	*BaseLogic
}

func newTemplateRawLogic() *templateRawLogic {
	logic := &templateRawLogic{
		BaseLogic: NewBaseLogic(
			models.SharedTemplateModel(),
		),
	}
	return logic
}

func (logic *templateRawLogic) GetSingle(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	template := &models.Template{}

	if err := db.Select(queryFields).First(template, parameters.ByName("id")).Error; err != nil {
		logging.Logger().Debug(err.Error())
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
