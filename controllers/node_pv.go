package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetNodePvs(c *gin.Context) {
	processMultiGet(c, models.NodeModel, logics.GetNodePvs)
}

func GetNodePv(c *gin.Context) {
	processSingleGet(c, models.NodeModel, logics.GetNodePv)
}

func CreateNodePv(c *gin.Context) {
	container := &models.Node{}
	processCreate(c, container, models.NodeModel, logics.CreateNodePv)
}

func UpdateNodePv(c *gin.Context) {
	container := &models.Node{}
	processUpdate(c, container, models.NodeModel, logics.UpdateNodePv)
}

func DeleteNodePv(c *gin.Context) {
	processDelete(c, models.NodeModel, logics.DeleteNodePv)
}
