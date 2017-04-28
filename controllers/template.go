package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type templateController struct {
	*BaseController
}

func newTemplateController() extensions.Controller {
	controller := &templateController{
		BaseController: NewBaseController(
			models.SharedTemplateModel(),
			logics.UniqueTemplateLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *templateController) RouteMap() map[int]map[string]gin.HandlerFunc {
	resourceSingleURL := controller.ResourceSingleURL()
	resourceMultiURL := controller.ResourceMultiURL()

	routeMap := map[int]map[string]gin.HandlerFunc{
		extensions.MethodGet: {
			resourceSingleURL: controller.GetSingle,
			resourceMultiURL:  controller.GetMulti,
		},
		extensions.MethodPost: {
			resourceMultiURL: controller.Create,
		},
		extensions.MethodPut: {
			resourceSingleURL: controller.Update,
		},
		extensions.MethodDelete: {
			resourceSingleURL: controller.Delete,
		},
		extensions.MethodPatch: {
			resourceSingleURL: controller.Patch,
		},
	}
	return routeMap
}

func (controller *templateController) OutputPatch(c *gin.Context, code int, result interface{}) {
	text := result.(string)
	c.String(code, text)
}

var uniqueTemplateController = newTemplateController()

func init() {
	extensions.RegisterController(uniqueTemplateController)
}
