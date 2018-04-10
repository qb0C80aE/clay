package model

import (
	"github.com/qb0C80aE/clay/extension"
)

// Template is the model class what represents templates to generate texts from models
type Template struct {
	Base
	ID                int                 `json:"id" form:"id" gorm:"primary_key;auto_increment"`
	Name              string              `json:"name" form:"name" gorm:"not null;unique"`
	TemplateContent   string              `json:"template_content" form:"template_content" sql:"type:text"`
	Description       string              `json:"description" form:"description" sql:"type:text"`
	TemplateArguments []*TemplateArgument `json:"template_arguments"`
}

// NewTemplate creates a template model instance
func NewTemplate() *Template {
	return ConvertContainerToModel(&Template{}).(*Template)
}

func init() {
	extension.RegisterModel(NewTemplate(), true)
	extension.RegisterDesignAccessor(NewTemplate())
}
