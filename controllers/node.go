package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetNodes(c *gin.Context) {
	processMultiGet(c, models.NodeModel, logics.GetNodes)
}

func GetNode(c *gin.Context) {
	processSingleGet(c, models.NodeModel, logics.GetNode)
}

func CreateNode(c *gin.Context) {
	container := &models.Node{}
	processCreate(c, container, models.NodeModel, logics.CreateNode)
}

func UpdateNode(c *gin.Context) {
	container := &models.Node{}
	processUpdate(c, container, models.NodeModel, logics.UpdateNode)
}

func DeleteNode(c *gin.Context) {
	processDelete(c, models.NodeModel, logics.DeleteNode)
}
