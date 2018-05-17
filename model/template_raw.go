package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
	"net/url"
)

// TemplateRaw is the model class what represents raw template
type TemplateRaw struct {
	Base
}

// NewTemplateRaw creates a template raw model instance
func NewTemplateRaw() *TemplateRaw {
	return &TemplateRaw{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *TemplateRaw) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *TemplateRaw) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	templateModel := NewTemplate()
	templateModelAsContainer := NewTemplate()

	db = db.New()

	newURLValues := url.Values{}
	if len(urlValues.Get("key_parameter")) > 0 {
		newURLValues.Set("key_parameter", urlValues.Get("key_parameter"))
	}

	container, err := templateModel.GetSingle(templateModel, db, parameters, newURLValues, "*")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	if err := mapstructutilpkg.GetUtility().MapToStruct(container, templateModelAsContainer); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return templateModelAsContainer.TemplateContent, nil
}

func init() {
	extension.RegisterModel(NewTemplateRaw())
}
