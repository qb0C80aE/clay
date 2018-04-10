package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"github.com/qb0C80aE/clay/version"
	"net/http"
	"reflect"
)

type userDefinedModelDefinitionController struct {
	BaseController
}

func newUserDefinedModelRegistrationController() *userDefinedModelDefinitionController {
	return CreateController(&userDefinedModelDefinitionController{}, model.NewUserDefinedModelDefinition()).(*userDefinedModelDefinitionController)
}

func (receiver *userDefinedModelDefinitionController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
			extension.URLMulti:  receiver.GetMulti,
		},
		extension.MethodPost: {
			extension.URLMulti: receiver.Create,
		},
	}
	return routeMap
}

func (receiver *userDefinedModelDefinitionController) Create(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	resourceName, err := receiver.model.GetResourceName(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	container, err := receiver.binder.Bind(c, resourceName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)
	result, err := receiver.model.Create(receiver.model, db, c.Params, c.Request.URL.Query(), container)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	typeName := reflect.ValueOf(container).Elem().FieldByName("TypeName").Interface().(string)
	newModel, err := extension.GetAssociatedModelWithTypeName(typeName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	newUserDefinedModel, ok := newModel.(*model.UserDefinedModel)
	if !ok {
		logging.Logger().Debug("model is not an UserDefinedModel")
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("model is not an UserDefinedModel"))
		return
	}

	if newUserDefinedModel.IsControllerEnabled() {
		apiPath := ""
		newController := newUserDefinedModelController(newModel)
		extension.RegisterController(newController)
		newControllerList := []extension.Controller{newController}
		if err := extension.SetupController(apiPath, newControllerList); err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputCreate(c, http.StatusCreated, result)
}

func init() {
	extension.RegisterController(newUserDefinedModelRegistrationController())
}
