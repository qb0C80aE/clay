package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type designController struct {
	*BaseController
}

func newDesignController() extensions.Controller {
	controller := &designController{
		BaseController: NewBaseController(
			models.SharedDesignModel(),
			logics.UniqueDesignLogic(),
		),
	}
	return controller
}

func (controller *designController) ResourceSingleURL() string {
	return fmt.Sprintf("%s/present", controller.ResourceName())
}

func (controller *designController) RouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extensions.MethodGet: {
			extensions.URLSingle: controller.GetSingle,
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

var uniqueDesignController = newDesignController()

func init() {
	extensions.RegisterController(uniqueDesignController)
}
