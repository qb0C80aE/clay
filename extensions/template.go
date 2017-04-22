package extensions

import (
	"fmt"
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

func init() {
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { return a / b },
		"mod": func(a, b int) int { return a % b },
		"slice": func(items ...interface{}) []interface{} {
			slice := []interface{}{}
			return append(slice, items...)
		},
		"subslice": func(slice []interface{}, begin int, end int) []interface{} {
			return slice[begin:end]
		},
		"append": func(slice []interface{}, item ...interface{}) []interface{} {
			return append(slice, item...)
		},
		"concatenate": func(slice1 []interface{}, slice2 []interface{}) []interface{} {
			return append(slice1, slice2...)
		},
		"map": func(pairs ...interface{}) (map[interface{}]interface{}, error) {
			if len(pairs)%2 == 1 {
				return nil, fmt.Errorf("Number of arguments must be even")
			}
			m := make(map[interface{}]interface{}, len(pairs)/2)
			for i := 0; i < len(pairs); i += 2 {
				m[pairs[i]] = pairs[i+1]
			}
			return m, nil
		},
		"put": func(target map[interface{}]interface{}, key interface{}, value interface{}) map[interface{}]interface{} {
			target[key] = value
			return target
		},
		"get": func(target map[interface{}]interface{}, key interface{}) interface{} {
			return target[key]
		},
		"delete": func(target map[interface{}]interface{}, key interface{}) map[interface{}]interface{} {
			delete(target, key)
			return target
		},
		"merge": func(source, destination map[interface{}]interface{}) map[interface{}]interface{} {
			for key, value := range source {
				destination[key] = value
			}
			return destination
		},
	}
	RegisterTemplateFuncMap(funcMap)
}
