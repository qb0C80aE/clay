package extensions

import (
	"github.com/jinzhu/gorm"
	"net/url"
)

// Logic is the interface what handles business processes between controllers and models
// * GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
// * GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
// * Create corresponds HTTP POST message and handles a request for multi resource to create a new information
// * Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
// * Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
// * Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
// * Options corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
type Logic interface {
	GetMulti(db *gorm.DB, parameters url.Values, queryString string) (interface{}, error)
	GetSingle(db *gorm.DB, id string, parameters url.Values, queryString string) (interface{}, error)
	Create(db *gorm.DB, parameters url.Values, model interface{}) (interface{}, error)
	Update(db *gorm.DB, id string, parameters url.Values, model interface{}) (interface{}, error)
	Delete(db *gorm.DB, id string, parameters url.Values) error
	Patch(db *gorm.DB, id string, parameters url.Values) (interface{}, error)
	Options(db *gorm.DB, parameters url.Values) error
}
