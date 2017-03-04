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
	this.resourceName = "node_group"
}

func (this *NodeGroupController) GetResourceName() string {
	return this.resourceName
}

func (this *NodeGroupController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
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

func (_ *NodeGroupController) GetMulti(c *gin.Context) {
	ProcessMultiGet(c, models.NodeGroupModel, logics.GetNodeGroups, OutputJsonError, OutputMultiJsonResult)
}

func (_ *NodeGroupController) GetSingle(c *gin.Context) {
	ProcessSingleGet(c, models.NodeGroupModel, logics.GetNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeGroupController) Create(c *gin.Context) {
	ProcessCreate(c, &models.NodeGroup{}, logics.CreateNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeGroupController) Update(c *gin.Context) {
	ProcessUpdate(c, &models.NodeGroup{}, logics.UpdateNodeGroup, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeGroupController) Delete(c *gin.Context) {
	ProcessDelete(c, logics.DeleteNodeGroup, OutputJsonError, OutputNothing)
}
