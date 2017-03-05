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
	this.resourceName = "design"
}

func (this *DesignController) GetResourceName() string {
	return this.resourceName
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

func (_ *DesignController) GetSingle(c *gin.Context) {
	ProcessSingleGet(c, models.DesignModel, logics.GetDesign, OutputJsonError, OutputSingleJsonResult)
}

func (_ *DesignController) Update(c *gin.Context) {
	db.DBInstance(c).Exec("pragma foreign_keys = off;")
	ProcessUpdate(c, &models.Design{}, logics.UpdateDesign, OutputJsonError, OutputSingleJsonResult)
	db.DBInstance(c).Exec("pragma foreign_keys = on;")
}

func (_ *DesignController) Delete(c *gin.Context) {
	db.DBInstance(c).Exec("pragma foreign_keys = off;")
	ProcessDelete(c, logics.DeleteDesign, OutputJsonError, OutputNothing)
	db.DBInstance(c).Exec("pragma foreign_keys = on;")
}
