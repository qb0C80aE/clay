package extension

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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

var routerPreInitializers = []func(*gin.Engine){}
var routerPostInitializers = []func(*gin.Engine){}
var endPoints = map[string]string{}
var designExtractors []func(*gorm.DB, map[string]interface{}) error = []func(*gorm.DB, map[string]interface{}) error{}
var designDeleters []func(*gorm.DB) error = []func(*gorm.DB) error{}
var designLoaders []func(*gorm.DB, interface{}) error = []func(*gorm.DB, interface{}) error{}

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

func RegisterRouterPreInitializer(initializer func(*gin.Engine)) {
	routerPreInitializers = append(routerPreInitializers, initializer)
}

func GetRouterPreInitializers() []func(*gin.Engine) {
	return deepcopy.Copy(routerPreInitializers).([]func(*gin.Engine))
}

func RegisterRouterPostInitializer(initializer func(*gin.Engine)) {
	routerPostInitializers = append(routerPostInitializers, initializer)
}

func GetRouterPostInitializers() []func(*gin.Engine) {
	return deepcopy.Copy(routerPostInitializers).([]func(*gin.Engine))
}

func RegisterDesignExtractor(extractor func(*gorm.DB, map[string]interface{}) error) {
	designExtractors = append(designExtractors, extractor)
}

func RegisterDesignDeleter(deleter func(*gorm.DB) error) {
	designDeleters = append(designDeleters, deleter)
}

func RegisterDesignLoader(loader func(*gorm.DB, interface{}) error) {
	designLoaders = append(designLoaders, loader)
}

func GetDesignExtractors() []func(*gorm.DB, map[string]interface{}) error {
	return deepcopy.Copy(designExtractors).([]func(*gorm.DB, map[string]interface{}) error)
}

func GetDesignDeleters() []func(*gorm.DB) error {
	return deepcopy.Copy(designDeleters).([]func(*gorm.DB) error)
}

func GetDesignLoaders() []func(*gorm.DB, interface{}) error {
	return deepcopy.Copy(designLoaders).([]func(*gorm.DB, interface{}) error)
}
