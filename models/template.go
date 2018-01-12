package models

import (
	"github.com/qb0C80aE/clay/extensions"
)

// Template is the model class what represents templates to generate texts from models
type Template struct {
	ID                int                 `json:"id" form:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name              string              `json:"name" form:"name" gorm:"not null;unique"`
	TemplateContent   string              `json:"template_content" form:"template_content"`
	Description       string              `json:"description" form:"description"`
	TemplateArguments []*TemplateArgument `json:"template_arguments"`
}

// NewTemplateModel creates a Template model instance
func NewTemplateModel() *Template {
	return &Template{}
}

var sharedTemplateModel = NewTemplateModel()

// SharedTemplateModel returns the template model instance used as a model prototype and type analysis
func SharedTemplateModel() *Template {
	return sharedTemplateModel
}

func init() {
	extensions.RegisterModel(sharedTemplateModel, true)
}
