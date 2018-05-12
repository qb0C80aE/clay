package model

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
	tplpkg "text/template"
)

// EphemeralTemplateGeneration is the model class what represents ephemeral template generation
type EphemeralTemplateGeneration struct {
	Base
	Name string `json:"name" yaml:"name" gorm:"primary_key"`
}

// NewEphemeralTemplateGeneration creates a ephemeral template generation model instance
func NewEphemeralTemplateGeneration() *EphemeralTemplateGeneration {
	return &EphemeralTemplateGeneration{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralTemplateGeneration) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralTemplateGeneration) GetSingle(model extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (interface{}, error) {
	modelKey, err := model.GetModelKey(model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	name := parameters.ByName(modelKey.KeyParameter)

	ephemeralTemplate, exists := nameEphemeralTemplateMap[name]

	if !exists {
		logging.Logger().Debug(errors.New("record not found"))
		return nil, errors.New("record not found")
	}

	tpl := tplpkg.New("ephemeral_template")

	templateFuncMaps := extension.GetRegisteredTemplateFuncMapList()
	for _, templateFuncMap := range templateFuncMaps {
		tpl = tpl.Funcs(templateFuncMap)
	}

	tpl, err = tpl.Parse(ephemeralTemplate.TemplateContent)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	templateParameter := &templateParameter{
		ModelStore: &modelStore{
			db: db,
		},
		Core:      &coreUtil{},
		Network:   &networkUtil{},
		Parameter: map[interface{}]interface{}{},
		Query:     urlValues,
	}

	var doc bytes.Buffer
	if err := tpl.Execute(&doc, templateParameter); err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	result := doc.String()

	return result, nil
}

func init() {
	extension.RegisterModel(NewEphemeralTemplateGeneration())
}
