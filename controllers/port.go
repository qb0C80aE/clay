package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetPorts(c *gin.Context) {
	processMultiGet(c, models.PortModel, logics.GetPorts, OutputJsonError, OutputMultiJsonResult)
}

func GetPort(c *gin.Context) {
	processSingleGet(c, models.PortModel, logics.GetPort, OutputJsonError, OutputSingleJsonResult)
}

func CreatePort(c *gin.Context) {
	container := &models.Port{}
	processCreate(c, container, logics.CreatePort, OutputJsonError, OutputSingleJsonResult)
}

func UpdatePort(c *gin.Context) {
	container := &models.Port{}
	processUpdate(c, container, logics.UpdatePort, OutputJsonError, OutputSingleJsonResult)
}

func DeletePort(c *gin.Context) {
	processDelete(c, logics.DeletePort, OutputJsonError)
}
