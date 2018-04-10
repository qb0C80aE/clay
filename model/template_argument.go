package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/util/conversion"
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
	ID           int    `json:"id" gorm:"primary_key;auto_increment"`
	TemplateID   int    `json:"template_id" gorm:"unique_index:template_id_name" sql:"type:integer references templates(id)"`
	Name         string `json:"name" gorm:"unique_index:template_id_name"`
	Description  string `json:"description" form:"description" sql:"type:text"`
	Type         int    `json:"type"`
	DefaultValue string `json:"default_value"`
	ToBeDeleted  bool   `json:"to_be_deleted,omitempty" sql:"-"`
}

// NewTemplateArgument creates a template argument model instance
func NewTemplateArgument() *TemplateArgument {
	return ConvertContainerToModel(&TemplateArgument{}).(*TemplateArgument)
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

func (receiver *TemplateArgument) checkDefaultValueBeforeStore() error {
	var err error

	switch receiver.Type {
	case TemplateArgumentTypeInt:
		_, err = conversion.ToInt64Interface(receiver.DefaultValue)
	case TemplateArgumentTypeFloat:
		_, err = conversion.ToFloat64Interface(receiver.DefaultValue)
	case TemplateArgumentTypeBool:
		_, err = conversion.ToBooleanInterface(receiver.DefaultValue)
	case TemplateArgumentTypeString:
	default:
		err = fmt.Errorf("invalid type: %v", receiver.Type)
	}

	return err
}

// BeforeCreate is executed before db.Create with the model
func (receiver *TemplateArgument) BeforeCreate(tx *gorm.DB) error {
	if err := receiver.Base.BeforeCreate(tx); err != nil {
		return err
	}
	return receiver.checkDefaultValueBeforeStore()
}

// BeforeSave is executed before db.Save with the model
func (receiver *TemplateArgument) BeforeSave(tx *gorm.DB) error {
	if err := receiver.Base.BeforeSave(tx); err != nil {
		return err
	}
	return receiver.checkDefaultValueBeforeStore()
}

func init() {
	extension.RegisterModel(NewTemplateArgument(), true)
	extension.RegisterInitializer(NewTemplateArgument())
	extension.RegisterDesignAccessor(NewTemplateArgument())
}
