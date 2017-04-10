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

// SortRecords configures the order of retrieved records from db
func SortRecords(sorts string, db *gorm.DB) *gorm.DB {
	if sorts == "" {
		return db
	}

	for _, sort := range strings.Split(sorts, ",") {
		db = db.Order(convertPrefixToQuery(sort))
	}

	return db
}
