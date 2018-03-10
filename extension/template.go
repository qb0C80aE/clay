package extension

import (
	"text/template"
)

var templateFuncMapList = []template.FuncMap{}

// RegisterTemplateFuncMap registers a template FuncMap used in the template logic
func RegisterTemplateFuncMap(templateFuncMap template.FuncMap) {
	templateFuncMapList = append(templateFuncMapList, templateFuncMap)
}

// GetRegisteredTemplateFuncMapList returns the registered template FuncMap
func GetRegisteredTemplateFuncMapList() []template.FuncMap {
	result := []template.FuncMap{}
	result = append(result, templateFuncMapList...)
	return result
}
