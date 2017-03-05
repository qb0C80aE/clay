package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type TemplateExternalParameterController struct {
	BaseController
}

type TemplateController struct {
	BaseController
}

func init() {
	extension.RegisterController(NewTemplateExternalParameterController())
	extension.RegisterController(NewTemplateController())
}

func NewTemplateExternalParameterController() *TemplateExternalParameterController {
	controller := &TemplateExternalParameterController{}
	controller.Initialize()
	return controller
}

func (this *TemplateExternalParameterController) Initialize() {
	this.ResourceName = "template_external_parameter"
	this.Model = models.TemplateExternalParameterModel
	this.Logic = logics.NewTemplateExternalParameterLogic()
	this.Outputter = this
}

func (this *TemplateExternalParameterController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
	resourceSingleUrl := extension.GetResourceSingleUrl(this.ResourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(this.ResourceName)

	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			resourceSingleUrl: this.GetSingle,
			resourceMultiUrl:  this.GetMulti,
		},
		extension.MethodPost: {
			resourceMultiUrl: this.Create,
		},
		extension.MethodPut: {
			resourceSingleUrl: this.Update,
		},
		extension.MethodDelete: {
			resourceSingleUrl: this.Delete,
		},
	}
	return routeMap
}

func NewTemplateController() *TemplateController {
	controller := &TemplateController{}
	controller.Initialize()
	return controller
}

func (this *TemplateController) Initialize() {
	this.ResourceName = "template"
	this.Model = models.TemplateModel
	this.Logic = logics.NewTemplateLogic()
	this.Outputter = this
}

func (this *TemplateController) GetResourceName() string {
	return this.ResourceName
}

func (this *TemplateController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
	resourceSingleUrl := extension.GetResourceSingleUrl(this.ResourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(this.ResourceName)

	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			resourceSingleUrl: this.GetSingle,
			resourceMultiUrl:  this.GetMulti,
		},
		extension.MethodPost: {
			resourceMultiUrl: this.Create,
		},
		extension.MethodPut: {
			resourceSingleUrl: this.Update,
		},
		extension.MethodDelete: {
			resourceSingleUrl: this.Delete,
		},
		extension.MethodPatch: {
			resourceSingleUrl: this.Patch,
		},
	}
	return routeMap
}

func (this *TemplateController) OutputPatch(c *gin.Context, code int, result interface{}) {
	text := result.(string)
	c.String(code, text)
}
