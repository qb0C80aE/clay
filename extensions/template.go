package extensions

import (
	"github.com/jinzhu/gorm"
	"text/template"
)

var templateParameterGenerators = []TemplateParameterGenerator{}
var templateFuncMaps = []template.FuncMap{}

type TemplateParameterGenerator interface {
	GenerateTemplateParameter(db *gorm.DB) (string, interface{}, error)
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
