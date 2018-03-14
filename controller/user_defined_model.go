package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/model"
)

type userDefinedModelController struct {
	BaseController
}

func newUserDefinedModelController(model extension.Model) *userDefinedModelController {
	return CreateController(&userDefinedModelController{}, model).(*userDefinedModelController)
}

func (receiver *userDefinedModelController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	model := receiver.model.(*model.UserDefinedModel)
	var routeMap map[int]map[int]gin.HandlerFunc
	if model.ToBeMigrated() {
		routeMap = map[int]map[int]gin.HandlerFunc{
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
	} else {
		routeMap = map[int]map[int]gin.HandlerFunc{
			extension.MethodGet: {
				extension.URLSingle: receiver.GetSingle,
				extension.URLMulti:  receiver.GetMulti,
			},
		}
	}
	return routeMap
}

func init() {
}
