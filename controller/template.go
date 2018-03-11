package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/model"
)

type templateController struct {
	BaseController
}

func newTemplateController() *templateController {
	return CreateController(&templateController{}, model.NewTemplate()).(*templateController)
}

func (receiver *templateController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
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
	extension.RegisterController(newTemplateController())
}
