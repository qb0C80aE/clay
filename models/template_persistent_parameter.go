package models

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
)

// TemplatePersistentParameter is the model class what represents model-independent parameters used in templates
type TemplatePersistentParameter struct {
	ID          int             `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	TemplateID  int             `json:"template_id" gorm:"unique_index:template_id_name" sql:"type:integer references templates(id)"`
	Name        string          `json:"name" gorm:"unique_index:template_id_name"`
	Description string          `json:"description" form:"description"`
	ValueInt    sql.NullInt64   `json:"value_int"`
	ValueFloat  sql.NullFloat64 `json:"value_float"`
	ValueBool   sql.NullBool    `json:"value_bool"`
	ValueString sql.NullString  `json:"value_string"`
	ToBeDeleted bool            `json:"to_be_deleted,omitempty" sql:"-"`
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

// SetupInitialData setups the initial data
func (templatePersistentParameter *TemplatePersistentParameter) SetupInitialData(db *gorm.DB) error {
	return db.Exec(`
		create trigger if not exists DeleteTemplatePersistentParameterBeforeTemplateDeletion before delete on templates
		begin
			delete from template_persistent_parameters where template_id = old.id;
		end;
	`).Error
}

func init() {
	extensions.RegisterInitialDataLoader(sharedTemplatePersistentParameterModel)
	extensions.RegisterModel(sharedTemplatePersistentParameterModel, true)
}
