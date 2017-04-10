package extensions

import (
	"github.com/jinzhu/gorm"
	"text/template"
)

var templateParameterGenerators = []TemplateParameterGenerator{}
var templateFuncMaps = []template.FuncMap{}

// TemplateParameterGenerator is the interface what extracts the models from db and generates parameters used in template processes
// * GenerateTemplateParameter generates parameters based on specific models used in templates
type TemplateParameterGenerator interface {
	GenerateTemplateParameter(db *gorm.DB) (string, interface{}, error)
}

// RegisterTemplateParameterGenerator registers a template parameter generator used in the template logic
func RegisterTemplateParameterGenerator(templateParameterGenerator TemplateParameterGenerator) {
	templateParameterGenerators = append(templateParameterGenerators, templateParameterGenerator)
}

// RegisteredTemplateParameterGenerators returns the registered template parameter generators
func RegisteredTemplateParameterGenerators() []TemplateParameterGenerator {
	result := []TemplateParameterGenerator{}
	result = append(result, templateParameterGenerators...)
	return result
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
