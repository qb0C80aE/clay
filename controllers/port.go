package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type PortController struct {
	BaseController
}

func init() {
	extension.RegisterController(NewPortController())
}

func NewPortController() *PortController {
	controller := &PortController{}
	controller.Initialize()
	return controller
}

func (this *PortController) Initialize() {
	this.resourceName = "port"
}

func (this *PortController) GetResourceName() string {
	return this.resourceName
}

func (this *PortController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
	resourceSingleUrl := extension.GetResourceSingleUrl(this.resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(this.resourceName)

	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			resourceMultiUrl:  this.GetSingle,
			resourceSingleUrl: this.GetMulti,
		},
		extension.MethodPost: {
			resourceMultiUrl: this.Create,
		},
		extension.MethodPut: {
			resourceSingleUrl: this.Update,
		},
		extension.MethodDelete: {
			resourceSingleUrl: this.Delete,
		},
	}
	return routeMap
}

func (_ *PortController) GetMulti(c *gin.Context) {
	ProcessMultiGet(c, models.PortModel, logics.GetPorts, OutputJsonError, OutputMultiJsonResult)
}

func (_ *PortController) GetSingle(c *gin.Context) {
	ProcessSingleGet(c, models.PortModel, logics.GetPort, OutputJsonError, OutputSingleJsonResult)
}

func (_ *PortController) Create(c *gin.Context) {
	ProcessCreate(c, &models.Port{}, logics.CreatePort, OutputJsonError, OutputSingleJsonResult)
}

func (_ *PortController) Update(c *gin.Context) {
	ProcessUpdate(c, &models.Port{}, logics.UpdatePort, OutputJsonError, OutputSingleJsonResult)
}

func (_ *PortController) Delete(c *gin.Context) {
	ProcessDelete(c, logics.DeletePort, OutputJsonError, OutputNothing)
}
