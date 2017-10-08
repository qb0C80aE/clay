package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/db"
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
	controller.SetOutputter(controller)
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

// Caution: Even if you input the inconsistent data like foreign keys do not exist,
//          it will be registered, and never be checked this time.
//          Todo: It requires order resolution logic like "depends on" between models.
func (controller *designController) Update(c *gin.Context) {
	db.Instance(c).Exec("pragma foreign_keys = off;")
	controller.BaseController.Update(c)
	db.Instance(c).Exec("pragma foreign_keys = on;")
}

func (controller *designController) Delete(c *gin.Context) {
	db.Instance(c).Exec("pragma foreign_keys = off;")
	controller.BaseController.Delete(c)
	db.Instance(c).Exec("pragma foreign_keys = on;")
}

var uniqueDesignController = newDesignController()

func init() {
	extensions.RegisterController(uniqueDesignController)
}
