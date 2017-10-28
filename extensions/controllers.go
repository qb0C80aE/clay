package extensions

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
)

var controllers = []Controller{}
var controllerMapByResourceName = map[string]Controller{}

// MethodSomething integer constants are used as the key value of HTTP methods in maps
const (
	MethodGet     = 1
	MethodPost    = 2
	MethodPut     = 3
	MethodDelete  = 4
	MethodPatch   = 5
	MethodOptions = 6
	URLSingle     = 1
	URLMulti      = 2
)

var methodStringMap = map[int]string{
	MethodGet:     http.MethodGet,
	MethodPost:    http.MethodPost,
	MethodPut:     http.MethodPut,
	MethodDelete:  http.MethodDelete,
	MethodPatch:   http.MethodPatch,
	MethodOptions: http.MethodOptions,
}

// Controller is the interface what is mapped into specific urls and handles the requests from HTTP clients
// * ResourceName returns its resource name in REST
// * ResourceSingleURL builds a resource url what represents a single resource based on the argument
// * ResourceMultiURL builds a resource url what represents multi resources based on the argument
// * RouteMap returns the map its key is resource url and its value is request handler of the controller
type Controller interface {
	ResourceName() string
	ResourceSingleURL() string
	ResourceMultiURL() string
	RouteMap() map[int]map[int]gin.HandlerFunc
}

// QueryCustomizer is the interface what handles query parameters used as Parameter struct in GetSingle and GetMulti
// * Queries returns query parameters
type QueryCustomizer interface {
	Queries(c *gin.Context) url.Values
}

// Outputter is the interface what handles outputs of the results from the logic to controllers
// * OutputError handles an error output
// * OutputGetSingle corresponds HTTP GET message and handles the output of a single result from logic classes
// * OutputGetMulti corresponds HTTP GET message and handles the output of multiple result from logic classes
// * OutputCreate corresponds HTTP POST message and handles the output of a single result from logic classes
// * OutputUpdate corresponds HTTP PUT message and handles the output of a single result from logic classes
// * OutputDelete corresponds HTTP DELETE message and handles the code result from logic classes
// * OutputPatch corresponds HTTP PATCH message and handles the output of a single result from logic classes
// * OutputOptions corresponds HTTP DELETE message and handles the code result from logic classes, as well as OutputDelete
type Outputter interface {
	OutputError(c *gin.Context, code int, err error)
	OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{})
	OutputGetMulti(c *gin.Context, code int, result interface{}, total int, fields map[string]interface{})
	OutputCreate(c *gin.Context, code int, result interface{})
	OutputUpdate(c *gin.Context, code int, result interface{})
	OutputDelete(c *gin.Context, code int)
	OutputPatch(c *gin.Context, code int, result interface{})
	OutputOptions(c *gin.Context, code int)
}

// LookUpMethodName returns the HTTP method name string corresponded to the argument
func LookUpMethodName(method int) string {
	return methodStringMap[method]
}

// RegisterController registers a controller used in the router
func RegisterController(controller Controller) {
	controllers = append(controllers, controller)
}

// RegisteredControllers returns the registered controllers
func RegisteredControllers() []Controller {
	result := []Controller{}
	result = append(result, controllers...)
	return result
}

// AssociateControllerWithResourceName associates a resource name and a controller
func AssociateControllerWithResourceName(resourceName string, controller Controller) {
	controllerMapByResourceName[resourceName] = controller
}

// AssociatedControllerWithResourceName returns the registered controller related to the given resource name
func AssociatedControllerWithResourceName(resourceName string) Controller {
	return controllerMapByResourceName[resourceName]
}
