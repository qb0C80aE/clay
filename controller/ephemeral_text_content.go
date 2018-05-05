package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type ephemeralTextContentController struct {
	BaseController
}

func newEphemeralTextContentController() *ephemeralTextContentController {
	return CreateController(&ephemeralTextContentController{}, model.NewEphemeralTextContent()).(*ephemeralTextContentController)
}

func (receiver *ephemeralTextContentController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	ephemeralTextResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewEphemeralText())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s/content", ephemeralTextResourceName, modelKey.KeyParameter), nil
}

func (receiver *ephemeralTextContentController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *ephemeralTextContentController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	OutputTextWithType(c, code, result)
}

func init() {
	extension.RegisterController(newEphemeralTextContentController())
}
