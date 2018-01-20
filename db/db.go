package db

import (
	"os"
	"strings"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite" // Need to avoid "Got error when connect database, the error is 'sql: unknown driver "sqlite3" (forgotten import?)'"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/logging"
	"github.com/serenize/snaker"
)

// Connect connects to its database and returns the instance
func Connect() *gorm.DB {
	dbMode := os.Getenv("DB_MODE")
	var dbPath string
	switch dbMode {
	case "memory":
		dbPath = ":memory:"
	case "", "file":
		if dbFilePath := os.Getenv("DB_FILE_PATH"); dbFilePath != "" {
			dbPath = dbFilePath
		} else {
			dbPath = "clay.db"
		}
	default:
		logging.Logger().Criticalf("invalid DB_MODE '%s'", dbMode)
		os.Exit(1)
	}

	db, err := gorm.Open("sqlite3", dbPath)

	if err != nil {
		logging.Logger().Criticalf("got an error when connect to the database, the error is '%v'", err)
		os.Exit(1)
	}

	if gin.IsDebugging() {
		db.LogMode(true)
	}

	db.Exec("pragma foreign_keys = on")

	registeredModelsToBeMigrated := extensions.RegisteredModelsToBeMigrated()
	if err := db.AutoMigrate(registeredModelsToBeMigrated...).Error; err != nil {
		logging.Logger().Criticalf("AutoMigration failed: '%s'", err.Error())
		os.Exit(1)
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
			logging.Logger().Criticalf("failed to run the initial data loader: %s", err)
			os.Exit(1)
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
					switch {
					case regexpLike.MatchString(v):
						db = db.Where(fmt.Sprintf("%s LIKE ?", columnName), v)
					case regexNotEquals.MatchString(v):
						parameter := v[1:]
						if parameter == "null" {
							db = db.Where(fmt.Sprintf("%s is not null", columnName))
						} else {
							db = db.Where(fmt.Sprintf("%s NOT IN (?)", columnName), strings.Split(parameter, ","))
						}
					default:
						if v == "null" {
							db = db.Where(fmt.Sprintf("%s is null", columnName))
						} else {
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
