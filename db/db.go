package db

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/serenize/snaker"
	"log"
	"os"
	"strings"
)

func Connect() *gorm.DB {
	dbMode := os.Getenv("DB_MODE")
	var dbPath string
	switch dbMode {
	case "", "memory":
		dbPath = ":memory:"
	case "file":
		if dbFilePath := os.Getenv("DB_FILE_PATH"); dbFilePath != "" {
			dbPath = dbFilePath
		} else {
			dbPath = "clay.db"
		}
	default:
		log.Fatalf("Invalid DB_MODE '%s'", dbMode)
	}

	db, err := gorm.Open("sqlite3", dbPath)

	if err != nil {
		log.Fatalf("Got error when connect database, the error is '%v'", err)
	}

	db.Exec("pragma foreign_keys = on")
	db.LogMode(true)

	if gin.IsDebugging() {
		db.LogMode(true)
	}

	db.AutoMigrate(extensions.GetModels()...)

	return db
}

func DBInstance(c *gin.Context) *gorm.DB {
	return c.MustGet("DB").(*gorm.DB)
}

func SetPreloads(preloads string, db *gorm.DB) *gorm.DB {
	if preloads == "" {
		return db
	}

	for _, preload := range strings.Split(preloads, ",") {
		var a []string

		for _, s := range strings.Split(preload, ".") {
			a = append(a, snaker.SnakeToCamel(s))
		}

		db = db.Preload(strings.Join(a, "."))
	}

	return db
}
