package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type NodeGroupController struct {
	BaseController
}

func init() {
	extension.RegisterController(NewNodeGroupController())
}

func NewNodeGroupController() *NodeGroupController {
	controller := &NodeGroupController{}
	controller.Initialize()
	return controller
}

func (this *NodeGroupController) Initialize() {
	this.ResourceName = "node_group"
	this.Model = models.NodeGroupModel
	this.Logic = logics.NewNodeGroupLogic()
	this.Outputter = this
}

func (this *NodeGroupController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
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
