package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetPhysicalDiagram(c *gin.Context) {
	processSingleGet(c, models.DiagramModel, logics.GetPhysicalDiagram)
}

func GetLogicalDiagram(c *gin.Context) {
	processSingleGet(c, models.DiagramModel, logics.GetLogicalDiagram)

}
