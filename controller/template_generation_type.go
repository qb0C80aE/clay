package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"net/http"
)

type templateGenerationTypeController struct {
	BaseController
}

func newTemplateGenerationTypeController() *templateGenerationTypeController {
	return CreateController(&templateGenerationTypeController{}, model.NewTemplateGeneration()).(*templateGenerationTypeController)
}

func (receiver *templateGenerationTypeController) GetResourceSingleURL() (string, error) {
	templateResourceName, err := extension.GetAssociatedResourceNameWithModel(model.NewTemplate())
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}
	return fmt.Sprintf("%s/:id/generation/:type", templateResourceName), nil
}

func (receiver *templateGenerationTypeController) GetRouteMap() map[int]map[int]gin.HandlerFunc {
	routeMap := map[int]map[int]gin.HandlerFunc{
		extension.MethodGet: {
			extension.URLSingle: receiver.GetSingle,
		},
	}
	return routeMap
}

func (receiver *templateGenerationTypeController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	text := result.(string)

	outputType := c.Param("type")

	switch outputType {
	case "text":
		c.String(code, text)
	case "html":
		c.Data(code, "text/html; charset=utf-8", []byte(text))
	default:
		c.JSON(http.StatusBadRequest, fmt.Errorf("invalid output type: %s", outputType))
	}
}

func init() {
	extension.RegisterController(newTemplateGenerationTypeController())
}
