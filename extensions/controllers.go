package extensions

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
