package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
	"reflect"
)

// TemplateRawByName is the model class what represents template raw
type TemplateRawByName struct {
	TemplateRaw
	Name string `json:"name" clay:"key_parameter"`
}

// NewTemplateRawByName creates a template raw model instance
func NewTemplateRawByName() *TemplateRawByName {
	return &TemplateRawByName{}
}

// GenerateTableName generates its resource/table name in URL/DB
func (receiver *TemplateRawByName) GenerateTableName(_ extension.Model, _ *gorm.DB) string {
	return "template_raws_by_name"
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *TemplateRawByName) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle get text data of registered templates
func (receiver *TemplateRawByName) GetSingle(_ extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, _ string) (interface{}, error) {
	templateModel := NewTemplate()

	templateName, exists := parameters.Get("name")
	if !exists {
		logging.Logger().Debug("parameter name does not exist")
		return nil, errors.New("parameter name does not exist")
	}

	newURLValues := url.Values{}
	newURLValues.Set("q[name]", templateName)

	dbParameter, err := dbpkg.NewParameter(newURLValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	db = dbParameter.FilterFields(db)
	db = dbParameter.SetPreloads(db)

	container, err := templateModel.GetMulti(templateModel, db, parameters, newURLValues, "*")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	containerValue := reflect.ValueOf(container)
	if containerValue.Len() == 0 {
		logging.Logger().Debug("record not found")
		return nil, errors.New("record not found")
	}

	result := reflect.ValueOf(container).Index(0).Elem().FieldByName("TemplateContent").String()

	return result, nil
}

func init() {
	extension.RegisterModel(NewTemplateRawByName())
}
