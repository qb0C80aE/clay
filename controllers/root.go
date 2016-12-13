package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func APIEndpoints(c *gin.Context) {
	reqScheme := "http"

	if c.Request.TLS != nil {
		reqScheme = "https"
	}

	reqHost := c.Request.Host
	baseURL := fmt.Sprintf("%s://%s/%s", reqScheme, reqHost, "v1")

	resources := map[string]string{
		"nodes_url":       baseURL + "/nodes",
		"node_url":        baseURL + "/nodes/{id}",
		"node_groups_url": baseURL + "/node_groups",
		"node_group_url":  baseURL + "/node_groups/{id}",
		"node_pvs_url":    baseURL + "/node_pvs",
		"node_pv_url":     baseURL + "/node_pvs/{id}",
		"node_types_url":  baseURL + "/node_types",
		"node_type_url":   baseURL + "/node_types/{id}",
		"ports_url":       baseURL + "/ports",
		"port_url":        baseURL + "/ports/{id}",
		"designs_url":     baseURL + "/designs/present",
		"segments_url":    baseURL + "/segments",
		"diagrams_url":    baseURL + "/diagrams/{type}",
	}

	c.IndentedJSON(http.StatusOK, resources)
}
