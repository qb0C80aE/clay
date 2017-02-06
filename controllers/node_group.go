package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetNodeGroups(c *gin.Context) {
	processMultiGet(c, models.NodeGroupModel, logics.GetNodeGroups, OutputJsonError, OutputMultiJsonResult)
}

func GetNodeGroup(c *gin.Context) {
	processSingleGet(c, models.NodeGroupModel, logics.GetNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func CreateNodeGroup(c *gin.Context) {
	container := &models.NodeGroup{}
	processCreate(c, container, BindJson, logics.CreateNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func UpdateNodeGroup(c *gin.Context) {
	container := &models.NodeGroup{}
	processUpdate(c, container, BindJson, logics.UpdateNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func DeleteNodeGroup(c *gin.Context) {
	processDelete(c, logics.DeleteNodeGroup, OutputJsonError)
}
