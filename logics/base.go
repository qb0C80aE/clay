package logics

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"net/url"
)

// BaseLogic is the base class that all logic classes inherit
type BaseLogic struct {
	model interface{}
}

// NewBaseLogic creates a new instance of BaseLogic
func NewBaseLogic(model interface{}) *BaseLogic {
	logic := &BaseLogic{
		model: model,
	}
	return logic
}

// ResourceName returns its target resource name
func (logic *BaseLogic) ResourceName() string {
	return extensions.RegisteredResourceName(logic.model)
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (logic *BaseLogic) GetSingle(_ *gorm.DB, _ gin.Params, _ url.Values, _ string) (interface{}, error) {
	return nil, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (logic *BaseLogic) GetMulti(_ *gorm.DB, _ gin.Params, _ url.Values, _ string) (interface{}, error) {
	return nil, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (logic *BaseLogic) Create(_ *gorm.DB, _ gin.Params, _ url.Values, _ interface{}) (interface{}, error) {
	return nil, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (logic *BaseLogic) Update(_ *gorm.DB, _ gin.Params, _ url.Values, _ interface{}) (interface{}, error) {
	return nil, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (logic *BaseLogic) Delete(_ *gorm.DB, _ gin.Params, _ url.Values) error {
	return nil
}

// Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
func (logic *BaseLogic) Patch(_ *gorm.DB, _ gin.Params, _ url.Values) (interface{}, error) {
	return nil, nil
}

// Options corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
func (logic *BaseLogic) Options(_ *gorm.DB, _ gin.Params, _ url.Values) error {
	return nil
}

// Total returns the count of for multi resource
func (logic *BaseLogic) Total(db *gorm.DB, model interface{}) (int, error) {
	var total int
	if err := db.Model(model).Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func init() {
}
