package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func init() {
	resourceName := "node"
	extension.RegisterEndpoint(resourceName)

	resourceSingleUrl := extension.GetResourceSingleUrl(resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(resourceName)
	extension.RegisterRoute(extension.MethodGet, resourceMultiUrl, GetNodes)
	extension.RegisterRoute(extension.MethodGet, resourceSingleUrl, GetNode)
	extension.RegisterRoute(extension.MethodPost, resourceMultiUrl, CreateNode)
	extension.RegisterRoute(extension.MethodPut, resourceSingleUrl, UpdateNode)
	extension.RegisterRoute(extension.MethodDelete, resourceSingleUrl, DeleteNode)
}

func GetNodes(c *gin.Context) {
	ProcessMultiGet(c, models.NodeModel, logics.GetNodes, OutputJsonError, OutputMultiJsonResult)
}

func GetNode(c *gin.Context) {
	ProcessSingleGet(c, models.NodeModel, logics.GetNode, OutputJsonError, OutputSingleJsonResult)
}

func CreateNode(c *gin.Context) {
	ProcessCreate(c, &models.Node{}, logics.CreateNode, OutputJsonError, OutputSingleJsonResult)
}

func UpdateNode(c *gin.Context) {
	ProcessUpdate(c, &models.Node{}, logics.UpdateNode, OutputJsonError, OutputSingleJsonResult)
}

func DeleteNode(c *gin.Context) {
	ProcessDelete(c, logics.DeleteNode, OutputJsonError, OutputNothing)
}
