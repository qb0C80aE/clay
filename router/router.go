package router

import (
	"github.com/qb0C80aE/clay/controllers"

	"github.com/gin-gonic/gin"
	"net/http"
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

		api.GET("/template_external_parameters", controllers.GetTemplateExternalParameters)
		api.GET("/template_external_parameters/:id", controllers.GetTemplateExternalParameter)
		api.POST("/template_external_parameters", controllers.CreateTemplateExternalParameter)
		api.PUT("/template_external_parameters/:id", controllers.UpdateTemplateExternalParameter)
		api.DELETE("/template_external_parameters/:id", controllers.DeleteTemplateExternalParameter)

		api.GET("/templates", controllers.GetTemplates)
		api.GET("/templates/:id", controllers.GetTemplate)
		api.POST("/templates", controllers.CreateTemplate)
		api.PUT("/templates/:id", controllers.UpdateTemplate)
		api.DELETE("/templates/:id", controllers.DeleteTemplate)
		api.PATCH("/templates/:id", controllers.ApplyTemplate)

	}

	r.Static("ui/files", "ui/files")
	r.LoadHTMLGlob("ui/templates/*.tmpl")
	ui := r.Group("/ui")
	ui.GET("/", func(c *gin.Context) { c.HTML(http.StatusOK, "index.tmpl", gin.H{}) })
	ui.GET("/network", func(c *gin.Context) { c.HTML(http.StatusOK, "network.tmpl", gin.H{}) })
	ui.GET("/diagram", func(c *gin.Context) { c.HTML(http.StatusOK, "diagram.tmpl", gin.H{}) })

}
