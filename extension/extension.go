package extension

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mohae/deepcopy"
	"net/http"
	"reflect"
)

const (
	MethodGet    = 1
	MethodPost   = 2
	MethodPut    = 3
	MethodDelete = 4
	MethodPatch  = 5
)

var typeMap = map[string]reflect.Type{}
var modelMap = map[reflect.Type]interface{}{}

var methodMap = map[int]map[string]gin.HandlerFunc{
	MethodGet:    {},
	MethodPost:   {},
	MethodPut:    {},
	MethodDelete: {},
	MethodPatch:  {},
}

var routePostInitializers = []func(*gin.Engine){}

func RegisterModelType(reflectType reflect.Type) {
	typeMap[reflectType.String()] = reflectType
	modelMap[reflectType] = reflect.New(reflectType).Elem().Interface()
}

func GetModels() []interface{} {
	result := []interface{}{}
	for _, model := range modelMap {
		result = append(result, model)
	}
	return result
}

var endPoints = map[string]string{}

func RegisterEndpoint(resourceName string) {
	resourceSingleTitle := fmt.Sprintf("%s_url", resourceName)
	resourceMultiTitle := fmt.Sprintf("%ss_url", resourceName)
	resourceSingleUrl := fmt.Sprintf("%ss/{id}", resourceName)
	resourceMultiUrl := fmt.Sprintf("%ss", resourceName)
	endPoints[resourceSingleTitle] = resourceSingleUrl
	endPoints[resourceMultiTitle] = resourceMultiUrl
}

func RegisterUniqueEndpoint(title, url string) {
	endPoints[title] = url
}

func APIEndpoints(c *gin.Context) {
	reqScheme := "http"

	if c.Request.TLS != nil {
		reqScheme = "https"
	}

	reqHost := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s/%s", reqScheme, reqHost, "v1")

	resources := map[string]string{}
	for title, url := range endPoints {
		resources[title] = fmt.Sprintf("%s/%s", baseURL, url)
	}

	c.IndentedJSON(http.StatusOK, resources)
}

func RegisterRoute(method int, relativePath string, handlerFunc gin.HandlerFunc) {
	methodMap[method][relativePath] = handlerFunc
}

func GetRoutes(method int) map[string]gin.HandlerFunc {
	return deepcopy.Copy(methodMap[method]).(map[string]gin.HandlerFunc)
}

func GetResourceSingleUrl(resourceName string) string {
	return fmt.Sprintf("/%ss/:id", resourceName)
}

func GetResourceMultiUrl(resourceName string) string {
	return fmt.Sprintf("/%ss", resourceName)
}

func RegisterRoutePostInitializer(initializer func(*gin.Engine)) {
	routePostInitializers = append(routePostInitializers, initializer)
}

func GetRoutePostInitializers() []func(*gin.Engine) {
	return deepcopy.Copy(routePostInitializers).([]func(*gin.Engine))
}
