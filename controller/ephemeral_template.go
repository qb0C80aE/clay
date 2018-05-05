package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type ephemeralTemplateController struct {
	BaseController
}

func newEphemeralTemplateController() *ephemeralTemplateController {
	return CreateController(&ephemeralTemplateController{}, model.NewEphemeralTemplate()).(*ephemeralTemplateController)
}

func (receiver *ephemeralTemplateController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	resourceName, err := receiver.GetResourceName()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s", resourceName, modelKey.KeyParameter), nil
}

func (receiver *ephemeralTemplateController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
			extension.URLMulti:  receiver.GetMulti,
		},
		extension.MethodPost: {
			extension.URLMulti: receiver.Create,
		},
		extension.MethodPut: {
			extension.URLSingle: receiver.Update,
		},
		extension.MethodDelete: {
			extension.URLSingle: receiver.Delete,
		},
	}
	return routeMap
}

func init() {
	extension.RegisterController(newEphemeralTemplateController())
}
