package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"net/http"
)

type urlAliasController struct {
	BaseController
	name    string
	from    string
	to      string
	query   string
	methods []*model.URLAliasMethodURLTypeDefinition
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

func (receiver *urlAliasController) GetResourceMultiURL() (string, error) {
	return receiver.from, nil
}

func (receiver *urlAliasController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{}

	for _, method := range receiver.methods {
		switch method.Method {
		case "GET":
			var route map[int]gin.HandlerFunc
			var exists bool
			if route, exists = routeMap[extension.MethodGet]; !exists {
				route = map[int]gin.HandlerFunc{}
			}
			switch method.TargetURLType {
			case "single":
				route[extension.URLSingle] = receiver.GetSingle
			case "multi":
				route[extension.URLMulti] = receiver.GetMulti
			}
			if len(route) > 1 {
				panic("an url alias cannot support both single and multi a method")
			}
			routeMap[extension.MethodGet] = route
		case "POST":
			var route map[int]gin.HandlerFunc
			var exists bool
			if route, exists = routeMap[extension.MethodPost]; !exists {
				route = map[int]gin.HandlerFunc{}
			}
			switch method.TargetURLType {
			case "single":
				route[extension.URLSingle] = receiver.Create
			case "multi":
				route[extension.URLMulti] = receiver.Create
			}
			if len(route) > 1 {
				panic("an url alias cannot support both single and multi a method")
			}
			routeMap[extension.MethodPost] = route
		case "PUT":
			var route map[int]gin.HandlerFunc
			var exists bool
			if route, exists = routeMap[extension.MethodPut]; !exists {
				route = map[int]gin.HandlerFunc{}
			}
			switch method.TargetURLType {
			case "single":
				route[extension.URLSingle] = receiver.Update
			case "multi":
				route[extension.URLMulti] = receiver.Update
			}
			if len(route) > 1 {
				panic("an url alias cannot support both single and multi a method")
			}
			routeMap[extension.MethodPut] = route
		case "DELETE":
			var route map[int]gin.HandlerFunc
			var exists bool
			if route, exists = routeMap[extension.MethodDelete]; !exists {
				route = map[int]gin.HandlerFunc{}
			}
			switch method.TargetURLType {
			case "single":
				route[extension.URLSingle] = receiver.Delete
			case "multi":
				route[extension.URLMulti] = receiver.Delete
			}
			if len(route) > 1 {
				panic("an url alias cannot support both single and multi a method")
			}
			routeMap[extension.MethodDelete] = route
		case "PATCH":
			var route map[int]gin.HandlerFunc
			var exists bool
			if route, exists = routeMap[extension.MethodPatch]; !exists {
				route = map[int]gin.HandlerFunc{}
			}
			switch method.TargetURLType {
			case "single":
				route[extension.URLSingle] = receiver.Patch
			case "multi":
				route[extension.URLMulti] = receiver.Patch
			}
			if len(route) > 1 {
				panic("an url alias cannot support both single and multi a method")
			}
			routeMap[extension.MethodPatch] = route
		case "OPTIONS":
			var route map[int]gin.HandlerFunc
			var exists bool
			if route, exists = routeMap[extension.MethodOptions]; !exists {
				route = map[int]gin.HandlerFunc{}
			}
			switch method.TargetURLType {
			case "single":
				route[extension.URLSingle] = receiver.GetOptions
			case "multi":
				route[extension.URLMulti] = receiver.GetOptions
			}
			if len(route) > 1 {
				panic("an url alias cannot support both single and multi a method")
			}
			routeMap[extension.MethodOptions] = route
		}
	}

	return routeMap
}

func (receiver *urlAliasController) route(c *gin.Context, method int, selfURLType int) {
	controller, err := extension.GetAssociatedControllerWithPath(receiver.to)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	resourceName, err := controller.GetResourceName()
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	methodName := extension.LookUpMethodName(method)

	routeMap := controller.GetRouteMap()
	methodRoute, exists := routeMap[method]
	if !exists {
		logging.Logger().Debugf("the controller for %s does not support %s", resourceName, methodName)
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, fmt.Errorf("the controller for %s does not support %s", resourceName, methodName))
		return
	}

	switch selfURLType {
	case extension.URLSingle:
		handler, exists := methodRoute[selfURLType]
		if !exists {
			logging.Logger().Debug("the controller for %s does not support %s", resourceName, methodName)
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, fmt.Errorf("the controller for %s does not support %s", resourceName, methodName))
			return
		}

		url, err := controller.GetResourceSingleURL()
		if err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}

		c.Params, err = extension.CreateParametersFromPathAndRoute(receiver.to, url)
		if err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}

		c.Request.URL.RawQuery = fmt.Sprintf("%s&%s", receiver.query, c.Request.URL.RawQuery)
		c.Request.URL.Path = receiver.to

		handler(c)
	case extension.URLMulti:
		handler, exists := methodRoute[selfURLType]
		if !exists {
			logging.Logger().Debug("the controller does not support the given method")
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("the controller does not support the given method"))
			return
		}

		url, err := controller.GetResourceMultiURL()
		if err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}

		c.Params, err = extension.CreateParametersFromPathAndRoute(receiver.to, url)
		if err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}

		c.Request.URL.RawQuery = fmt.Sprintf("%s&%s", receiver.query, c.Request.URL.RawQuery)
		c.Request.URL.Path = receiver.to

		handler(c)
	default:
		logging.Logger().Debug("unknown url type")
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("unknown url type"))
	}
}

func (receiver *urlAliasController) GetSingle(c *gin.Context) {
	receiver.route(c, extension.MethodGet, extension.URLSingle)
}

func (receiver *urlAliasController) GetMulti(c *gin.Context) {
	receiver.route(c, extension.MethodGet, extension.URLMulti)
}

func (receiver *urlAliasController) Create(c *gin.Context) {
	selfRouteMap := receiver.GetRouteMap()
	var selfURLType int
	for key := range selfRouteMap[extension.MethodPost] {
		selfURLType = key
		break
	}
	receiver.route(c, extension.MethodPost, selfURLType)
}

func (receiver *urlAliasController) Update(c *gin.Context) {
	selfRouteMap := receiver.GetRouteMap()
	var selfURLType int
	for key := range selfRouteMap[extension.MethodPut] {
		selfURLType = key
		break
	}
	receiver.route(c, extension.MethodPut, selfURLType)
}

func (receiver *urlAliasController) Delete(c *gin.Context) {
	selfRouteMap := receiver.GetRouteMap()
	var selfURLType int
	for key := range selfRouteMap[extension.MethodDelete] {
		selfURLType = key
		break
	}
	receiver.route(c, extension.MethodDelete, selfURLType)
}

func (receiver *urlAliasController) Patch(c *gin.Context) {
	selfRouteMap := receiver.GetRouteMap()
	var selfURLType int
	for key := range selfRouteMap[extension.MethodPatch] {
		selfURLType = key
		break
	}
	receiver.route(c, extension.MethodPatch, selfURLType)
}

func (receiver *urlAliasController) GetOptions(c *gin.Context) {
	selfRouteMap := receiver.GetRouteMap()
	var selfURLType int
	for key := range selfRouteMap[extension.MethodOptions] {
		selfURLType = key
		break
	}
	receiver.route(c, extension.MethodOptions, selfURLType)
}

func init() {
}
