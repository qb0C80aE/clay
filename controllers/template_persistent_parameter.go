package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type templatePersistentParameterController struct {
	*BaseController
}

func newTemplatePersistentParameterController() extensions.Controller {
	controller := &templatePersistentParameterController{
		BaseController: NewBaseController(
			models.SharedTemplatePersistentParameterModel(),
			logics.UniqueTemplatePersistentParameterLogic(),
		),
	}
	return controller
}

func (controller *templatePersistentParameterController) RouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extensions.MethodGet: {
			extensions.URLSingle: controller.GetSingle,
			extensions.URLMulti:  controller.GetMulti,
		},
		extensions.MethodPost: {
			extensions.URLMulti: controller.Create,
		},
		extensions.MethodPut: {
			extensions.URLSingle: controller.Update,
		},
		extensions.MethodDelete: {
			extensions.URLSingle: controller.Delete,
		},
	}
	return routeMap
}

var uniqueTemplatePersistentParameterController = newTemplatePersistentParameterController()

func init() {
	extensions.RegisterController(uniqueTemplatePersistentParameterController)
}
