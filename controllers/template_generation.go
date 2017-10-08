package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logics"
	"github.com/qb0C80aE/clay/models"
)

type templateGenerationController struct {
	*BaseController
}

func newTemplateGenerationController() extensions.Controller {
	controller := &templateGenerationController{
		BaseController: NewBaseController(
			models.SharedTemplateGenerationModel(),
			logics.UniqueTemplateGenerationLogic(),
		),
	}
	controller.SetOutputter(controller)
	return controller
}

func (controller *templateGenerationController) ResourceSingleURL() string {
	templateResourceName := extensions.RegisteredResourceName(models.SharedTemplateModel())
	return fmt.Sprintf("%s/:id/generation", templateResourceName)
}

func (controller *templateGenerationController) RouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extensions.MethodGet: {
			extensions.URLSingle: controller.GetSingle,
		},
	}
	return routeMap
}

func (controller *templateGenerationController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	text := result.(string)
	c.String(code, text)
}

var uniqueTemplateGenerationController = newTemplateGenerationController()

func init() {
	extensions.RegisterController(uniqueTemplateGenerationController)
}
