package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
)

const (
	// TemplateArgumentTypeInt indicates that a type of argument is integer
	TemplateArgumentTypeInt = 1
	// TemplateArgumentTypeFloat indicates that a type of argument is float
	TemplateArgumentTypeFloat = 2
	// TemplateArgumentTypeBool indicates that a type of argument is bool
	TemplateArgumentTypeBool = 3
	// TemplateArgumentTypeString indicates that a type of argument is string
	TemplateArgumentTypeString = 4
)

// TemplateArgument is the model class what represents model-independent arguments used in templates
type TemplateArgument struct {
	ID                 int             `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TemplateID         int             `json:"template_id" gorm:"unique_index:template_id_name" sql:"type:integer references templates(id)"`
	Name               string          `json:"name" gorm:"unique_index:template_id_name"`
	Description        string          `json:"description" form:"description"`
	Type               int             `json:"type"`
	DefaultValueInt    sql.NullInt64   `json:"default_value_int"`
	DefaultValueFloat  sql.NullFloat64 `json:"default_value_float"`
	DefaultValueBool   sql.NullBool    `json:"default_value_bool"`
	DefaultValueString sql.NullString  `json:"default_value_string"`
	ToBeDeleted        bool            `json:"to_be_deleted,omitempty" sql:"-"`
}

// NewTemplateArgumentModel creates a template argument model instance
func NewTemplateArgumentModel() *TemplateArgument {
	return &TemplateArgument{}
}

var sharedTemplateArgumentModel = NewTemplateArgumentModel()

// SharedTemplateArgumentModel returns the template argument model instance used as a model prototype and type analysis
func SharedTemplateArgumentModel() *TemplateArgument {
	return sharedTemplateArgumentModel
}

// SetupInitialData setups the initial data
func (templateArgument *TemplateArgument) SetupInitialData(db *gorm.DB) error {
	return db.Exec(`
		create trigger if not exists DeleteTemplateArgumentBeforeTemplateDeletion before delete on templates
		begin
			delete from template_arguments where template_id = old.id;
		end;
	`).Error
}

func init() {
	extensions.RegisterInitialDataLoader(sharedTemplateArgumentModel)
	extensions.RegisterModel(sharedTemplateArgumentModel, true)
}
