package db

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Paginate sets the parameter related to pagination to db
func (parameter *Parameter) Paginate(db *gorm.DB) (*gorm.DB, error) {
	if parameter == nil {
		return nil, errors.New("Parameter struct got nil")
	}

	if parameter.IsLastID {
		if parameter.Order == "asc" {
			return db.Where("id > ?", parameter.LastID).Limit(parameter.Limit).Order("id asc"), nil
		}

		return db.Where("id < ?", parameter.LastID).Limit(parameter.Limit).Order("id desc"), nil
	}

	return db.Offset(parameter.Limit * (parameter.Page - 1)).Limit(parameter.Limit), nil
}

// SetHeaderLink sets the header links to HTTP header
func (parameter *Parameter) SetHeaderLink(c *gin.Context, index int) error {
	if parameter == nil {
		return errors.New("Parameter struct got nil")
	}

	var pretty, filters, preloads string
	reqScheme := "http"

	if c.Request.TLS != nil {
		reqScheme = "https"
	}

	if _, ok := c.GetQuery("pretty"); ok {
		pretty = "&pretty"
	}

	if len(parameter.Filters) != 0 {
		filters = parameter.GetRawFilterQuery()
	}

	if parameter.Preloads != "" {
		preloads = fmt.Sprintf("&preloads=%v", parameter.Preloads)
	}

	if parameter.IsLastID {
		c.Header("Link", fmt.Sprintf("<%s://%v%v?limit=%v%s%s&last_id=%v&order=%v%s>; rel=\"next\"", reqScheme, c.Request.Host, c.Request.URL.Path, parameter.Limit, filters, preloads, index, parameter.Order, pretty))
		return nil
	}

	if parameter.Page == 1 {
		c.Header("Link", fmt.Sprintf("<%s://%v%v?limit=%v%s%s&page=%v%s>; rel=\"next\"", reqScheme, c.Request.Host, c.Request.URL.Path, parameter.Limit, filters, preloads, parameter.Page+1, pretty))
		return nil
	}

	c.Header("Link", fmt.Sprintf(
		"<%s://%v%v?limit=%v%s%s&page=%v%s>; rel=\"next\",<%s://%v%v?limit=%v%s%s&page=%v%s>; rel=\"prev\"", reqScheme,
		c.Request.Host, c.Request.URL.Path, parameter.Limit, filters, preloads, parameter.Page+1, pretty, reqScheme, c.Request.Host, c.Request.URL.Path, parameter.Limit, filters, preloads, parameter.Page-1, pretty))
	return nil
}
