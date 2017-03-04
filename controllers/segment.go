package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type SegmentController struct {
	BaseController
}

func init() {
	extension.RegisterController(NewSegmentController())
}

func NewSegmentController() *SegmentController {
	controller := &SegmentController{}
	controller.Initialize()
	return controller
}

func (this *SegmentController) Initialize() {
	this.resourceName = "segment"
}

func (this *SegmentController) GetResourceName() string {
	return this.resourceName
}

func (this *SegmentController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
	resourceMultiUrl := extension.GetResourceMultiUrl(this.resourceName)
	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			resourceMultiUrl: this.GetMulti,
		},
	}
	return routeMap
}

func (_ *SegmentController) GetMulti(c *gin.Context) {
	ProcessMultiGet(c, models.SegmentModel, logics.GetSegments, OutputJsonError, OutputMultiJsonResult)
}
