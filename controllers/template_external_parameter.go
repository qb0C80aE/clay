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
			"template_external_parameter",
			models.SharedTemplateExternalParameterModel(),
			logics.UniqueTemplateExternalParameterLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *templateExternalParameterController) RouteMap() map[int]map[string]gin.HandlerFunc {
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
	}
	return routeMap
}

var uniqueTemplateExternalParameterController = newTemplateExternalParameterController()

func init() {
	extensions.RegisterController(uniqueTemplateExternalParameterController)
}
