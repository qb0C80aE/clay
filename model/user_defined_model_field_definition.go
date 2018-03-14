package model

import (
	"github.com/qb0C80aE/clay/extension"
)

// UserDefinedModelFieldDefinition is the model class what represents raw template
type UserDefinedModelFieldDefinition struct {
	Base
	Name     string `json:"name" binding:"required"`
	Tag      string `json:"tag"`
	TypeName string `json:"type_name" binding:"required"`
	IsSlice  bool   `json:"is_slice"`
}

// NewUserDefinedModelFieldDefinition creates a template raw model instance
func NewUserDefinedModelFieldDefinition() *UserDefinedModelFieldDefinition {
	return &UserDefinedModelFieldDefinition{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *UserDefinedModelFieldDefinition) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

func init() {
	extension.RegisterModel(NewUserDefinedModelFieldDefinition())
}
