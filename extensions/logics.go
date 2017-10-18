package extensions

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/url"
)

var logics = []Logic{}
var logicMapByResourceName = map[string]Logic{}

// Logic is the interface what handles business processes between logics and models
// * ResourceName returns its target resource name
// * GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
// * GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
// * Create corresponds HTTP POST message and handles a request for multi resource to create a new information
// * Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
// * Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
// * Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
// * Options corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
// * Total returns the count of for multi resource
type Logic interface {
	ResourceName() string
	GetMulti(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryString string) (interface{}, error)
	GetSingle(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryString string) (interface{}, error)
	Create(db *gorm.DB, parameters gin.Params, urlValues url.Values, model interface{}) (interface{}, error)
	Update(db *gorm.DB, parameters gin.Params, urlValues url.Values, model interface{}) (interface{}, error)
	Delete(db *gorm.DB, parameters gin.Params, urlValues url.Values) error
	Patch(db *gorm.DB, parameters gin.Params, urlValues url.Values) (interface{}, error)
	Options(db *gorm.DB, parameters gin.Params, urlValues url.Values) error
	Total(db *gorm.DB, model interface{}) (int, error)
}

// RegisterLogic registers a logic used in the router
func RegisterLogic(logic Logic) {
	logics = append(logics, logic)
	logicMapByResourceName[logic.ResourceName()] = logic
}

// RegisteredLogics returns the registered logics
func RegisteredLogics() []Logic {
	result := []Logic{}
	result = append(result, logics...)
	return result
}

// RegisteredLogicByResourceName returns the registered logic related to the given resource name
func RegisteredLogicByResourceName(resourceName string) Logic {
	return logicMapByResourceName[resourceName]
}
