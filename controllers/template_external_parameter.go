package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type templateExternalParameterController struct {
	*BaseController
}

func newTemplateExternalParameterController() extensions.Controller {
	controller := &templateExternalParameterController{
		BaseController: NewBaseController(
			models.SharedTemplateExternalParameterModel(),
			logics.UniqueTemplateExternalParameterLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *templateExternalParameterController) RouteMap() map[int]map[string]gin.HandlerFunc {
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
	}
	return routeMap
}

var uniqueTemplateExternalParameterController = newTemplateExternalParameterController()

func init() {
	extensions.RegisterController(uniqueTemplateExternalParameterController)
}
