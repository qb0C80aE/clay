package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
)

// TemplateRaw is the model class what represents raw template
type TemplateRaw struct {
	Base
}

// NewTemplateRaw creates a template raw model instance
func NewTemplateRaw() *TemplateRaw {
	return ConvertContainerToModel(&TemplateRaw{}).(*TemplateRaw)
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *TemplateRaw) GetSingle(db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	template := NewTemplate()

	if err := db.Select(queryFields).First(template, parameters.ByName("id")).Error; err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return template.TemplateContent, nil
}

func init() {
	extension.RegisterModel(NewTemplateRaw(), false)
}
