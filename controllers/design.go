package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetDesign(c *gin.Context) {
	processSingleGet(c, models.DesignModel, logics.GetDesign, OutputJsonError, OutputSingleJsonResult)
}

func UpdateDesign(c *gin.Context) {
	container := &models.Design{}
	processUpdate(c, container, BindJson, logics.UpdateDesign, OutputJsonError, OutputSingleJsonResult)
}

func DeleteDesign(c *gin.Context) {
	processDelete(c, logics.DeleteDesign, OutputJsonError, OutputNothing)
}
