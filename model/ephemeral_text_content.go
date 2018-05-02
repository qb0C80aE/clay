package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"net/url"
)

// EphemeralTextContent is the model class what represents ephemeral text content
type EphemeralTextContent struct {
	Base
}

// NewEphemeralTextContent creates a ephemeral text content model instance
func NewEphemeralTextContent() *EphemeralTextContent {
	return &EphemeralTextContent{}
}

// GetContainerForMigration returns its container for migration, if no need to be migrated, just return null
func (receiver *EphemeralTextContent) GetContainerForMigration() (interface{}, error) {
	return nil, nil
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *EphemeralTextContent) GetSingle(_ extension.Model, db *gorm.DB, parameters gin.Params, _ url.Values, queryFields string) (interface{}, error) {
	name := parameters.ByName("name")

	result, exists := nameEphemeralTextMap[name]

	if !exists {
		return nil, errors.New("record not found")
	}

	return result.Content, nil
}

func init() {
	extension.RegisterModel(NewEphemeralTextContent())
}
