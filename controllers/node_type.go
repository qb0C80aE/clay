package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetNodeTypes(c *gin.Context) {
	processMultiGet(c, models.NodeModel, logics.GetNodeTypes)
}

func GetNodeType(c *gin.Context) {
	processSingleGet(c, models.NodeModel, logics.GetNodeType)
}

func CreateNodeType(c *gin.Context) {
	container := &models.Node{}
	processCreate(c, container, models.NodeModel, logics.CreateNodeType)
}

func UpdateNodeType(c *gin.Context) {
	container := &models.Node{}
	processUpdate(c, container, models.NodeModel, logics.UpdateNodeType)
}

func DeleteNodeType(c *gin.Context) {
	processDelete(c, models.NodeModel, logics.DeleteNodeType)
}
