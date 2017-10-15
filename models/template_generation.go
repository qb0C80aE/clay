package models

import "github.com/qb0C80aE/clay/extensions"

// TemplateGeneration is the mock model class what represents template generations
type TemplateGeneration struct {
	ID int `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
}

// NewTemplateGenerationModel creates a TemplateGeneration model instance
func NewTemplateGenerationModel() *TemplateGeneration {
	return &TemplateGeneration{}
}

var sharedTemplateGenerationModel = NewTemplateGenerationModel()

// SharedTemplateGenerationModel returns the template generation model instance used as a model prototype and type analysis
func SharedTemplateGenerationModel() *TemplateGeneration {
	return sharedTemplateGenerationModel
}

func init() {
	extensions.RegisterModel(sharedTemplateGenerationModel, true)
}
