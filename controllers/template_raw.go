package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type templateRawController struct {
	*BaseController
}

func newTemplateRawController() extensions.Controller {
	controller := &templateRawController{
		BaseController: NewBaseController(
			models.SharedTemplateModel(),
			logics.UniqueTemplateRawLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *templateRawController) ResourceSingleURL() string {
	templateResourceName := extensions.RegisteredResourceName(models.SharedTemplateModel())
	return fmt.Sprintf("%s/:id/raw", templateResourceName)
}

func (controller *templateRawController) RouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extensions.MethodGet: {
			extensions.URLSingle: controller.GetSingle,
		},
	}
	return routeMap
}

func (controller *templateRawController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	text := result.(string)
	c.String(code, text)
}

var uniqueTemplateRawController = newTemplateRawController()

func init() {
	extensions.RegisterController(uniqueTemplateRawController)
}
