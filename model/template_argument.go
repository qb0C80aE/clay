package model

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
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
	Base
	ID           int       `json:"id" gorm:"primary_key;auto_increment"`
	TemplateID   int       `json:"template_id" gorm:"unique_index:template_id_name" sql:"type:integer references templates(id)"`
	Template     *Template `json:"template" gorm:"ForeignKey:TemplateID"`
	Name         string    `json:"name" gorm:"unique_index:template_id_name"`
	Description  string    `json:"description" form:"description" sql:"type:text"`
	Type         int       `json:"type" binding:"required,min=1,max=4"`
	DefaultValue string    `json:"default_value"`
	ToBeDeleted  bool      `json:"to_be_deleted,omitempty" sql:"-"`
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
