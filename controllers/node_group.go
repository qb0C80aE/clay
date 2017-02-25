package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func init() {
	resourceName := "node_group"
	extension.RegisterEndpoint(resourceName)

	resourceSingleUrl := extension.GetResourceSingleUrl(resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(resourceName)
	extension.RegisterRoute(extension.MethodGet, resourceMultiUrl, GetNodeGroups)
	extension.RegisterRoute(extension.MethodGet, resourceSingleUrl, GetNodeGroup)
	extension.RegisterRoute(extension.MethodPost, resourceMultiUrl, CreateNodeGroup)
	extension.RegisterRoute(extension.MethodPut, resourceSingleUrl, UpdateNodeGroup)
	extension.RegisterRoute(extension.MethodDelete, resourceSingleUrl, DeleteNodeGroup)
}

func GetNodeGroups(c *gin.Context) {
	ProcessMultiGet(c, models.NodeGroupModel, logics.GetNodeGroups, OutputJsonError, OutputMultiJsonResult)
}

func GetNodeGroup(c *gin.Context) {
	ProcessSingleGet(c, models.NodeGroupModel, logics.GetNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func CreateNodeGroup(c *gin.Context) {
	ProcessCreate(c, &models.NodeGroup{}, logics.CreateNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func UpdateNodeGroup(c *gin.Context) {
	ProcessUpdate(c, &models.NodeGroup{}, logics.UpdateNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func DeleteNodeGroup(c *gin.Context) {
	ProcessDelete(c, logics.DeleteNodeGroup, OutputJsonError, OutputNothing)
}
