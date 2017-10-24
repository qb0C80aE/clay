package db

import (
	"log"
	"os"
	"strings"

	"fmt"
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

	if gin.IsDebugging() {
		db.LogMode(true)
	}

	db.Exec("pragma foreign_keys = on")

	registeredModelsToBeMigrated := extensions.RegisteredModelsToBeMigrated()
	if err := db.AutoMigrate(registeredModelsToBeMigrated...).Error; err != nil {
		log.Fatalf("AutoMigration failed: '%s'", err.Error())
	}

	for _, model := range extensions.RegisteredModels() {
		tableName := db.NewScope(model).TableName()
		extensions.RegisterResourceName(model, tableName)
	}

	initialDataLoaders := extensions.RegisteredInitialDataLoaders()
	// Caution: Even if you input the inconsistent data like foreign keys do not exist,
	//          it will be registered, and never be checked this time.
	//          Todo: It requires order resolution logic like "depends on" between models.
	db.Exec("pragma foreign_keys = off;")

	tx := db.Begin()
	for _, initialDataLoader := range initialDataLoaders {
		err := initialDataLoader.SetupInitialData(tx)
		if err != nil {
			tx.Rollback()
			log.Fatalf("Failed to load the initial data: %s", err)
		}
	}
	tx.Commit()

	db.Exec("pragma foreign_keys = on;")

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

		if m, exists := parameter.PreloadsFilterMap[preload]; exists {
			db = db.Preload(strings.Join(a, "."), func(db *gorm.DB) *gorm.DB {
				for k, v := range m {
					columnName := snaker.CamelToSnake(k)
					if strings.Contains(v, "%") {
						db = db.Where(fmt.Sprintf("%s LIKE ?", columnName), v)
					} else {
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
			})
		} else {
			db = db.Preload(strings.Join(a, "."))
		}
	}

	return db
}
