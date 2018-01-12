package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type templateArgumentController struct {
	*BaseController
}

func newTemplateArgumentController() extensions.Controller {
	controller := &templateArgumentController{
		BaseController: NewBaseController(
			models.SharedTemplateArgumentModel(),
			logics.UniqueTemplateArgumentLogic(),
		),
	}
	return controller
}

func (controller *templateArgumentController) RouteMap() map[int]map[int]gin.HandlerFunc {
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

var uniqueTemplateArgumentController = newTemplateArgumentController()

func init() {
	extensions.RegisterController(uniqueTemplateArgumentController)
}
