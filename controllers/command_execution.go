package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type commandExecutionController struct {
	*BaseController
}

func newCommandExecutionController() extensions.Controller {
	controller := &commandExecutionController{
		BaseController: NewBaseController(
			models.SharedCommandModel(),
			logics.UniqueCommandExecutionLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *commandExecutionController) ResourceSingleURL() string {
	commandResourceName := extensions.RegisteredResourceName(models.SharedCommandModel())
	return fmt.Sprintf("%s/:id/execution", commandResourceName)
}

func (controller *commandExecutionController) RouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extensions.MethodPost: {
			extensions.URLSingle: controller.Create,
		},
		extensions.MethodDelete: {
			extensions.URLSingle: controller.Delete,
		},
	}
	return routeMap
}

var uniqueCommandExecutionController = newCommandExecutionController()

func init() {
	extensions.RegisterController(uniqueCommandExecutionController)
}
