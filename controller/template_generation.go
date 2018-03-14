package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type templateGenerationController struct {
	BaseController
}

func newTemplateGenerationController() *templateGenerationController {
	return CreateController(&templateGenerationController{}, model.NewTemplateGeneration()).(*templateGenerationController)
}

func (receiver *templateGenerationController) GetResourceSingleURL() (string, error) {
	templateResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewTemplate())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}
	return fmt.Sprintf("%s/:id/generation", templateResourceName), nil
}

func (receiver *templateGenerationController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *templateGenerationController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	text := result.(string)
	c.String(code, text)
}

func init() {
	extension.RegisterController(newTemplateGenerationController())
}
