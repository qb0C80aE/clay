package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type NodePvController struct {
	BaseController
}

func init() {
	extension.RegisterController(NewNodePvController())
}

func NewNodePvController() *NodePvController {
	controller := &NodePvController{}
	controller.Initialize()
	return controller
}

func (this *NodePvController) Initialize() {
	this.resourceName = "node_pv"
}

func (this *NodePvController) GetResourceName() string {
	return this.resourceName
}

func (this *NodePvController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
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

func (_ *NodePvController) GetMulti(c *gin.Context) {
	ProcessMultiGet(c, models.NodePvModel, logics.GetNodePvs, OutputJsonError, OutputMultiJsonResult)
}

func (_ *NodePvController) GetSingle(c *gin.Context) {
	ProcessSingleGet(c, models.NodePvModel, logics.GetNodePv, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodePvController) Create(c *gin.Context) {
	ProcessCreate(c, &models.NodePv{}, logics.CreateNodePv, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodePvController) Update(c *gin.Context) {
	ProcessUpdate(c, &models.NodePv{}, logics.UpdateNodePv, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodePvController) Delete(c *gin.Context) {
	ProcessDelete(c, logics.DeleteNodePv, OutputJsonError, OutputNothing)
}
