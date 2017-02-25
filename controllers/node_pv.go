package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func init() {
	resourceName := "node_pv"
	extension.RegisterEndpoint(resourceName)

	resourceSingleUrl := extension.GetResourceSingleUrl(resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(resourceName)
	extension.RegisterRoute(extension.MethodGet, resourceMultiUrl, GetNodePvs)
	extension.RegisterRoute(extension.MethodGet, resourceSingleUrl, GetNodePv)
	extension.RegisterRoute(extension.MethodPost, resourceMultiUrl, CreateNodePv)
	extension.RegisterRoute(extension.MethodPut, resourceSingleUrl, UpdateNodePv)
	extension.RegisterRoute(extension.MethodDelete, resourceSingleUrl, DeleteNodePv)
}

func GetNodePvs(c *gin.Context) {
	ProcessMultiGet(c, models.NodePvModel, logics.GetNodePvs, OutputJsonError, OutputMultiJsonResult)
}

func GetNodePv(c *gin.Context) {
	ProcessSingleGet(c, models.NodePvModel, logics.GetNodePv, OutputJsonError, OutputSingleJsonResult)
}

func CreateNodePv(c *gin.Context) {
	ProcessCreate(c, &models.NodePv{}, logics.CreateNodePv, OutputJsonError, OutputSingleJsonResult)
}

func UpdateNodePv(c *gin.Context) {
	ProcessUpdate(c, &models.NodePv{}, logics.UpdateNodePv, OutputJsonError, OutputSingleJsonResult)
}

func DeleteNodePv(c *gin.Context) {
	ProcessDelete(c, logics.DeleteNodePv, OutputJsonError, OutputNothing)
}
