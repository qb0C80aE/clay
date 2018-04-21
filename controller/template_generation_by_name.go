package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type templateGenerationByNameController struct {
	BaseController
}

func newTemplateGenerationByNameController() *templateGenerationByNameController {
	return CreateController(&templateGenerationByNameController{}, model.NewTemplateGenerationByName()).(*templateGenerationByNameController)
}

func (receiver *templateGenerationByNameController) GetResourceSingleURL() (string, error) {
	resourceName, err := receiver.GetResourceName()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	modelKey, err := extension.GetRegisteredModelKey(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s", resourceName, modelKey.KeyParameter), nil
}

func (receiver *templateGenerationByNameController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *templateGenerationByNameController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	OutputTextWithType(c, code, result)
}

func init() {
	extension.RegisterController(newTemplateGenerationByNameController())
}
