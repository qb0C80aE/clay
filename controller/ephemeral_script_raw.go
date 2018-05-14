package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type ephemeralScriptRawController struct {
	BaseController
}

func newEphemeralScriptRawController() *ephemeralScriptRawController {
	return CreateController(&ephemeralScriptRawController{}, model.NewEphemeralScriptRaw()).(*ephemeralScriptRawController)
}

func (receiver *ephemeralScriptRawController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	ephemeralScriptResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewEphemeralScript())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s/raw", ephemeralScriptResourceName, modelKey.KeyParameter), nil
}

func (receiver *ephemeralScriptRawController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *ephemeralScriptRawController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	OutputTextWithContentType(c, code, result)
}

func init() {
	extension.RegisterController(newEphemeralScriptRawController())
}
