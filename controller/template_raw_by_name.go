package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type templateRawByNameController struct {
	BaseController
}

func newTemplateRawByNameController() *templateRawByNameController {
	return CreateController(&templateRawByNameController{}, model.NewTemplateRawByName()).(*templateRawByNameController)
}

func (receiver *templateRawByNameController) GetResourceSingleURL() (string, error) {
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

func (receiver *templateRawByNameController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *templateRawByNameController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	OutputTextWithType(c, code, result)
}

func init() {
	extension.RegisterController(newTemplateRawByNameController())
}
