package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetNodeTypes(c *gin.Context) {
	processMultiGet(c, models.NodeModel, logics.GetNodeTypes, OutputJsonError, OutputMultiJsonResult)
}

func GetNodeType(c *gin.Context) {
	processSingleGet(c, models.NodeModel, logics.GetNodeType, OutputJsonError, OutputSingleJsonResult)
}

func CreateNodeType(c *gin.Context) {
	container := &models.Node{}
	processCreate(c, container, logics.CreateNodeType, OutputJsonError, OutputSingleJsonResult)
}

func UpdateNodeType(c *gin.Context) {
	container := &models.Node{}
	processUpdate(c, container, models.NodeModel, logics.UpdateNodeType, OutputJsonError, OutputSingleJsonResult)
}

func DeleteNodeType(c *gin.Context) {
	processDelete(c, logics.DeleteNodeType, OutputJsonError)
}
