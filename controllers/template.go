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
			"template",
			models.SharedTemplateModel(),
			logics.UniqueTemplateLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *templateController) RouteMap() map[int]map[string]gin.HandlerFunc {
	resourceSingleUrl := extensions.GetResourceSingleUrl(controller.ResourceName())
	resourceMultiUrl := extensions.GetResourceMultiUrl(controller.ResourceName())

	routeMap := map[int]map[string]gin.HandlerFunc{
		extensions.MethodGet: {
			resourceSingleUrl: controller.GetSingle,
			resourceMultiUrl:  controller.GetMulti,
		},
		extensions.MethodPost: {
			resourceMultiUrl: controller.Create,
		},
		extensions.MethodPut: {
			resourceSingleUrl: controller.Update,
		},
		extensions.MethodDelete: {
			resourceSingleUrl: controller.Delete,
		},
		extensions.MethodPatch: {
			resourceSingleUrl: controller.Patch,
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
