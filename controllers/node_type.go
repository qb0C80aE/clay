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
	this.resourceName = "node_type"
}

func (this *NodeTypeController) GetResourceName() string {
	return this.resourceName
}

func (this *NodeTypeController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
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

func (_ *NodeTypeController) GetMulti(c *gin.Context) {
	ProcessMultiGet(c, models.NodeTypeModel, logics.GetNodeTypes, OutputJsonError, OutputMultiJsonResult)
}

func (_ *NodeTypeController) GetSingle(c *gin.Context) {
	ProcessSingleGet(c, models.NodeTypeModel, logics.GetNodeType, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeTypeController) Create(c *gin.Context) {
	ProcessCreate(c, &models.NodeType{}, logics.CreateNodeType, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeTypeController) Update(c *gin.Context) {
	ProcessUpdate(c, &models.NodeType{}, logics.UpdateNodeType, OutputJsonError, OutputSingleJsonResult)
}

func (_ *NodeTypeController) Delete(c *gin.Context) {
	ProcessDelete(c, logics.DeleteNodeType, OutputJsonError, OutputNothing)
}
