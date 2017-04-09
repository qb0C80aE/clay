package extensions

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var controllers = []Controller{}

// MethodSomething integer constants are used as the key value of HTTP methods in maps
const (
	MethodGet     = 1
	MethodPost    = 2
	MethodPut     = 3
	MethodDelete  = 4
	MethodPatch   = 5
	MethodOptions = 6
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
// * RouteMap returns the map its key is resource url and its value is request handler of the controller
type Controller interface {
	ResourceName() string
	RouteMap() map[int]map[string]gin.HandlerFunc
}

// LookUpMethodName returns the HTTP method name string corresponded to the argument
func LookUpMethodName(method int) string {
	return methodStringMap[method]
}

// BuildResourceSingleURL builds a resource url what represents a single resource based on the argument
func BuildResourceSingleURL(resourceName string) string {
	return fmt.Sprintf("/%ss/:id", resourceName)
}

// BuildResourceMultiURL builds a resource url what represents multi resources based on the argument
func BuildResourceMultiURL(resourceName string) string {
	return fmt.Sprintf("/%ss", resourceName)
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
