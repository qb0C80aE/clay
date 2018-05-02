package model

import "github.com/qb0C80aE/clay/extension"

// URLAlias is the model class what represents urlAlias to execute something
type URLAlias struct {
	Base
}

// NewURLAlias creates a urlAlias model instance
func NewURLAlias() *URLAlias {
	return &URLAlias{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *URLAlias) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

func init() {
	extension.RegisterModel(NewURLAlias())
}
