package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type designController struct {
	*BaseController
}

func newDesignController() extension.Controller {
	controller := &designController{
		BaseController: NewBaseController(
			"design",
			models.SharedDesignModel(),
			logics.UniqueDesignLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *designController) RouteMap() map[int]map[string]gin.HandlerFunc {
	url := "designs/present"
	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			url: controller.GetSingle,
		},
		extension.MethodPut: {
			url: controller.Update,
		},
		extension.MethodDelete: {
			url: controller.Delete,
		},
	}
	return routeMap
}

func (controller *designController) Update(c *gin.Context) {
	db.DBInstance(c).Exec("pragma foreign_keys = off;")
	controller.BaseController.Update(c)
	db.DBInstance(c).Exec("pragma foreign_keys = on;")
}

func (controller *designController) Delete(c *gin.Context) {
	db.DBInstance(c).Exec("pragma foreign_keys = off;")
	controller.BaseController.Delete(c)
	db.DBInstance(c).Exec("pragma foreign_keys = on;")
}

var uniqueDesignController = newDesignController()

func init() {
	extension.RegisterController(uniqueDesignController)
}
