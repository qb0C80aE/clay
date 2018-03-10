package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/model"
)

type templateRawController struct {
	*BaseController
}

func newTemplateRawController() *templateRawController {
	return CreateController(&templateRawController{}, model.NewTemplateRaw()).(*templateRawController)
}

func (receiver *templateRawController) GetResourceSingleURL() string {
	templateResourceName := extension.GetAssociateResourceNameWithModel(model.NewTemplate())
	return fmt.Sprintf("%s/:id/raw", templateResourceName)
}

func (receiver *templateRawController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *templateRawController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	text := result.(string)
	c.String(code, text)
}

func init() {
	extension.RegisterController(newTemplateRawController())
}
