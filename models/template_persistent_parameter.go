package models

import (
	"database/sql"
	"github.com/qb0C80aE/clay/extensions"
)

// TemplatePersistentParameter is the model class what represents model-independent parameters used in templates
type TemplatePersistentParameter struct {
	ID          int             `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TemplateID  int             `json:"template_id" gorm:"unique_index:template_id_name" sql:"type:integer references templates(id) on delete cascade"`
	Name        string          `json:"name" gorm:"unique_index:template_id_name"`
	ValueInt    sql.NullInt64   `json:"value_int"`
	ValueFloat  sql.NullFloat64 `json:"value_float"`
	ValueBool   sql.NullBool    `json:"value_bool"`
	ValueString sql.NullString  `json:"value_string"`
}

// NewTemplatePersistentParameterModel creates a template persistent parameter model instance
func NewTemplatePersistentParameterModel() *TemplatePersistentParameter {
	return &TemplatePersistentParameter{}
}

var sharedTemplatePersistentParameterModel = NewTemplatePersistentParameterModel()

// SharedTemplatePersistentParameterModel returns the template persistent parameter model instance used as a model prototype and type analysis
func SharedTemplatePersistentParameterModel() *TemplatePersistentParameter {
	return sharedTemplatePersistentParameterModel
}

func init() {
	extensions.RegisterModel(sharedTemplatePersistentParameterModel)
}
