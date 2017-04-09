package models

import (
	"github.com/qb0C80aE/clay/extensions"
)

type Template struct {
	ID                         int                          `json:"id" form:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name                       string                       `json:"name" form:"name"`
	TemplateContent            string                       `json:"template_content" form:"template_content"`
	TemplateExternalParameters []*TemplateExternalParameter `json:"template_external_parameters"`
}

func NewTemplateModel() *Template {
	return &Template{}
}

var sharedTemplateModel = NewTemplateModel()

func SharedTemplateModel() *Template {
	return sharedTemplateModel
}

func init() {
	extensions.RegisterModelType(sharedTemplateModel)
}
