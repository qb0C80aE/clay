package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"github.com/qb0C80aE/clay/version"
	"net/http"
)

type ephemeralScriptExecutionController struct {
	BaseController
}

func newEphemeralScriptExecutionController() *ephemeralScriptExecutionController {
	return CreateController(&ephemeralScriptExecutionController{}, model.NewEphemeralScriptExecution()).(*ephemeralScriptExecutionController)
}

func (receiver *ephemeralScriptExecutionController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	ephemeralScriptResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewEphemeralScript())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s/execution", ephemeralScriptResourceName, modelKey.KeyParameter), nil
}

func (receiver *ephemeralScriptExecutionController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodPost: {
			extension.URLSingle: receiver.Create,
		},
		extension.MethodDelete: {
			extension.URLSingle: receiver.Delete,
		},
	}
	return routeMap
}

func (receiver *ephemeralScriptExecutionController) Create(c *gin.Context) {
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

	// do not start transaction, handle it just as READ operation kile GetSingle
	// if transaction begins here, DB read operation in async script execution will fail after making commit or rollback in main task
	// script only can READ from db, so no need to begin transaction
	result, err := receiver.model.Create(receiver.model, db, c.Params, c.Request.URL.Query(), container)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputCreate(c, http.StatusCreated, result)
}

func init() {
	extension.RegisterController(newEphemeralScriptExecutionController())
}
