package extensions

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"text/template"
)

var templateParameterGeneratorMap = map[string]TemplateParameterGenerator{}
var templateFuncMaps = []template.FuncMap{}

// TemplateParameterGenerator is the interface what extracts the models from db and generates parameters used in template processes
// * GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
// * GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
type TemplateParameterGenerator interface {
	GetMulti(db *gorm.DB, queryString string) (interface{}, error)
	GetSingle(db *gorm.DB, id string, queryString string) (interface{}, error)
}

// RegisterTemplateParameterGenerator registers a template parameter generator used in the template logic
func RegisterTemplateParameterGenerator(name string, templateParameterGenerator TemplateParameterGenerator) {
	templateParameterGeneratorMap[name] = templateParameterGenerator
}

// RegisteredTemplateParameterGeneratorMap returns the map what contains registered template parameter generators
func RegisteredTemplateParameterGeneratorMap() map[string]TemplateParameterGenerator {
	result := map[string]TemplateParameterGenerator{}
	for key, value := range templateParameterGeneratorMap {
		result[key] = value
	}
	return result
}

// RegisteredTemplateParameterGenerator returns the registered template parameter generator related to given name
func RegisteredTemplateParameterGenerator(name string) (TemplateParameterGenerator, error) {
	result, exist := templateParameterGeneratorMap[name]
	if !exist {
		return nil, fmt.Errorf("the TemplateParameterGenerator related to given name %s does not exist", name)
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
