package model

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"net/url"
)

// TemplateGenerationByName is the model class what represents template generation
type TemplateGenerationByName struct {
	TemplateGeneration
	Name string `json:"name" clay:"key_parameter"`
}

// NewTemplateGenerationByName creates a template generation model instance
func NewTemplateGenerationByName() *TemplateGenerationByName {
	return &TemplateGenerationByName{}
}

// GenerateTableName generates its resource/table name in URL/DB
func (receiver *TemplateGenerationByName) GenerateTableName(_ extension.Model, _ *gorm.DB) string {
	return "template_generations_by_name"
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *TemplateGenerationByName) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle generates text data based on registered templates
// parameters must be given as p[...]=...
func (receiver *TemplateGenerationByName) GetSingle(_ extension.Model, db *gorm.DB, parameters gin.Params, urlValues url.Values, _ string) (interface{}, error) {
	return receiver.GenerateTemplate(db, parameters, urlValues)
}

func init() {
	extension.RegisterModel(NewTemplateGenerationByName())
}
