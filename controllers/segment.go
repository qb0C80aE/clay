package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

func init() {
	resourceName := "segment"
	extension.RegisterEndpoint(resourceName)

	resourceMultiUrl := extension.GetResourceMultiUrl(resourceName)
	extension.RegisterRoute(extension.MethodGet, resourceMultiUrl, GetSegments)
}

func GetSegments(c *gin.Context) {
	ProcessMultiGet(c, models.SegmentModel, logics.GetSegments, OutputJsonError, OutputMultiJsonResult)
}
