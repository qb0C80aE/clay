package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetNodePvs(c *gin.Context) {
	processMultiGet(c, models.NodeModel, logics.GetNodePvs, OutputJsonError, OutputMultiJsonResult)
}

func GetNodePv(c *gin.Context) {
	processSingleGet(c, models.NodeModel, logics.GetNodePv, OutputJsonError, OutputSingleJsonResult)
}

func CreateNodePv(c *gin.Context) {
	container := &models.Node{}
	processCreate(c, container, logics.CreateNodePv, OutputJsonError, OutputSingleJsonResult)
}

func UpdateNodePv(c *gin.Context) {
	container := &models.Node{}
	processUpdate(c, container, logics.UpdateNodePv, OutputJsonError, OutputSingleJsonResult)
}

func DeleteNodePv(c *gin.Context) {
	processDelete(c, logics.DeleteNodePv, OutputJsonError, OutputNothing)
}
