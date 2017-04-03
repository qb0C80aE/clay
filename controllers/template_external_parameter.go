package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type templateExternalParameterController struct {
	*BaseController
}

func newTemplateExternalParameterController() extension.Controller {
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
	resourceSingleUrl := extension.GetResourceSingleUrl(controller.ResourceName())
	resourceMultiUrl := extension.GetResourceMultiUrl(controller.ResourceName())

	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			resourceSingleUrl: controller.GetSingle,
			resourceMultiUrl:  controller.GetMulti,
		},
		extension.MethodPost: {
			resourceMultiUrl: controller.Create,
		},
		extension.MethodPut: {
			resourceSingleUrl: controller.Update,
		},
		extension.MethodDelete: {
			resourceSingleUrl: controller.Delete,
		},
	}
	return routeMap
}

var uniqueTemplateExternalParameterController = newTemplateExternalParameterController()

func init() {
	extension.RegisterController(uniqueTemplateExternalParameterController)
}
