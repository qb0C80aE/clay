package model

import (
	"github.com/qb0C80aE/clay/extension"
)

// Template is the model class what represents templates to generate texts from models
type Template struct {
	*Base             `json:"base,omitempty"`
	ID                int                 `json:"id" form:"id" gorm:"primary_key;auto_increment"`
	Name              string              `json:"name" form:"name" gorm:"not null;unique"`
	TemplateContent   string              `json:"template_content" form:"template_content"`
	Description       string              `json:"description" form:"description"`
	TemplateArguments []*TemplateArgument `json:"template_arguments"`
}

// NewTemplate creates a template model instance
func NewTemplate() *Template {
	return CreateModel(&Template{}).(*Template)
}

func init() {
	extension.RegisterModel(NewTemplate(), true)
	extension.RegisterDesignAccessor(NewTemplate())
}