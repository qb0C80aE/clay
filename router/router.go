package router

import (
	"github.com/qb0C80aE/clay/controllers"

	"github.com/gin-gonic/gin"
)

func Initialize(r *gin.Engine) {

	r.GET("/", controllers.APIEndpoints)

	api := r.Group("/v1")
	{

		api.GET("/nodes", controllers.GetNodes)
		api.GET("/nodes/:id", controllers.GetNode)
		api.POST("/nodes", controllers.CreateNode)
		api.PUT("/nodes/:id", controllers.UpdateNode)
		api.DELETE("/nodes/:id", controllers.DeleteNode)

		api.GET("/node_groups", controllers.GetNodeGroups)
		api.GET("/node_groups/:id", controllers.GetNodeGroup)
		api.POST("/node_groups", controllers.CreateNodeGroup)
		api.PUT("/node_groups/:id", controllers.UpdateNodeGroup)
		api.DELETE("/node_groups/:id", controllers.DeleteNodeGroup)

		api.GET("/node_pvs", controllers.GetNodePvs)
		api.GET("/node_pvs/:id", controllers.GetNodePv)
		api.POST("/node_pvs", controllers.CreateNodePv)
		api.PUT("/node_pvs/:id", controllers.UpdateNodePv)
		api.DELETE("/node_pvs/:id", controllers.DeleteNodePv)

		api.GET("/node_types", controllers.GetNodeTypes)
		api.GET("/node_types/:id", controllers.GetNodeType)
		api.POST("/node_types", controllers.CreateNodeType)
		api.PUT("/node_types/:id", controllers.UpdateNodeType)
		api.DELETE("/node_types/:id", controllers.DeleteNodeType)

		api.GET("/ports", controllers.GetPorts)
		api.GET("/ports/:id", controllers.GetPort)
		api.POST("/ports", controllers.CreatePort)
		api.PUT("/ports/:id", controllers.UpdatePort)
		api.DELETE("/ports/:id", controllers.DeletePort)

		api.GET("/designs/present", controllers.GetDesign)
		api.PUT("/designs/present", controllers.UpdateDesign)
		api.DELETE("/designs/present", controllers.DeleteDesign)

		api.GET("/diagrams/physical", controllers.GetPhysicalDiagram)
		api.GET("/diagrams/logical", controllers.GetLogicalDiagram)

		api.GET("/segments", controllers.GetSegments)

	}

}
