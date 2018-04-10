package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/model"
)

type templateArgumentController struct {
	BaseController
}

func newTemplateArgumentController() *templateArgumentController {
	return CreateController(&templateArgumentController{}, model.NewTemplateArgument()).(*templateArgumentController)
}

func (receiver *templateArgumentController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
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
	extension.RegisterController(newTemplateArgumentController())
}
