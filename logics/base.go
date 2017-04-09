package logics

import (
	"github.com/jinzhu/gorm"
)

// BaseLogic is the base class that all logic classes inherit
type BaseLogic struct {
}

func newBaseLogic() *BaseLogic {
	return &BaseLogic{}
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (logic *BaseLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {
	return nil, nil
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (logic *BaseLogic) GetMulti(db *gorm.DB, queryFields string) ([]interface{}, error) {
	return nil, nil
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (logic *BaseLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {
	return nil, nil
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (logic *BaseLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {
	return nil, nil
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (logic *BaseLogic) Delete(db *gorm.DB, id string) error {
	return nil
}

// Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
func (logic *BaseLogic) Patch(db *gorm.DB, id string) (interface{}, error) {
	return nil, nil
}

// Options corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
func (logic *BaseLogic) Options(db *gorm.DB) error {
	return nil
}

var baseLogic = newBaseLogic()

// SharedBaseLogic returns the base logic instance used as the embedded instance in all logic classes
func SharedBaseLogic() *BaseLogic {
	return baseLogic
}

func init() {
}
