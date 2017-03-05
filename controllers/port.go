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
	this.ResourceName = "port"
	this.Model = models.PortModel
	this.Logic = logics.NewPortLogic()
	this.Outputter = this
}

func (this *PortController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
	resourceSingleUrl := extension.GetResourceSingleUrl(this.ResourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(this.ResourceName)

	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			resourceSingleUrl: this.GetSingle,
			resourceMultiUrl:  this.GetMulti,
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
