package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type ephemeralBinaryObjectController struct {
	BaseController
}

func newEphemeralBinaryObjectController() *ephemeralBinaryObjectController {
	return CreateController(&ephemeralBinaryObjectController{}, model.NewEphemeralBinaryObject()).(*ephemeralBinaryObjectController)
}

func (receiver *ephemeralBinaryObjectController) GetResourceSingleURL() (string, error) {
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

func (receiver *ephemeralBinaryObjectController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
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
	extension.RegisterController(newEphemeralBinaryObjectController())
}
