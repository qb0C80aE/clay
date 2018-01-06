package extensions

import (
	"text/template"
)

var templateFuncMaps = []template.FuncMap{}

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
