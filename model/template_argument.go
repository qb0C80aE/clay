package model

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
)

const (
	// TemplateArgumentTypeInt indicates that a type of argument is integer
	TemplateArgumentTypeInt = "int"
	// TemplateArgumentTypeFloat indicates that a type of argument is float
	TemplateArgumentTypeFloat = "float"
	// TemplateArgumentTypeBool indicates that a type of argument is bool
	TemplateArgumentTypeBool = "bool"
	// TemplateArgumentTypeString indicates that a type of argument is string
	TemplateArgumentTypeString = "string"
)

// TemplateArgument is the model class what represents model-independent arguments used in templates
type TemplateArgument struct {
	Base
	ID           int       `json:"id" yaml:"id" gorm:"primary_key;auto_increment"`
	TemplateID   int       `json:"template_id" yaml:"template_id" gorm:"unique_index:template_id_name" sql:"type:integer references templates(id)"`
	Template     *Template `json:"template" yaml:"template" gorm:"ForeignKey:TemplateID"`
	Name         string    `json:"name" yaml:"name" gorm:"unique_index:template_id_name" validate:"required"`
	Description  string    `json:"description" yaml:"description" form:"description" sql:"type:text"`
	Type         string    `json:"type" yaml:"type" validate:"oneof=int float bool string"`
	DefaultValue string    `json:"default_value" yaml:"default_value" validate:"required"`
	ToBeDeleted  bool      `json:"to_be_deleted,omitempty" yaml:"to_be_deleted,omitempty" sql:"-"`
}

// NewTemplateArgument creates a template argument model instance
func NewTemplateArgument() *TemplateArgument {
	return &TemplateArgument{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *TemplateArgument) GetContainerForMigration() (interface{}, error) {
	return extension.CreateContainerByTypeName(receiver.GetTypeName(receiver))
}

// DoAfterDBMigration execute initialization process after DB migration
func (receiver *TemplateArgument) DoAfterDBMigration(db *gorm.DB) error {
	return db.Exec(`
		create trigger if not exists DeleteTemplateArgumentBeforeTemplateDeletion before delete on templates
		begin
			delete from template_arguments where template_id = old.id;
		end;
	`).Error
}

func init() {
	extension.RegisterModel(NewTemplateArgument())
	extension.RegisterInitializer(NewTemplateArgument())
	extension.RegisterDesignAccessor(NewTemplateArgument())
}
