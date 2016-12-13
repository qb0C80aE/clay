package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func GetSegments(c *gin.Context) {
	processMultiGet(c, models.SegmentModel, logics.GetSegments)
}
