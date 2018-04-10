package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/model"
)

type userDefinedManyToManyAssociationModelController struct {
	BaseController
}

func newUserDefinedManyToManyAssociationModelController(model extension.Model) *userDefinedManyToManyAssociationModelController {
	return CreateController(&userDefinedManyToManyAssociationModelController{}, model).(*userDefinedManyToManyAssociationModelController)
}

func (receiver *userDefinedManyToManyAssociationModelController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	model := receiver.model.(*model.UserDefinedManyToManyAssociationModel)
	var routeMap map[int]map[int]gin.HandlerFunc
	if model.ToBeMigrated() {
		routeMap = map[int]map[int]gin.HandlerFunc{
			extension.MethodGet: {
				extension.URLMulti: receiver.GetMulti,
			},
			extension.MethodPost: {
				extension.URLMulti: receiver.Create,
			},
		}
	} else {
		routeMap = map[int]map[int]gin.HandlerFunc{
			extension.MethodGet: {
				extension.URLMulti: receiver.GetMulti,
			},
		}
	}
	return routeMap
}

func init() {
}
