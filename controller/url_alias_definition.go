package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	mapstructutilpkg "github.com/qb0C80aE/clay/util/mapstruct"
	"github.com/qb0C80aE/clay/version"
	"net/http"
	"net/url"
	"strings"
)

type urlAliasDefinitionController struct {
	BaseController
}

func newURLAliasDefinitionController() *urlAliasDefinitionController {
	return CreateController(&urlAliasDefinitionController{}, model.NewURLAliasDefinition()).(*urlAliasDefinitionController)
}

func (receiver *urlAliasDefinitionController) GetResourceSingleURL() (string, error) {
	modelKey, err := receiver.model.GetModelKey(receiver.model, "")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	resourceName, err := receiver.GetResourceName()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s", resourceName, modelKey.KeyParameter), nil
}

func (receiver *urlAliasDefinitionController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
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

func (receiver *urlAliasDefinitionController) Create(c *gin.Context) {
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

	newURLAliasDefinition := &model.URLAliasDefinition{}
	if err := mapstructutilpkg.GetUtility().MapToStruct(container, newURLAliasDefinition); err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	apiPath := ""
	newController := newURLAliasController(model.NewURLAlias())
	newController.name = newURLAliasDefinition.Name
	newController.from = strings.Trim(newURLAliasDefinition.From, "/")
	newController.to = strings.Trim(newURLAliasDefinition.To, "/")
	newController.methods = newURLAliasDefinition.Methods

	newControllerHandlerMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: newController.GetSingle,
			extension.URLMulti:  newController.GetMulti,
		},
		extension.MethodPost: {
			extension.URLSingle: newController.Create,
			extension.URLMulti:  newController.Create,
		},
		extension.MethodPut: {
			extension.URLSingle: newController.Update,
			extension.URLMulti:  newController.Update,
		},
		extension.MethodDelete: {
			extension.URLSingle: newController.Delete,
			extension.URLMulti:  newController.Delete,
		},
		extension.MethodPatch: {
			extension.URLSingle: newController.Patch,
			extension.URLMulti:  newController.Patch,
		},
		extension.MethodOptions: {
			extension.URLSingle: newController.GetOptions,
			extension.URLMulti:  newController.GetOptions,
		},
	}
	newController.routeMap = map[int]map[int]gin.HandlerFunc{}
	newController.acceptMap = map[int]map[int]*acceptSet{}
	for _, method := range newController.methods {
		var route map[int]gin.HandlerFunc
		var accept map[int]*acceptSet
		var exists bool

		methodInt := extension.LookUpMethodInt(method.Method)

		if route, exists = newController.routeMap[methodInt]; !exists {
			route = map[int]gin.HandlerFunc{}
		}
		if accept, exists = newController.acceptMap[methodInt]; !exists {
			accept = map[int]*acceptSet{}
		}
		switch method.TargetURLType {
		case "single":
			route[extension.URLSingle] = newControllerHandlerMap[methodInt][extension.URLSingle]
			newAcceptSet := &acceptSet{}
			if len(strings.Trim(method.Accept, " ")) == 0 {
				newAcceptSet.accept = extension.AcceptAll
			} else {
				newAcceptSet.accept = method.Accept
			}
			newAcceptSet.acceptCharset = method.AcceptCharset

			accept[extension.URLSingle] = newAcceptSet
		case "multi":
			route[extension.URLMulti] = newControllerHandlerMap[methodInt][extension.URLMulti]
			newAcceptSet := &acceptSet{}
			if len(strings.Trim(method.Accept, " ")) == 0 {
				newAcceptSet.accept = extension.AcceptAll
			} else {
				newAcceptSet.accept = method.Accept
			}
			newAcceptSet.acceptCharset = method.AcceptCharset

			accept[extension.URLMulti] = newAcceptSet
		}
		if (len(route) > 1) || (len(accept) > 1) {
			panic("an url alias cannot support both single and multi a method")
		}
		newController.routeMap[methodInt] = route
		newController.acceptMap[methodInt] = accept
	}

	if _, err := url.ParseQuery(newURLAliasDefinition.Query); err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}
	newController.query = newURLAliasDefinition.Query
	extension.RegisterController(newController)
	newControllerList := []extension.Controller{newController}
	if err := extension.SetupController(apiPath, newControllerList); err != nil {
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
	extension.RegisterController(newURLAliasDefinitionController())
}
