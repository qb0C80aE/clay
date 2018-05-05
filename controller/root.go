package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/http"
	"os"
	"sort"
	"strings"
)

// GetAPIEndpoints returns endpoints of Clay API
func GetAPIEndpoints(c *gin.Context) {
	reqScheme := "http"

	if c.Request.TLS != nil {
		reqScheme = "https"
	}

	reqHost := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s", reqScheme, reqHost)
	resourceList := []string{}

	controllerList := extension.GetRegisteredControllerList()
	for _, controller := range controllerList {
		routeMap := controller.GetRouteMap()
		for method, routes := range routeMap {
			resourceName, err := controller.GetResourceName()
			if err != nil {
				logging.Logger().Critical(err.Error())
				os.Exit(1)
			}
			modelKey, err := controller.GetModel().GetModelKey(controller.GetModel(), "")
			if err != nil {
				logging.Logger().Critical(err.Error())
				os.Exit(1)
			}
			title := fmt.Sprintf("%s_url [%s]", resourceName, extension.LookUpMethodName(method))
			for pathType := range routes {
				switch pathType {
				case extension.URLSingle:
					resourceSingleURL, err := controller.GetResourceSingleURL()
					if err != nil {
						logging.Logger().Critical(err.Error())
						os.Exit(1)
					}
					resourceSingleURL = strings.Replace(resourceSingleURL, ":key_parameter", fmt.Sprintf(":key_parameter(default=%s)", modelKey.KeyParameter), 1)
					resourceList = append(resourceList, fmt.Sprintf("%s %s/%s", title, baseURL, resourceSingleURL))
				case extension.URLMulti:
					resourceMultiURL, err := controller.GetResourceMultiURL()
					if err != nil {
						logging.Logger().Critical(err.Error())
						os.Exit(1)
					}
					resourceMultiURL = strings.Replace(resourceMultiURL, ":key_parameter", fmt.Sprintf(":key_parameter(default=%s)", modelKey.KeyParameter), 1)
					resourceList = append(resourceList, fmt.Sprintf("%s %s/%s", title, baseURL, resourceMultiURL))
				default:
					logging.Logger().Criticalf("invalid url type: %d", pathType)
					os.Exit(1)
				}
			}
		}
	}
	sort.Strings(resourceList)
	c.IndentedJSON(http.StatusOK, resourceList)
}
