package extension

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"reflect"
	"text/template"
)

const (
	MethodGet     = 1
	MethodPost    = 2
	MethodPut     = 3
	MethodDelete  = 4
	MethodPatch   = 5
	MethodOptions = 6
)

var methodNameMap = map[int]string{
	MethodGet:     "GET",
	MethodPost:    "POST",
	MethodPut:     "PUT",
	MethodDelete:  "DELETE",
	MethodPatch:   "PATCH",
	MethodOptions: "OPTIONS",
}

type Controller interface {
	Initialize()
	GetResourceName() string
	GetRouteMap() map[int]map[string]gin.HandlerFunc
}

type Logic interface {
	GetMulti(*gorm.DB, string) ([]interface{}, error)
	GetSingle(*gorm.DB, string, string) (interface{}, error)
	Create(*gorm.DB, interface{}) (interface{}, error)
	Update(*gorm.DB, string, interface{}) (interface{}, error)
	Delete(*gorm.DB, string) error
	Patch(*gorm.DB, string, string) (interface{}, error)
	Options(*gorm.DB) error
}

type RouterInitializer interface {
	InitializeEarly(*gin.Engine) error
	InitializeLate(*gin.Engine) error
}

type DesignAccessor interface {
	ExtractFromDesign(*gorm.DB, map[string]interface{}) error
	DeleteFromDesign(*gorm.DB) error
	LoadToDesign(*gorm.DB, interface{}) error
}

var typeMap = map[string]reflect.Type{}
var modelMap = map[reflect.Type]interface{}{}
var controllers = []Controller{}
var routerInitializers = []RouterInitializer{}
var designAccessors = []DesignAccessor{}
var templateFuncMaps = []template.FuncMap{}

func GetMethodName(method int) string {
	return methodNameMap[method]
}

func RegisterModelType(model interface{}) {
	reflectType := reflect.TypeOf(model)
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

func GetResourceSingleUrl(resourceName string) string {
	return fmt.Sprintf("/%ss/:id", resourceName)
}

func GetResourceMultiUrl(resourceName string) string {
	return fmt.Sprintf("/%ss", resourceName)
}

func RegisterController(controller Controller) {
	controllers = append(controllers, controller)
}

func GetControllers() []Controller {
	result := []Controller{}
	result = append(result, controllers...)
	return result
}

func RegisterRouterInitializer(initializer RouterInitializer) {
	routerInitializers = append(routerInitializers, initializer)
}

func GetRouterInitializers() []RouterInitializer {
	result := []RouterInitializer{}
	result = append(result, routerInitializers...)
	return result
}

func RegisterDesignAccessor(designAccessor DesignAccessor) {
	designAccessors = append(designAccessors, designAccessor)
}

func GetDesignAccessos() []DesignAccessor {
	result := []DesignAccessor{}
	result = append(result, designAccessors...)
	return result
}

func RegisterTemplateFuncMap(templateFuncMap template.FuncMap) {
	templateFuncMaps = append(templateFuncMaps, templateFuncMap)
}

func GetTemplateFuncMaps() []template.FuncMap {
	result := []template.FuncMap{}
	result = append(result, templateFuncMaps...)
	return result
}
