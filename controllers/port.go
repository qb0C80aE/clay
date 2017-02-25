package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func init() {
	resourceName := "port"
	extension.RegisterEndpoint(resourceName)

	resourceSingleUrl := extension.GetResourceSingleUrl(resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(resourceName)
	extension.RegisterRoute(extension.MethodGet, resourceMultiUrl, GetPorts)
	extension.RegisterRoute(extension.MethodGet, resourceSingleUrl, GetPort)
	extension.RegisterRoute(extension.MethodPost, resourceMultiUrl, CreatePort)
	extension.RegisterRoute(extension.MethodPut, resourceSingleUrl, UpdatePort)
	extension.RegisterRoute(extension.MethodDelete, resourceSingleUrl, DeletePort)
}

func GetPorts(c *gin.Context) {
	ProcessMultiGet(c, models.PortModel, logics.GetPorts, OutputJsonError, OutputMultiJsonResult)
}

func GetPort(c *gin.Context) {
	ProcessSingleGet(c, models.PortModel, logics.GetPort, OutputJsonError, OutputSingleJsonResult)
}

func CreatePort(c *gin.Context) {
	ProcessCreate(c, &models.Port{}, logics.CreatePort, OutputJsonError, OutputSingleJsonResult)
}

func UpdatePort(c *gin.Context) {
	ProcessUpdate(c, &models.Port{}, logics.UpdatePort, OutputJsonError, OutputSingleJsonResult)
}

func DeletePort(c *gin.Context) {
	ProcessDelete(c, logics.DeletePort, OutputJsonError, OutputNothing)
}
