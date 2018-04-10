package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type userDefinedManyToManyAssociationChildModelController struct {
	BaseController
	left bool
}

func newUserDefinedManyToManyAssociationChildModelController(model extension.Model, left bool) *userDefinedManyToManyAssociationChildModelController {
	controller := CreateController(&userDefinedManyToManyAssociationChildModelController{}, model).(*userDefinedManyToManyAssociationChildModelController)
	controller.left = left
	return controller
}

// GetResourceSingleURL builds a resource url what represents a single resource based on the argument
func (receiver *userDefinedManyToManyAssociationChildModelController) GetResourceSingleURL() (string, error) {
	userDefinedManyToManyAssociationModel, ok := receiver.model.(*model.UserDefinedManyToManyAssociationModel)
	if !ok {
		logging.Logger().Debug("the actual model is not *UserDefinedManyToManyAssociationModel")
		return "", errors.New("the actual model is not *UserDefinedManyToManyAssociationModel")
	}

	resourceName, err := receiver.GetResourceName()
	if err != nil {
		return "", err
	}

	url := ""
	if receiver.left {
		url = fmt.Sprintf("%s/%s/:%s/:%s",
			resourceName,
			userDefinedManyToManyAssociationModel.GetLeftResourceName(),
			userDefinedManyToManyAssociationModel.GetLeftResourceKeyName(),
			userDefinedManyToManyAssociationModel.GetRightResourceKeyName(),
		)
	} else {
		url = fmt.Sprintf("%s/%s/:%s/:%s",
			resourceName,
			userDefinedManyToManyAssociationModel.GetRightResourceName(),
			userDefinedManyToManyAssociationModel.GetRightResourceKeyName(),
			userDefinedManyToManyAssociationModel.GetLeftResourceKeyName(),
		)
	}

	return url, nil
}

// GetResourceMultiURL builds a resource url what represents multi resources based on the argument
func (receiver *userDefinedManyToManyAssociationChildModelController) GetResourceMultiURL() (string, error) {
	userDefinedManyToManyAssociationModel, ok := receiver.model.(*model.UserDefinedManyToManyAssociationModel)
	if !ok {
		logging.Logger().Debug("the actual model is not *UserDefinedManyToManyAssociationModel")
		return "", errors.New("the actual model is not *UserDefinedManyToManyAssociationModel")
	}

	resourceName, err := receiver.GetResourceName()
	if err != nil {
		return "", err
	}

	url := ""
	if receiver.left {
		url = fmt.Sprintf("%s/%s/:%s",
			resourceName,
			userDefinedManyToManyAssociationModel.GetLeftResourceName(),
			userDefinedManyToManyAssociationModel.GetLeftResourceKeyName(),
		)
	} else {
		url = fmt.Sprintf("%s/%s/:%s",
			resourceName,
			userDefinedManyToManyAssociationModel.GetRightResourceName(),
			userDefinedManyToManyAssociationModel.GetRightResourceKeyName(),
		)
	}

	return url, nil
}

func (receiver *userDefinedManyToManyAssociationChildModelController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	model := receiver.model.(*model.UserDefinedManyToManyAssociationModel)
	var routeMap map[int]map[int]gin.HandlerFunc
	if model.ToBeMigrated() {
		routeMap = map[int]map[int]gin.HandlerFunc{
			extension.MethodGet: {
				extension.URLMulti: receiver.GetMulti,
			},
			extension.MethodDelete: {
				extension.URLSingle: receiver.Delete,
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
