package db

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Need to avoid "Got error when connect database, the error is 'sql: unknown driver "sqlite3" (forgotten import?)'"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/serenize/snaker"
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

	registeredModels := extensions.RegisteredModels()
	db.AutoMigrate(registeredModels...)

	for _, model := range registeredModels {
		extensions.RegisterResourceName(model, db.NewScope(model).TableName())
	}

	return db
}

// Instance returns the connected db instance
func Instance(c *gin.Context) *gorm.DB {
	return c.MustGet("DB").(*gorm.DB)
}

// SetPreloads configures the preload settings of the connected db
func (parameter *Parameter) SetPreloads(db *gorm.DB) *gorm.DB {
	if parameter.Preloads == "" {
		return db
	}

	for _, preload := range strings.Split(parameter.Preloads, ",") {
		var a []string

		for _, s := range strings.Split(preload, ".") {
			a = append(a, snaker.SnakeToCamel(s))
		}

		db = db.Preload(strings.Join(a, "."))
	}

	return db
}
