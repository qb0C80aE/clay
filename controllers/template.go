package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/models"
	"github.com/qb0C80aE/clay/logics"
)

func GetTemplateExternalParameters(c *gin.Context) {
	processMultiGet(c, models.TemplateExternalParameterModel, logics.GetTemplateExternalParameters, OutputJsonError, OutputMultiJsonResult)
}

func GetTemplateExternalParameter(c *gin.Context) {
	processSingleGet(c, models.TemplateExternalParameterModel, logics.GetTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func CreateTemplateExternalParameter(c *gin.Context) {
	container := &models.TemplateExternalParameter{}
	processCreate(c, container, logics.CreateTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func UpdateTemplateExternalParameter(c *gin.Context) {
	container := &models.TemplateExternalParameter{}
	processUpdate(c, container, logics.UpdateTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func DeleteTemplateExternalParameter(c *gin.Context) {
	processDelete(c, logics.DeleteTemplateExternalParameter, OutputJsonError, OutputNothing)
}

func CreateTemplate(c *gin.Context) {
	container := &models.Template{}
	processCreate(c, container, logics.CreateTemplate, OutputJsonError, OutputNothing)
}

func UpdateTemplate(c *gin.Context) {
	container := &models.Template{}
	processUpdate(c, container, logics.UpdateTemplate, OutputJsonError, OutputNothing)
}

func GetTemplate(c *gin.Context) {
	processSingleGet(c, models.TemplateModel, logics.GetTemplate, OutputJsonError, OutputTextResult)
}
