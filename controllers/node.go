package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type NodeController struct {
	BaseController
}

func init() {
	extension.RegisterController(NewNodeController())
}

func NewNodeController() *NodeController {
	controller := &NodeController{}
	controller.Initialize()
	return controller
}

func (this *NodeController) Initialize() {
	this.resourceName = "node"
}

func (this *NodeController) GetResourceName() string {
	return this.resourceName
}

func (this *NodeController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
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

func (_ *NodeController) GetSingle(c *gin.Context) {
	ProcessMultiGet(c, models.NodeModel, logics.GetNodes, OutputJsonError, OutputMultiJsonResult)
}

func (_ *NodeController) GetMulti(c *gin.Context) {
	ProcessSingleGet(c, models.NodeModel, logics.GetNode, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeController) Create(c *gin.Context) {
	ProcessCreate(c, &models.Node{}, logics.CreateNode, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeController) Update(c *gin.Context) {
	ProcessUpdate(c, &models.Node{}, logics.UpdateNode, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeController) Delete(c *gin.Context) {
	ProcessDelete(c, logics.DeleteNode, OutputJsonError, OutputNothing)
}
