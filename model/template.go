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
	TemplateArguments []*TemplateArgument `json:"template_arguments" gorm:"ForeignKey:template_id"`
}

// NewTemplate creates a template model instance
func NewTemplate() *Template {
	return &Template{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *Template) GetContainerForMigration() (interface{}, error) {
	return extension.CreateContainerByTypeName(receiver.GetTypeName(receiver))
}

func init() {
	extension.RegisterModel(NewTemplate())
	extension.RegisterDesignAccessor(NewTemplate())
}
