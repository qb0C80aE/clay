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
	this.resourceName = "template_external_parameter"
}

func (this *TemplateExternalParameterController) GetResourceName() string {
	return this.resourceName
}

func (this *TemplateExternalParameterController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
	resourceSingleUrl := extension.GetResourceSingleUrl(this.resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(this.resourceName)

	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			resourceMultiUrl:  this.GetSingle,
			resourceSingleUrl: this.GetMulti,
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
	this.resourceName = "template"
}

func (this *TemplateController) GetResourceName() string {
	return this.resourceName
}

func (this *TemplateController) GetRouteMap() map[int]map[string]gin.HandlerFunc {
	resourceSingleUrl := extension.GetResourceSingleUrl(this.resourceName)
	resourceMultiUrl := extension.GetResourceMultiUrl(this.resourceName)

	routeMap := map[int]map[string]gin.HandlerFunc{
		extension.MethodGet: {
			resourceMultiUrl:  this.GetSingle,
			resourceSingleUrl: this.GetMulti,
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

func (_ *TemplateExternalParameterController) GetMulti(c *gin.Context) {
	ProcessMultiGet(c, models.TemplateExternalParameterModel, logics.GetTemplateExternalParameters, OutputJsonError, OutputMultiJsonResult)
}

func (_ *TemplateExternalParameterController) GetSingle(c *gin.Context) {
	ProcessSingleGet(c, models.TemplateExternalParameterModel, logics.GetTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func (_ *TemplateExternalParameterController) Create(c *gin.Context) {
	ProcessCreate(c, &models.TemplateExternalParameter{}, logics.CreateTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func (_ *TemplateExternalParameterController) Update(c *gin.Context) {
	ProcessUpdate(c, &models.TemplateExternalParameter{}, logics.UpdateTemplateExternalParameter, OutputJsonError, OutputSingleJsonResult)
}

func (_ *TemplateExternalParameterController) Delete(c *gin.Context) {
	ProcessDelete(c, logics.DeleteTemplateExternalParameter, OutputJsonError, OutputNothing)
}

func (_ *TemplateController) GetMulti(c *gin.Context) {
	ProcessMultiGet(c, models.TemplateModel, logics.GetTemplates, OutputJsonError, OutputMultiJsonResult)
}

func (_ *TemplateController) GetSingle(c *gin.Context) {
	ProcessSingleGet(c, models.TemplateModel, logics.GetTemplate, OutputJsonError, OutputSingleJsonResult)
}

func (_ *TemplateController) Create(c *gin.Context) {
	ProcessCreate(c, &models.Template{}, logics.CreateTemplate, OutputJsonError, OutputNothing)
}

func (_ *TemplateController) Update(c *gin.Context) {
	ProcessUpdate(c, &models.Template{}, logics.UpdateTemplate, OutputJsonError, OutputNothing)
}

func (_ *TemplateController) Delete(c *gin.Context) {
	ProcessDelete(c, logics.DeleteTemplate, OutputJsonError, OutputNothing)
}

func (_ *TemplateController) Patch(c *gin.Context) {
	ProcessSingleGet(c, &models.Template{}, logics.ApplyTemplate, OutputJsonError, OutputTextResult)
}
