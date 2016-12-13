package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetPorts(c *gin.Context) {
	processMultiGet(c, models.PortModel, logics.GetPorts)
}

func GetPort(c *gin.Context) {
	processSingleGet(c, models.PortModel, logics.GetPort)
}

func CreatePort(c *gin.Context) {
	container := &models.Port{}
	processCreate(c, container, models.PortModel, logics.CreatePort)
}

func UpdatePort(c *gin.Context) {
	container := &models.Port{}
	processUpdate(c, container, models.PortModel, logics.UpdatePort)
}

func DeletePort(c *gin.Context) {
	processDelete(c, models.PortModel, logics.DeletePort)
}
