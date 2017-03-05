package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type DesignController struct {
	BaseController
}

func init() {
	extension.RegisterController(NewDesignController())
}

func NewDesignController() *DesignController {
	controller := &DesignController{}
	controller.Initialize()
	return controller
}

func (this *DesignController) Initialize() {
	this.ResourceName = "design"
	this.Model = models.DesignModel
	this.Logic = logics.NewDesignLogic()
	this.Outputter = this
}

func (this *DesignController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
	url := "designs/present"
	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			url: this.GetSingle,
		},
		extension.MethodPut: {
			url: this.Update,
		},
		extension.MethodDelete: {
			url: this.Delete,
		},
	}
	return routeMap
}

func (this *DesignController) Update(c *gin.Context) {
	db.DBInstance(c).Exec("pragma foreign_keys = off;")
	this.BaseController.Update(c)
	db.DBInstance(c).Exec("pragma foreign_keys = on;")
}

func (this *DesignController) Delete(c *gin.Context) {
	db.DBInstance(c).Exec("pragma foreign_keys = off;")
	this.BaseController.Delete(c)
	db.DBInstance(c).Exec("pragma foreign_keys = on;")
}
