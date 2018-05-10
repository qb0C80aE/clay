package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/http"
)

type urlAliasController struct {
	BaseController
	name  string
	from  string
	to    string
	query string
}

func newURLAliasController(model extension.Model) *urlAliasController {
	return CreateController(&urlAliasController{}, model).(*urlAliasController)
}

func (receiver *urlAliasController) GetResourceName() (string, error) {
	return fmt.Sprintf("%s [/%s : /%s] alias", receiver.name, receiver.from, receiver.to), nil
}

func (receiver *urlAliasController) GetResourceSingleURL() (string, error) {
	return receiver.from, nil
}

func (receiver *urlAliasController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *urlAliasController) GetSingle(c *gin.Context) {
	controller, err := extension.GetAssociatedControllerWithPath(receiver.to)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	routeMap := controller.GetRouteMap()
	methodGet, exists := routeMap[extension.MethodGet]
	if !exists {
		logging.Logger().Debug("the controller does not support GET")
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("the controller does not support GET"))
		return
	}

	handler, exists := methodGet[extension.URLSingle]
	if !exists {
		logging.Logger().Debug("the controller does not support single GET")
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("the controller does not support single GET"))
		return
	}

	singleURL, err := controller.GetResourceSingleURL()
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	c.Params, err = extension.CreateParametersFromPathAntRoute(receiver.to, singleURL)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	c.Request.URL.RawQuery = fmt.Sprintf("%s&%s", receiver.query, c.Request.URL.RawQuery)
	c.Request.URL.Path = receiver.to

	handler(c)
}

func init() {
}
