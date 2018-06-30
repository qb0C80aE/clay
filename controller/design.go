package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type designController struct {
	BaseController
}

func newDesignController() *designController {
	return CreateController(&designController{}, model.NewDesign()).(*designController)
}

func (receiver *designController) GetResourceSingleURL() (string, error) {
	resourceName, err := receiver.model.GetResourceName(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}
	return fmt.Sprintf("%s/present", resourceName), nil
}

func (receiver *designController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
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
	extension.RegisterController(newDesignController())
}
