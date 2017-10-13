package db

import (
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/serenize/snaker"
	"net/url"
	"regexp"
)

func filterToMap(query url.Values) (map[string]string, map[string]map[string]string) {
	p, _ := regexp.Compile("^q\\[[a-zA-Z\\._]+\\]$")
	filters := make(map[string]string)
	preloadsFilterMap := make(map[string]map[string]string)
	for queryKey := range query {
		if p.MatchString(queryKey) {
			if strings.Contains(queryKey, ".") {
				s := queryKey[2 : len(queryKey)-1]
				a := strings.Split(s, ".")
				filterKey := a[len(a)-1]
				a = a[:len(a)-1]
				fieldName := strings.Join(a, ".")
				if _, exists := preloadsFilterMap[fieldName]; !exists {
					preloadsFilterMap[fieldName] = make(map[string]string)
				}
				preloadsFilterMap[fieldName][filterKey] = query.Get(queryKey)
			} else {
				filterKey := queryKey[2 : len(queryKey)-1]
				filters[filterKey] = query.Get(queryKey)
			}
		}
	}

	return filters, preloadsFilterMap
}

// FilterFields filters fields
func (parameter *Parameter) FilterFields(db *gorm.DB) *gorm.DB {
	for k, v := range parameter.Filters {
		if (v != "") && !(strings.Contains(k, ".")) {
			columnName := snaker.CamelToSnake(k)
			switch v {
			case "null":
				db = db.Where(fmt.Sprintf("%s is null", columnName))
			case "not_null":
				db = db.Where(fmt.Sprintf("%s is not null", columnName))
			default:
				db = db.Where(fmt.Sprintf("%s IN (?)", columnName), strings.Split(v, ","))
			}
		}
	}

	return db
}

// GetRawFilterQuery generates a filter query string
func (parameter *Parameter) GetRawFilterQuery() string {
	var s string

	for k, v := range parameter.Filters {
		if v != "" {
			s += "&q[" + k + "]=" + v
		}
	}

	return s
}
