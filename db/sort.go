package db

import (
	"strings"

	"github.com/jinzhu/gorm"
)

func convertPrefixToQuery(sort string) string {
	if strings.HasPrefix(sort, "-") {
		return strings.TrimLeft(sort, "-") + " desc"
	}
	return strings.TrimLeft(sort, " ") + " asc"
}

// SortRecords set the sort method to the db
func (parameter *Parameter) SortRecords(db *gorm.DB) *gorm.DB {
	if parameter.Sort == "" {
		return db
	}

	for _, sort := range strings.Split(parameter.Sort, ",") {
		db = db.Order(convertPrefixToQuery(sort))
	}

	return db
}
