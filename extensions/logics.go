package extensions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/logging"
	"net/url"
	"reflect"
)

var logicMap = map[reflect.Type]Logic{}

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
func RegisterLogic(model interface{}, logic Logic) {
	modelType := ModelType(model)
	logicMap[modelType] = logic
}

// RegisteredLogic returns the registered logic related to the given resource name
func RegisteredLogic(model interface{}) (Logic, error) {
	modelType := ModelType(model)
	result, exist := logicMap[modelType]
	if !exist {
		logging.Logger().Debugf("the logic related to given name %s does not exist", modelType.Name())
		return nil, fmt.Errorf("the logic related to given name %s does not exist", modelType.Name())
	}
	return result, nil
}
