package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetNodes(c *gin.Context) {
	processMultiGet(c, models.NodeModel, logics.GetNodes, OutputJsonError, OutputMultiJsonResult)
}

func GetNode(c *gin.Context) {
	processSingleGet(c, models.NodeModel, logics.GetNode, OutputJsonError, OutputSingleJsonResult)
}

func CreateNode(c *gin.Context) {
	container := &models.Node{}
	processCreate(c, container, logics.CreateNode, OutputJsonError, OutputSingleJsonResult)
}

func UpdateNode(c *gin.Context) {
	container := &models.Node{}
	processUpdate(c, container, logics.UpdateNode, OutputJsonError, OutputSingleJsonResult)
}

func DeleteNode(c *gin.Context) {
	processDelete(c, logics.DeleteNode, OutputJsonError, OutputNothing)
}
