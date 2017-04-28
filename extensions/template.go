package extensions

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"reflect"
	"text/template"
)

var templateParameterGeneratorMap = map[reflect.Type]TemplateParameterGenerator{}
var templateFuncMaps = []template.FuncMap{}

// TemplateParameterGenerator is the interface what extracts the models from db and generates parameters used in template processes
// * GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
// * GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
type TemplateParameterGenerator interface {
	GetMulti(db *gorm.DB, queryString string) (interface{}, error)
	GetSingle(db *gorm.DB, id string, queryString string) (interface{}, error)
}

// RegisterTemplateParameterGenerator registers a template parameter generator used in the template logic
func RegisterTemplateParameterGenerator(model interface{}, templateParameterGenerator TemplateParameterGenerator) {
	modelType := ModelType(model)
	templateParameterGeneratorMap[modelType] = templateParameterGenerator
}

// RegisteredTemplateParameterGenerator returns the registered template parameter generator related to given name
func RegisteredTemplateParameterGenerator(model interface{}) (TemplateParameterGenerator, error) {
	modelType := ModelType(model)
	result, exist := templateParameterGeneratorMap[modelType]
	if !exist {
		return nil, fmt.Errorf("the TemplateParameterGenerator related to given name %s does not exist", modelType.Name())
	}
	return result, nil
}

// RegisterTemplateFuncMap registers a template FuncMap used in the template logic
func RegisterTemplateFuncMap(templateFuncMap template.FuncMap) {
	templateFuncMaps = append(templateFuncMaps, templateFuncMap)
}

// RegisteredTemplateFuncMaps returns the registered template FuncMap
func RegisteredTemplateFuncMaps() []template.FuncMap {
	result := []template.FuncMap{}
	result = append(result, templateFuncMaps...)
	return result
}
