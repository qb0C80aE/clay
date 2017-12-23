package models

import "github.com/qb0C80aE/clay/extensions"

// TemplateRaw is the mock model class what represents template raws
type TemplateRaw struct {
}

// NewTemplateRawModel creates a TemplateRaw model instance
func NewTemplateRawModel() *TemplateRaw {
	return &TemplateRaw{}
}

var sharedTemplateRawModel = NewTemplateRawModel()

// SharedTemplateRawModel returns the template raw model instance used as a model prototype and type analysis
func SharedTemplateRawModel() *TemplateRaw {
	return sharedTemplateRawModel
}

func init() {
	extensions.RegisterModel(sharedTemplateRawModel, false)
}
