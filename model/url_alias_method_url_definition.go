package model

import (
	"github.com/qb0C80aE/clay/extension"
)

// URLAliasMethodURLTypeDefinition is the model class what represents url alias method-urltype definition
type URLAliasMethodURLTypeDefinition struct {
	Base
	Method        string `json:"method" yaml:"method" validate:"required,oneof=GET POST PUT DELETE PATCH OPTIONS"`
	TargetURLType string `json:"target_url_type" yaml:"target_url_type" validate:"required,oneof=single multi"`
	Accept        string `json:"accept" yaml:"accept"`
	AcceptCharset string `json:"accept_charset" yaml:"accept_charset"`
}

// NewURLAliasMethodURLTypeDefinition creates a template raw model instance
func NewURLAliasMethodURLTypeDefinition() *URLAliasMethodURLTypeDefinition {
	return &URLAliasMethodURLTypeDefinition{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *URLAliasMethodURLTypeDefinition) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

func init() {
	extension.RegisterModel(NewURLAliasMethodURLTypeDefinition())
}
