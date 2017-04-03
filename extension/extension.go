package extension

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
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
	MethodGet:     http.MethodGet,
	MethodPost:    http.MethodPost,
	MethodPut:     http.MethodPut,
	MethodDelete:  http.MethodDelete,
	MethodPatch:   http.MethodPatch,
	MethodOptions: http.MethodOptions,
}

type Outputter interface {
	OutputError(c *gin.Context, code int, err error)
	OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{})
	OutputGetMulti(c *gin.Context, code int, result []interface{}, fields map[string]interface{})
	OutputCreate(c *gin.Context, code int, result interface{})
	OutputUpdate(c *gin.Context, code int, result interface{})
	OutputDelete(c *gin.Context, code int)
	OutputPatch(c *gin.Context, code int, result interface{})
	OutputOptions(c *gin.Context, code int)
}

type Controller interface {
	ResourceName() string
	RouteMap() map[int]map[string]gin.HandlerFunc
}

type Logic interface {
	GetMulti(db *gorm.DB, queryString string) ([]interface{}, error)
	GetSingle(db *gorm.DB, id string, queryString string) (interface{}, error)
	Create(db *gorm.DB, model interface{}) (interface{}, error)
	Update(db *gorm.DB, id string, model interface{}) (interface{}, error)
	Delete(db *gorm.DB, id string) error
	Patch(db *gorm.DB, id string) (interface{}, error)
	Options(db *gorm.DB) error
}

type RouterInitializer interface {
	InitializeEarly(*gin.Engine) error
	InitializeLate(*gin.Engine) error
}

type DesignAccessor interface {
	ExtractFromDesign(*gorm.DB) (string, interface{}, error)
	DeleteFromDesign(*gorm.DB) error
	LoadToDesign(*gorm.DB, interface{}) error
}

type TemplateParameterGenerator interface {
	GenerateTemplateParameter(*gorm.DB) (string, interface{}, error)
}

var typeMap = map[string]reflect.Type{}
var modelMap = map[reflect.Type]interface{}{}
var controllers = []Controller{}
var routerInitializers = []RouterInitializer{}
var designAccessors = []DesignAccessor{}
var templateParameterGenerators = []TemplateParameterGenerator{}
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

func RegisterTemplateParameterGenerator(templateParameterGenerator TemplateParameterGenerator) {
	templateParameterGenerators = append(templateParameterGenerators, templateParameterGenerator)
}

func GetTemplateParameterGenerators() []TemplateParameterGenerator {
	result := []TemplateParameterGenerator{}
	result = append(result, templateParameterGenerators...)
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
