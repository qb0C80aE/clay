package extension

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logging"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type regExpControllerPair struct {
	regExp     *regexp.Regexp
	controller Controller
}

var controllerList = []Controller{}
var pathControllerMap = map[string]*regExpControllerPair{}

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
// * GetModel returns its model
// * GetResourceName returns its resource name in REST
// * GetResourceSingleURL builds a resource url what represents a single resource based on the argument
// * GetResourceMultiURL builds a resource url what represents multi resources based on the argument
// * GetRouteMap returns the map its key is resource url and its value is request handler of the controller
type Controller interface {
	GetModel() Model
	GetResourceName() string
	GetResourceSingleURL() string
	GetResourceMultiURL() string
	GetRouteMap() map[int]map[int]gin.HandlerFunc
}

// Binder is the interface what handles binding of input data
// * Bind binds input data to model instance
type Binder interface {
	Bind(c *gin.Context, container interface{}) error
}

// QueryCustomizer is the interface what handles query parameters used as Parameter struct in GetSingle and GetMulti
// * GetQueries returns query parameters
type QueryCustomizer interface {
	GetQueries(c *gin.Context) url.Values
}

// OutputHandler is the interface what handles outputs of the results from the logic to controllers
// * OutputError handles an error output
// * OutputGetSingle corresponds HTTP GET message and handles the output of a single result from logic classes
// * OutputGetMulti corresponds HTTP GET message and handles the output of multiple result from logic classes
// * OutputCreate corresponds HTTP POST message and handles the output of a single result from logic classes
// * OutputUpdate corresponds HTTP PUT message and handles the output of a single result from logic classes
// * OutputDelete corresponds HTTP DELETE message and handles the code result from logic classes
// * OutputPatch corresponds HTTP PATCH message and handles the output of a single result from logic classes
// * OutputGetOptions corresponds HTTP DELETE message and handles the code result from logic classes, as well as OutputDelete
type OutputHandler interface {
	OutputError(c *gin.Context, code int, err error)
	OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{})
	OutputGetMulti(c *gin.Context, code int, result interface{}, total int, fields map[string]interface{})
	OutputCreate(c *gin.Context, code int, result interface{})
	OutputUpdate(c *gin.Context, code int, result interface{})
	OutputDelete(c *gin.Context, code int)
	OutputPatch(c *gin.Context, code int, result interface{})
	OutputGetOptions(c *gin.Context, code int)
}

// LookUpMethodName returns the HTTP method name string corresponded to the argument
func LookUpMethodName(method int) string {
	return methodStringMap[method]
}

// RegisterController registers a controller used in the router
func RegisterController(controller Controller) {
	controllerList = append(controllerList, controller)
}

// GetRegisteredControllerList returns the registered controllers
func GetRegisteredControllerList() []Controller {
	result := []Controller{}
	result = append(result, controllerList...)
	return result
}

// AssociateControllerWithPath associates a path and a controller
func AssociateControllerWithPath(path string, controller Controller) {
	pathElementList := strings.Split(strings.Trim(path, "/"), "/")

	newPathElementList := []string{}
	for _, pathElement := range pathElementList {
		if pathElement[:1] == ":" {
			newPathElementList = append(newPathElementList, "[0-9a-z_]+")
		} else {
			newPathElementList = append(newPathElementList, pathElement)
		}
	}

	newPath := strings.Join(newPathElementList, "/")

	if oldPair, exists := pathControllerMap[newPath]; exists {
		if oldPair.controller != controller {
			logging.Logger().Criticalf("path %s is already mapped to another controller", newPath)
			os.Exit(1)
		}
	}

	newRegExp, err := regexp.Compile(newPath)
	if err != nil {
		logging.Logger().Critical(err.Error())
		os.Exit(1)
	}

	pathControllerMap[newPath] = &regExpControllerPair{
		regExp:     newRegExp,
		controller: controller,
	}
}

// GetAssociatedControllerWithPath returns the registered controller related to the given resource name
func GetAssociatedControllerWithPath(path string) (Controller, error) {
	path = strings.Trim(path, "/")
	for _, pair := range pathControllerMap {
		if pair.regExp.MatchString(path) {
			return pair.controller, nil
		}
	}
	logging.Logger().Critical("no controller is associated with %s", path)
	return nil, fmt.Errorf("no controller is associated with %s", path)
}
