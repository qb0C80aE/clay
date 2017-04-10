package models

import (
	"github.com/qb0C80aE/clay/extensions"
)

// TemplateExternalParameter is the model class what represents model-independent parameters used in templates
type TemplateExternalParameter struct {
	ID         int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TemplateID int    `json:"template_id" gorm:"index" sql:"type:integer references templates(id) on delete cascade"`
	Name       string `json:"name"`
	Value      string `json:"value"`
}

// NewTemplateExternalParameterModel creates a Template external parameter model instance
func NewTemplateExternalParameterModel() *TemplateExternalParameter {
	return &TemplateExternalParameter{}
}

var sharedTemplateExternalParameterModel = NewTemplateExternalParameterModel()

// SharedTemplateExternalParameterModel returns the template external parameter model instance used as a model prototype and type analysis
func SharedTemplateExternalParameterModel() *TemplateExternalParameter {
	return sharedTemplateExternalParameterModel
}

func init() {
	extensions.RegisterModel(sharedTemplateExternalParameterModel)
}
