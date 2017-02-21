package controllers

import (
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/helper"

	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"io/ioutil"
	"bytes"
	"reflect"
)

func bind(c *gin.Context, container interface{}) error {
	if err := c.Bind(container); err != nil {
		return err
	}
	if c.Request.ParseMultipartForm(1024 * 1024) == nil {
		files := c.Request.MultipartForm.File
		for name, _ := range (files) {
			buffer := &bytes.Buffer{}
			file, _, err := c.Request.FormFile(name)
			if err != nil {
				return err
			}
			data, _ := ioutil.ReadAll(file)
			if err != nil {
				return err
			}
			buffer.Write(data)

			vs := reflect.ValueOf(container)
			for vs.Kind() == reflect.Ptr {
				vs = vs.Elem()
			}
			if !vs.IsValid() {
				return errors.New("invalid container.")
			}
			if !vs.CanInterface() {
				return nil
			}
			value := vs.Interface()

			t := reflect.TypeOf(value)
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				jsonTag := field.Tag.Get("json")
				formTag := field.Tag.Get("form")
				if (jsonTag == name || formTag != name) && (field.Type.Kind() == reflect.String) {
					vs.FieldByName(field.Name).SetString(buffer.String())
					break
				}
			}



		}
	}
	return nil
}

func OutputJsonError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
}

func OutputSingleJsonResult(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	if fields == nil {
		c.JSON(code, result)
	} else {
		fieldMap, err := helper.FieldToMap(result, fields)
		if err != nil {
			OutputJsonError(c, http.StatusBadRequest, err)
			return
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(code, fieldMap)
		} else {
			c.JSON(code, fieldMap)
		}
	}
}

func OutputMultiJsonResult(c *gin.Context, code int, result []interface{}, fields map[string]interface{}) {
	if fields == nil {
		c.JSON(code, result)
	} else {
		if _, ok := c.GetQuery("stream"); ok {
			enc := json.NewEncoder(c.Writer)
			c.Status(code)

			for _, item := range result {
				fieldMap, err := helper.FieldToMap(item, fields)

				if err != nil {
					OutputJsonError(c, http.StatusBadRequest, err)
					return
				}

				if err := enc.Encode(fieldMap); err != nil {
					OutputJsonError(c, http.StatusBadRequest, err)
					return
				}
			}
		} else {
			fieldMaps := []map[string]interface{}{}

			for _, item := range result {
				fieldMap, err := helper.FieldToMap(item, fields)

				if err != nil {
					OutputJsonError(c, http.StatusBadRequest, err)
					return
				}

				fieldMaps = append(fieldMaps, fieldMap)
			}

			if _, ok := c.GetQuery("pretty"); ok {
				c.IndentedJSON(code, fieldMaps)
			} else {
				c.JSON(code, fieldMaps)
			}
		}
	}
}

func OutputTextResult(c *gin.Context, code int, result interface{}, _ map[string]interface{}) {
	text := result.(string)
	c.String(code, text)
}

func OutputNothing(c *gin.Context, code int, _ interface{}, _ map[string]interface{}) {
	c.Writer.WriteHeader(code)
}

func processSingleGet(c *gin.Context,
	model interface{},
	actualLogic func(*gorm.DB, string, string) (interface{}, error),
	errorOutputFunction func(*gin.Context, int, error),
	resultOutputFunction func(*gin.Context, int, interface{}, map[string]interface{})) {

	id := c.Params.ByName("id")
	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(c.Query("preloads"), db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(model, fields)

	result, err := actualLogic(db, id, queryFields)
	if err != nil {
		errorOutputFunction(c, http.StatusNotFound, errors.New("item with id#"+id+" not found"))
		return
	}

	resultOutputFunction(c, http.StatusOK, result, fields)
}

func processMultiGet(c *gin.Context,
	model interface{},
	actualLogic func(*gorm.DB, string) ([]interface{}, error),
	errorOutputFunction func(*gin.Context, int, error),
	resultOutputFunction func(*gin.Context, int, []interface{}, map[string]interface{})) {

	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(c.Query("preloads"), db)
	db = dbpkg.SortRecords(c.Query("sort"), db)
	db = dbpkg.FilterFields(c, model, db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(model, fields)

	result, err := actualLogic(db, queryFields)
	if err != nil {
		errorOutputFunction(c, http.StatusBadRequest, err)
		return
	}

	resultOutputFunction(c, http.StatusOK, result, fields)
}

func processCreate(c *gin.Context,
	container interface{},
	actualLogic func(*gorm.DB, interface{}) (interface{}, error),
	errorOutputFunction func(*gin.Context, int, error),
	resultOutputFunction func(*gin.Context, int, interface{}, map[string]interface{})) {
	if err := bind(c, container); err != nil {
		errorOutputFunction(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := actualLogic(db, container)
	if err != nil {
		db.Rollback()
		errorOutputFunction(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	resultOutputFunction(c, http.StatusCreated, result, nil)
}

func processUpdate(c *gin.Context,
	container interface{},
	actualLogic func(*gorm.DB, string, interface{}) (interface{}, error),
	errorOutputFunction func(*gin.Context, int, error),
	resultOutputFunction func(*gin.Context, int, interface{}, map[string]interface{})) {
	id := c.Params.ByName("id")

	if err := bind(c, container); err != nil {
		errorOutputFunction(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := actualLogic(db, id, container)
	if err != nil {
		db.Rollback()
		errorOutputFunction(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	resultOutputFunction(c, http.StatusOK, result, nil)
}

func processDelete(c *gin.Context,
	actualLogic func(*gorm.DB, string) error,
	errorOutputFunction func(*gin.Context, int, error),
	resultOutputFunction func(*gin.Context, int, interface{}, map[string]interface{})) {
	id := c.Params.ByName("id")

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	err := actualLogic(db, id)
	if err != nil {
		db.Rollback()
		errorOutputFunction(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	resultOutputFunction(c, http.StatusNoContent, nil, nil)
}
