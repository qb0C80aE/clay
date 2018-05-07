package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type ephemeralTemplateGenerationController struct {
	BaseController
}

func newEphemeralTemplateGenerationController() *ephemeralTemplateGenerationController {
	return CreateController(&ephemeralTemplateGenerationController{}, model.NewEphemeralTemplateGeneration()).(*ephemeralTemplateGenerationController)
}

func (receiver *ephemeralTemplateGenerationController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	ephemeralTemplateResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewEphemeralTemplate())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s/generation", ephemeralTemplateResourceName, modelKey.KeyParameter), nil
}

func (receiver *ephemeralTemplateGenerationController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *ephemeralTemplateGenerationController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	OutputTextWithContentType(c, code, result)
}

func init() {
	extension.RegisterController(newEphemeralTemplateGenerationController())
}
