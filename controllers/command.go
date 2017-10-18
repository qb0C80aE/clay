package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type commandController struct {
	*BaseController
}

func newCommandController() extensions.Controller {
	controller := &commandController{
		BaseController: NewBaseController(
			models.SharedCommandModel(),
			logics.UniqueCommandLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *commandController) RouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extensions.MethodGet: {
			extensions.URLSingle: controller.GetSingle,
			extensions.URLMulti:  controller.GetMulti,
		},
		extensions.MethodPost: {
			extensions.URLMulti: controller.Create,
		},
		extensions.MethodDelete: {
			extensions.URLSingle: controller.Delete,
		},
	}
	return routeMap
}

var uniqueCommandController = newCommandController()

func init() {
	extensions.RegisterController(uniqueCommandController)
}
