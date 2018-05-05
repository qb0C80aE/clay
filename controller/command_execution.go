package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"github.com/qb0C80aE/clay/version"
	"net/http"
)

type commandExecutionController struct {
	BaseController
}

func newCommandExecutionController() *commandExecutionController {
	return CreateController(&commandExecutionController{}, model.NewCommandExecution()).(*commandExecutionController)
}

func (receiver *commandExecutionController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	commandResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewCommand())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s/execution", commandResourceName, modelKey.KeyParameter), nil
}

func (receiver *commandExecutionController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
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

func (receiver *commandExecutionController) Create(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	// don't bind, just call Create of the model. the Content-Type and Body have nothing
	// The result is Command
	result, err := receiver.model.Create(receiver.model, nil, c.Params, c.Request.URL.Query(), nil)
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
	extension.RegisterController(newCommandExecutionController())
}
