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

func (controller *templateController) RouteMap() map[int]map[int]gin.HandlerFunc {
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

var uniqueTemplateController = newTemplateController()

func init() {
	extensions.RegisterController(uniqueTemplateController)
}
