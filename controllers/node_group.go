package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetNodeGroups(c *gin.Context) {
	processMultiGet(c, models.NodeGroupModel, logics.GetNodeGroups)
}

func GetNodeGroup(c *gin.Context) {
	processSingleGet(c, models.NodeGroupModel, logics.GetNodeGroup)
}

func CreateNodeGroup(c *gin.Context) {
	container := &models.NodeGroup{}
	processCreate(c, container, models.NodeGroupModel, logics.CreateNodeGroup)
}

func UpdateNodeGroup(c *gin.Context) {
	container := &models.NodeGroup{}
	processUpdate(c, container, models.NodeGroupModel, logics.UpdateNodeGroup)
}

func DeleteNodeGroup(c *gin.Context) {
	processDelete(c, models.NodeGroupModel, logics.DeleteNodeGroup)
}
