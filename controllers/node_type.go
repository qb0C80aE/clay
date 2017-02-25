package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func init() {
	resourceName := "node_type"
	extension.RegisterEndpoint(resourceName)

	resourceSingleUrl := extension.GetResourceSingleUrl(resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(resourceName)
	extension.RegisterRoute(extension.MethodGet, resourceMultiUrl, GetNodeTypes)
	extension.RegisterRoute(extension.MethodGet, resourceSingleUrl, GetNodeType)
	extension.RegisterRoute(extension.MethodPost, resourceMultiUrl, CreateNodeType)
	extension.RegisterRoute(extension.MethodPut, resourceSingleUrl, UpdateNodeType)
	extension.RegisterRoute(extension.MethodDelete, resourceSingleUrl, DeleteNodeType)
}

func GetNodeTypes(c *gin.Context) {
	ProcessMultiGet(c, models.NodeTypeModel, logics.GetNodeTypes, OutputJsonError, OutputMultiJsonResult)
}

func GetNodeType(c *gin.Context) {
	ProcessSingleGet(c, models.NodeTypeModel, logics.GetNodeType, OutputJsonError, OutputSingleJsonResult)
}

func CreateNodeType(c *gin.Context) {
	ProcessCreate(c, &models.NodeType{}, logics.CreateNodeType, OutputJsonError, OutputSingleJsonResult)
}

func UpdateNodeType(c *gin.Context) {
	ProcessUpdate(c, &models.NodeType{}, logics.UpdateNodeType, OutputJsonError, OutputSingleJsonResult)
}

func DeleteNodeType(c *gin.Context) {
	ProcessDelete(c, logics.DeleteNodeType, OutputJsonError, OutputNothing)
}
