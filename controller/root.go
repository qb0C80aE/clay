package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"net/http"
	"os"
	"sort"
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
			title := fmt.Sprintf("%s_url [%s]", controller.GetResourceName(), extension.LookUpMethodName(method))
			for pathType := range routes {
				switch pathType {
				case extension.URLSingle:
					resourceList = append(resourceList, fmt.Sprintf("%s %s/%s", title, baseURL, controller.GetResourceSingleURL()))
				case extension.URLMulti:
					resourceList = append(resourceList, fmt.Sprintf("%s %s/%s", title, baseURL, controller.GetResourceMultiURL()))
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
