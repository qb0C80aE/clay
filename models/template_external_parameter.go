package models

import (
	"github.com/qb0C80aE/clay/extensions"
)

type TemplateExternalParameter struct {
	ID         int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TemplateID int    `json:"template_id" gorm:"index" sql:"type:integer references templates(id) on delete cascade"`
	Name       string `json:"name"`
	Value      string `json:"value"`
}

func NewTemplateExternalParameterModel() *TemplateExternalParameter {
	return &TemplateExternalParameter{}
}

var sharedTemplateExternalParameterModel = NewTemplateExternalParameterModel()

func SharedTemplateExternalParameterModel() *TemplateExternalParameter {
	return sharedTemplateExternalParameterModel
}

func init() {
	extensions.RegisterModelType(sharedTemplateExternalParameterModel)
}
