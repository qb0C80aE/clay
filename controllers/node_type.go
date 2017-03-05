package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type NodeTypeController struct {
	BaseController
}

func init() {
	extension.RegisterController(NewNodeTypeController())
}

func NewNodeTypeController() *NodeTypeController {
	controller := &NodeTypeController{}
	controller.Initialize()
	return controller
}

func (this *NodeTypeController) Initialize() {
	this.ResourceName = "node_type"
	this.Model = models.NodeTypeModel
	this.Logic = logics.NewNodeTypeLogic()
	this.Outputter = this
}

func (this *NodeTypeController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
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
