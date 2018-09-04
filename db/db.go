package db

import (
	"strings"

	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/sqlite" // Need to avoid "Got error when connect database, the error is 'sql: unknown driver "sqlite3" (forgotten import?)'"
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	networkutilpkg "github.com/qb0C80aE/clay/util/network"
	"github.com/serenize/snaker"
)

func defaultTableNameHandler(db *gorm.DB, defaultTableName string) string {
	if db.Value == nil {
		return defaultTableName
	}

	model, err := extension.GetRegisteredModelByContainer(db.Value)
	if err != nil {
		if len(defaultTableName) > 0 {
			return defaultTableName
		}

		logging.Logger().Critical(err.Error())
		panic(err)
	}

	resourceName, err := model.GetResourceName(model)
	if err != nil {
		logging.Logger().Critical(err.Error())
		panic(err)
	}

	return resourceName
}

// IsIpv4AddressIncluding checks if left includes right network address, used in sqlite3
func IsIpv4AddressIncluding(baseCIDR, targetCIDR string) bool {
	base, _ := networkutilpkg.GetUtility().ParseCIDR(baseCIDR)
	target, _ := networkutilpkg.GetUtility().ParseCIDR(targetCIDR)
	return base.IsIncluding(target)
}

// SetupCustomDBFunctions setups custom sql functions
func SetupCustomDBFunctions() {
	sql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("is_ipv4_address_including", IsIpv4AddressIncluding, true)
		},
	})
}

// Connect connects to its database and returns the instance
func Connect(dbMode string) (*gorm.DB, error) {
	environmentalVariableSet := extension.GetCurrentEnvironmentalVariableSet()
	defaultEnvironmentalVariableSet := extension.GetDefaultEnvironmentalVariableSet()
	var dbPath string
	switch dbMode {
	case "memory":
		dbPath = ":memory:"
	case "", "file":
		if environmentalVariableSet.GetClayDBFilePath() != "" {
			dbPath = environmentalVariableSet.GetClayDBFilePath()
		} else {
			dbPath = defaultEnvironmentalVariableSet.GetClayDBFilePath()
		}
	default:
		logging.Logger().Criticalf("invalid db mode '%s'", dbMode)
		return nil, fmt.Errorf("invalid mode'%s'", dbMode)
	}

	// to avoid fallback because of unsupported dialect
	dialect, ok := gorm.GetDialect("sqlite3")
	if !ok {
		logging.Logger().Critical("cannot get the dialect sqlite3")
		return nil, errors.New("cannot get the dialect sqlite3")
	}
	gorm.RegisterDialect("sqlite3_custom", dialect)

	db, err := gorm.Open("sqlite3_custom", dbPath)

	if err != nil {
		logging.Logger().Criticalf("got an error when connect to the database, the error is '%v'", err)
		return nil, fmt.Errorf("got an error when connect to the database, the error is '%v'", err)
	}

	if gin.IsDebugging() {
		db.LogMode(true)
	}

	if err := db.Exec("pragma foreign_keys = on").Error; err != nil {
		logging.Logger().Critical(err)
		return nil, err
	}

	gorm.DefaultTableNameHandler = defaultTableNameHandler

	return db, nil
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
