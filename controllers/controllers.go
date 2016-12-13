package controllers

import (
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/helper"
	"github.com/qb0C80aE/clay/version"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func processSingleGet(c *gin.Context, model interface{}, actualLogic func(*gorm.DB, string, string) (interface{}, error)) {
	ver, err := version.New(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Params.ByName("id")
	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(c.Query("preloads"), db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(model, fields)

	result, err := actualLogic(db, id, queryFields)
	if err != nil {
		content := gin.H{"error": "item with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	fieldMap, err := helper.FieldToMap(result, fields)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("pretty"); ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

func processMultiGet(c *gin.Context, model interface{}, actualLogic func(*gorm.DB, string) ([]interface{}, error)) {
	ver, err := version.New(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(c.Query("preloads"), db)
	db = dbpkg.SortRecords(c.Query("sort"), db)
	db = dbpkg.FilterFields(c, model, db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(model, fields)

	result, err := actualLogic(db, queryFields)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, item := range result {
			fieldMap, err := helper.FieldToMap(item, fields)

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if err := enc.Encode(fieldMap); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
	} else {
		fieldMaps := []map[string]interface{}{}

		for _, item := range result {
			fieldMap, err := helper.FieldToMap(item, fields)

			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMaps = append(fieldMaps, fieldMap)
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(200, fieldMaps)
		} else {
			c.JSON(200, fieldMaps)
		}
	}
}

func processCreate(c *gin.Context, container interface{}, model interface{}, actualLogic func(*gorm.DB, interface{}) (interface{}, error)) {
	ver, err := version.New(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := c.Bind(container); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := actualLogic(db, container)
	if err != nil {
		db.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(201, result)
}

func processUpdate(c *gin.Context, container interface{}, model interface{}, actualLogic func(*gorm.DB, string, interface{}) (interface{}, error)) {
	ver, err := version.New(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Params.ByName("id")

	if err := c.Bind(container); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := actualLogic(db, id, container)
	if err != nil {
		db.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.JSON(200, result)
}

func processDelete(c *gin.Context, model interface{}, actualLogic func(*gorm.DB, string) error) {
	ver, err := version.New(c)

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Params.ByName("id")

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	err = actualLogic(db, id)
	if err != nil {
		db.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	c.Writer.WriteHeader(http.StatusNoContent)
}
