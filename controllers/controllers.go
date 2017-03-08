package controllers

import (
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/helper"

	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extension"
	"io/ioutil"
	"net/http"
	"reflect"
)

func HookSubmodules() {
}

type BaseController struct {
	ResourceName string
	Model        interface{}
	Logic        extension.Logic
	Outputter    Outputter
}

type Outputter interface {
	OutputError(c *gin.Context, code int, err error)
	OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{})
	OutputGetMulti(c *gin.Context, code int, result []interface{}, fields map[string]interface{})
	OutputCreate(c *gin.Context, code int, result interface{})
	OutputUpdate(c *gin.Context, code int, result interface{})
	OutputDelete(c *gin.Context, code int)
	OutputPatch(c *gin.Context, code int, result interface{})
	OutputOptions(c *gin.Context, code int)
}

func (this *BaseController) bind(c *gin.Context, container interface{}) error {
	if err := c.Bind(container); err != nil {
		return err
	}
	if c.Request.ParseMultipartForm(1024*1024) == nil {
		files := c.Request.MultipartForm.File
		for name := range files {
			buffer := &bytes.Buffer{}
			file, _, err := c.Request.FormFile(name)
			if err != nil {
				return err
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			_, err = buffer.Write(data)
			if err != nil {
				return err
			}

			vs := reflect.ValueOf(container)
			for vs.Kind() == reflect.Ptr {
				vs = vs.Elem()
			}
			if !vs.IsValid() {
				return errors.New("invalid container.")
			}
			if !vs.CanInterface() {
				return errors.New("invalid container.")
			}
			value := vs.Interface()

			t := reflect.TypeOf(value)
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				jsonTag := field.Tag.Get("json")
				formTag := field.Tag.Get("form")
				if (jsonTag == name || formTag == name) && (field.Type.Kind() == reflect.String) {
					vs.FieldByName(field.Name).SetString(buffer.String())
					break
				}
			}

		}
	}
	return nil
}

func (this *BaseController) GetResourceName() string {
	return this.ResourceName
}

func (this *BaseController) OutputError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
}

func (this *BaseController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	if fields == nil {
		c.JSON(code, result)
	} else {
		fieldMap, err := helper.FieldToMap(result, fields)
		if err != nil {
			this.OutputError(c, http.StatusBadRequest, err)
			return
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(code, fieldMap)
		} else {
			c.JSON(code, fieldMap)
		}
	}
}

func (this *BaseController) OutputGetMulti(c *gin.Context, code int, result []interface{}, fields map[string]interface{}) {
	if fields == nil {
		c.JSON(code, result)
	} else {
		if _, ok := c.GetQuery("stream"); ok {
			enc := json.NewEncoder(c.Writer)
			c.Status(code)

			for _, item := range result {
				fieldMap, err := helper.FieldToMap(item, fields)

				if err != nil {
					this.OutputError(c, http.StatusBadRequest, err)
					return
				}

				if err := enc.Encode(fieldMap); err != nil {
					this.OutputError(c, http.StatusBadRequest, err)
					return
				}
			}
		} else {
			fieldMaps := []map[string]interface{}{}

			for _, item := range result {
				fieldMap, err := helper.FieldToMap(item, fields)

				if err != nil {
					this.OutputError(c, http.StatusBadRequest, err)
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

func (this *BaseController) OutputCreate(c *gin.Context, code int, result interface{}) {
	this.OutputGetSingle(c, code, result, nil)
}

func (this *BaseController) OutputUpdate(c *gin.Context, code int, result interface{}) {
	this.OutputGetSingle(c, code, result, nil)
}

func (this *BaseController) OutputDelete(c *gin.Context, code int) {
	c.Writer.WriteHeader(code)
}

func (this *BaseController) OutputPatch(c *gin.Context, code int, result interface{}) {
	this.OutputGetSingle(c, code, result, nil)
}

func (this *BaseController) OutputOptions(c *gin.Context, code int) {
	this.OutputDelete(c, code)
}

func (this *BaseController) GetSingle(c *gin.Context) {

	id := c.Params.ByName("id")
	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(c.Query("preloads"), db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(this.Model, fields)

	result, err := this.Logic.GetSingle(db, id, queryFields)
	if err != nil {
		this.OutputError(c, http.StatusNotFound, err)
		return
	}

	this.OutputGetSingle(c, http.StatusOK, result, fields)
}

func (this *BaseController) GetMulti(c *gin.Context) {

	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(c.Query("preloads"), db)
	db = dbpkg.SortRecords(c.Query("sort"), db)
	db = dbpkg.FilterFields(c, this.Model, db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(this.Model, fields)

	result, err := this.Logic.GetMulti(db, queryFields)
	if err != nil {
		this.OutputError(c, http.StatusBadRequest, err)
		return
	}

	this.Outputter.OutputGetMulti(c, http.StatusOK, result, fields)
}

func (this *BaseController) Create(c *gin.Context) {
	vs := reflect.ValueOf(this.Model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		this.OutputError(c, http.StatusBadRequest, errors.New("Invalid model."))
		return
	}
	if !vs.CanInterface() {
		this.OutputError(c, http.StatusBadRequest, errors.New("Invalid model."))
		return
	}
	container := reflect.New(reflect.TypeOf(vs.Interface())).Interface()

	if err := this.bind(c, container); err != nil {
		this.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := this.Logic.Create(db, container)
	if err != nil {
		db.Rollback()
		this.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	this.Outputter.OutputCreate(c, http.StatusCreated, result)
}

func (this *BaseController) Update(c *gin.Context) {
	vs := reflect.ValueOf(this.Model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		this.OutputError(c, http.StatusBadRequest, errors.New("Invalid model."))
		return
	}
	if !vs.CanInterface() {
		this.OutputError(c, http.StatusBadRequest, errors.New("Invalid model."))
		return
	}
	container := reflect.New(reflect.TypeOf(vs.Interface())).Interface()

	id := c.Params.ByName("id")

	if err := this.bind(c, container); err != nil {
		this.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := this.Logic.Update(db, id, container)
	if err != nil {
		db.Rollback()
		this.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	this.Outputter.OutputUpdate(c, http.StatusOK, result)
}

func (this *BaseController) Delete(c *gin.Context) {
	id := c.Params.ByName("id")

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	err := this.Logic.Delete(db, id)
	if err != nil {
		db.Rollback()
		this.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	this.Outputter.OutputDelete(c, http.StatusNoContent)
}

func (this *BaseController) Patch(c *gin.Context) {
	id := c.Params.ByName("id")

	db := dbpkg.DBInstance(c)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(this.Model, fields)

	db = db.Begin()
	result, err := this.Logic.Patch(db, id, queryFields)
	if err != nil {
		db.Rollback()
		this.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	this.Outputter.OutputPatch(c, http.StatusOK, result)
}

func (this *BaseController) Options(c *gin.Context) {
	db := dbpkg.DBInstance(c)

	db = db.Begin()
	err := this.Logic.Options(db)
	if err != nil {
		db.Rollback()
		this.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	this.Outputter.OutputOptions(c, http.StatusNoContent)
}
