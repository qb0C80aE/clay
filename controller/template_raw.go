package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type templateRawController struct {
	BaseController
}

func newTemplateRawController() *templateRawController {
	return CreateController(&templateRawController{}, model.NewTemplateRaw()).(*templateRawController)
}

func (receiver *templateRawController) GetResourceSingleURL() (string, error) {
	templateResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewTemplate())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:key_parameter/raw", templateResourceName), nil
}

func (receiver *templateRawController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *templateRawController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	receiver.BaseController.outputTextWithContentType(c, code, result)
}

func init() {
	extension.RegisterController(newTemplateRawController())
}
