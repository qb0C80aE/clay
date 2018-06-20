package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"net/http"
	"strings"
)

type acceptSet struct {
	accept        string
	acceptCharset string
}

type urlAliasController struct {
	BaseController
	name      string
	from      string
	to        string
	query     string
	methods   []*model.URLAliasMethodURLTypeDefinition
	routeMap  map[int]map[int]gin.HandlerFunc
	acceptMap map[int]map[int]*acceptSet
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
	return receiver.routeMap
}

func (receiver *urlAliasController) route(c *gin.Context, method int, selfURLType int) {
	targetController, err := extension.GetAssociatedControllerWithPath(receiver.to)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	resourceName, err := targetController.GetResourceName()
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	methodName := extension.LookUpMethodName(method)

	methodRoute, exists := targetController.GetRouteMap()[method]
	if !exists {
		logging.Logger().Debugf("the controller for %s does not support %s", resourceName, methodName)
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, fmt.Errorf("the controller for %s does not support %s", resourceName, methodName))
		return
	}

	methodAccept, exists := receiver.acceptMap[method]
	if !exists {
		logging.Logger().Debugf("the controller for %s does not support %s", resourceName, methodName)
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, fmt.Errorf("the controller for %s does not support %s", resourceName, methodName))
		return
	}

	var url string
	switch selfURLType {
	case extension.URLSingle:
		url, err = targetController.GetResourceSingleURL()
		if err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}
	case extension.URLMulti:
		url, err = targetController.GetResourceMultiURL()
		if err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}
	default:
		logging.Logger().Debug("unknown url type")
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("unknown url type"))
		return
	}

	handler, exists := methodRoute[selfURLType]
	if !exists {
		logging.Logger().Debugf("the controller for %s does not support %s", resourceName, methodName)
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, fmt.Errorf("the controller for %s does not support %s", resourceName, methodName))
		return
	}

	c.Params, err = extension.CreateParametersFromPathAndRoute(receiver.to, url)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	targetAcceptSet, exists := methodAccept[selfURLType]
	if !exists {
		logging.Logger().Debug("the controller for %s does not support %s", resourceName, methodName)
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, fmt.Errorf("the controller for %s does not support %s", resourceName, methodName))
		return
	}

	var acceptList []string
	originalAcceptList := strings.Split(c.Request.Header.Get("Accept"), ",")
	switch len(originalAcceptList) {
	case 0:
		acceptList = strings.Split(targetAcceptSet.accept, ",")
	case 1:
		if originalAcceptList[0] == extension.AcceptAll {
			acceptList = strings.Split(targetAcceptSet.accept, ",")
		} else {
			if targetAcceptSet.accept == extension.AcceptAll {
				acceptList = append(originalAcceptList, targetAcceptSet.accept)
			} else {
				acceptList = append(strings.Split(targetAcceptSet.accept, ","), originalAcceptList...)
			}
		}
	default:
		if targetAcceptSet.accept == extension.AcceptAll {
			acceptList = append(originalAcceptList, targetAcceptSet.accept)
		} else {
			acceptList = append(strings.Split(targetAcceptSet.accept, ","), originalAcceptList...)
		}
	}

	var acceptCharsetList []string
	originalAcceptCharsetList := strings.Split(c.Request.Header.Get("Accept-Charset"), ",")
	switch len(originalAcceptCharsetList) {
	case 0:
		acceptCharsetList = strings.Split(targetAcceptSet.acceptCharset, ",")
	case 1:
		if originalAcceptCharsetList[0] == extension.AcceptAll {
			acceptCharsetList = strings.Split(targetAcceptSet.acceptCharset, ",")
		} else {
			if targetAcceptSet.acceptCharset == extension.AcceptAll {
				acceptCharsetList = append(originalAcceptCharsetList, targetAcceptSet.acceptCharset)
			} else {
				acceptCharsetList = append(strings.Split(targetAcceptSet.acceptCharset, ","), originalAcceptCharsetList...)
			}
		}
	default:
		if targetAcceptSet.accept == extension.AcceptAll {
			acceptCharsetList = append(originalAcceptCharsetList, targetAcceptSet.acceptCharset)
		} else {
			acceptCharsetList = append(strings.Split(targetAcceptSet.acceptCharset, ","), originalAcceptCharsetList...)
		}
	}

	c.Request.URL.RawQuery = fmt.Sprintf("%s&%s", receiver.query, c.Request.URL.RawQuery)
	c.Request.URL.Path = receiver.to
	c.Request.Header.Set("Accept", strings.Join(acceptList, ","))
	c.Request.Header.Set("Accept-Charset", strings.Join(acceptCharsetList, ","))

	handler(c)
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
