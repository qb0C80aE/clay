package db

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Need to avoid "Got error when connect database, the error is 'sql: unknown driver "sqlite3" (forgotten import?)'"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/serenize/snaker"
	"log"
	"os"
	"strings"
)

// Connect connects to its database and returns the instance
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

	db.AutoMigrate(extensions.RegisteredModels()...)

	return db
}

// Instance returns the connected db instance
func Instance(c *gin.Context) *gorm.DB {
	return c.MustGet("DB").(*gorm.DB)
}

// SetPreloads configures the preload settings of the connected db
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
