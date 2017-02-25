package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func init() {
	resourceName := "template_external_parameter"
	extension.RegisterEndpoint(resourceName)

	resourceSingleUrl := extension.GetResourceSingleUrl(resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(resourceName)
	extension.RegisterRoute(extension.MethodGet, resourceMultiUrl, GetTemplateExternalParameters)
	extension.RegisterRoute(extension.MethodGet, resourceSingleUrl, GetTemplateExternalParameter)
	extension.RegisterRoute(extension.MethodPost, resourceMultiUrl, CreateTemplateExternalParameter)
	extension.RegisterRoute(extension.MethodPut, resourceSingleUrl, UpdateTemplateExternalParameter)
	extension.RegisterRoute(extension.MethodDelete, resourceSingleUrl, DeleteTemplateExternalParameter)

	resourceName = "template"
	extension.RegisterEndpoint(resourceName)

	resourceSingleUrl = extension.GetResourceSingleUrl(resourceName)
	resourceMultiUrl = extension.GetResourceMultiUrl(resourceName)
	extension.RegisterRoute(extension.MethodGet, resourceMultiUrl, GetTemplates)
	extension.RegisterRoute(extension.MethodGet, resourceSingleUrl, GetTemplate)
	extension.RegisterRoute(extension.MethodPost, resourceMultiUrl, CreateTemplate)
	extension.RegisterRoute(extension.MethodPut, resourceSingleUrl, UpdateTemplate)
	extension.RegisterRoute(extension.MethodDelete, resourceSingleUrl, DeleteTemplate)
	extension.RegisterRoute(extension.MethodPatch, resourceSingleUrl, ApplyTemplate)
}

func GetTemplateExternalParameters(c *gin.Context) {
	ProcessMultiGet(c, models.TemplateExternalParameterModel, logics.GetTemplateExternalParameters, OutputJsonError, OutputMultiJsonResult)
}

func GetTemplateExternalParameter(c *gin.Context) {
	ProcessSingleGet(c, models.TemplateExternalParameterModel, logics.GetTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func CreateTemplateExternalParameter(c *gin.Context) {
	ProcessCreate(c, &models.TemplateExternalParameter{}, logics.CreateTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func UpdateTemplateExternalParameter(c *gin.Context) {
	ProcessUpdate(c, &models.TemplateExternalParameter{}, logics.UpdateTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func DeleteTemplateExternalParameter(c *gin.Context) {
	ProcessDelete(c, logics.DeleteTemplateExternalParameter, OutputJsonError, OutputNothing)
}

func GetTemplates(c *gin.Context) {
	ProcessMultiGet(c, models.TemplateModel, logics.GetTemplates, OutputJsonError, OutputMultiJsonResult)
}

func GetTemplate(c *gin.Context) {
	ProcessSingleGet(c, models.TemplateModel, logics.GetTemplate, OutputJsonError, OutputSingleJsonResult)
}

func CreateTemplate(c *gin.Context) {
	ProcessCreate(c, &models.Template{}, logics.CreateTemplate, OutputJsonError, OutputNothing)
}

func UpdateTemplate(c *gin.Context) {
	ProcessUpdate(c, &models.Template{}, logics.UpdateTemplate, OutputJsonError, OutputNothing)
}

func DeleteTemplate(c *gin.Context) {
	ProcessDelete(c, logics.DeleteTemplate, OutputJsonError, OutputNothing)
}

func ApplyTemplate(c *gin.Context) {
	ProcessSingleGet(c, &models.Template{}, logics.ApplyTemplate, OutputJsonError, OutputTextResult)
}
