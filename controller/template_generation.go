package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
)

type templateGenerationController struct {
	BaseController
}

func newTemplateGenerationController() *templateGenerationController {
	return CreateController(&templateGenerationController{}, model.NewTemplateGeneration()).(*templateGenerationController)
}

func (receiver *templateGenerationController) GetResourceSingleURL() (string, error) {
	templateResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewTemplate())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	modelKey, err := extension.GetRegisteredModelKey(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:%s/generation", templateResourceName, modelKey.KeyParameter), nil
}

func (receiver *templateGenerationController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *templateGenerationController) OutputGetSingle(c *gin.Context, code int, result interface{}, _ map[string]interface{}) {
	OutputTextWithType(c, code, result)
}

// OutputTextWithType outputs the result text in the given type
func OutputTextWithType(c *gin.Context, code int, result interface{}) {
	text := result.(string)

	outputType := c.Request.URL.Query().Get("output_type")

	switch outputType {
	case "html":
		c.Data(code, "text/html; charset=utf-8", []byte(text))
	case "js":
		c.Data(code, "application/javascript; charset=utf-8", []byte(text))
	case "css":
		c.Data(code, "text/css; charset=utf-8", []byte(text))
	case "json":
		c.Data(code, "application/json; charset=utf-8", []byte(text))
	case "yaml":
		c.Data(code, "application/x-yaml; charset=utf-8", []byte(text))
	default:
		c.String(code, text)
	}
}

func init() {
	extension.RegisterController(newTemplateGenerationController())
}
