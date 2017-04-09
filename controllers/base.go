package controllers

import (
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/helper"

	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/qb0C80aE/clay/extensions"
	"io/ioutil"
	"net/http"
	"reflect"
)

type BaseController struct {
	resourceName string
	model        interface{}
	logic        extensions.Logic
	outputter    extensions.Outputter
}

func NewBaseController(resourceName string, model interface{}, logic extensions.Logic) *BaseController {
	controller := &BaseController{
		resourceName: resourceName,
		model:        model,
		logic:        logic,
	}
	return controller
}

func (controller *BaseController) bind(c *gin.Context, container interface{}) error {
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
				return errors.New("invalid container")
			}
			if !vs.CanInterface() {
				return errors.New("invalid container")
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

func (controller *BaseController) ResourceName() string {
	return controller.resourceName
}

func (controller *BaseController) OutputError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
}

func (controller *BaseController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	if fields == nil {
		c.JSON(code, result)
	} else {
		fieldMap, err := helper.FieldToMap(result, fields)
		if err != nil {
			controller.OutputError(c, http.StatusBadRequest, err)
			return
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(code, fieldMap)
		} else {
			c.JSON(code, fieldMap)
		}
	}
}

func (controller *BaseController) OutputGetMulti(c *gin.Context, code int, result []interface{}, fields map[string]interface{}) {
	if fields == nil {
		c.JSON(code, result)
	} else {
		if _, ok := c.GetQuery("stream"); ok {
			enc := json.NewEncoder(c.Writer)
			c.Status(code)

			for _, item := range result {
				fieldMap, err := helper.FieldToMap(item, fields)

				if err != nil {
					controller.OutputError(c, http.StatusBadRequest, err)
					return
				}

				if err := enc.Encode(fieldMap); err != nil {
					controller.OutputError(c, http.StatusBadRequest, err)
					return
				}
			}
		} else {
			fieldMaps := []map[string]interface{}{}

			for _, item := range result {
				fieldMap, err := helper.FieldToMap(item, fields)

				if err != nil {
					controller.OutputError(c, http.StatusBadRequest, err)
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

func (controller *BaseController) OutputCreate(c *gin.Context, code int, result interface{}) {
	controller.OutputGetSingle(c, code, result, nil)
}

func (controller *BaseController) OutputUpdate(c *gin.Context, code int, result interface{}) {
	controller.OutputGetSingle(c, code, result, nil)
}

func (controller *BaseController) OutputDelete(c *gin.Context, code int) {
	c.Writer.WriteHeader(code)
}

func (controller *BaseController) OutputPatch(c *gin.Context, code int, result interface{}) {
	controller.OutputGetSingle(c, code, result, nil)
}

func (controller *BaseController) OutputOptions(c *gin.Context, code int) {
	controller.OutputDelete(c, code)
}

func (controller *BaseController) GetSingle(c *gin.Context) {
	id := c.Params.ByName("id")
	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(c.Query("preloads"), db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(controller.model, fields)

	result, err := controller.logic.GetSingle(db, id, queryFields)
	if err != nil {
		controller.OutputError(c, http.StatusNotFound, err)
		return
	}

	controller.OutputGetSingle(c, http.StatusOK, result, fields)
}

func (controller *BaseController) GetMulti(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	db = dbpkg.SetPreloads(c.Query("preloads"), db)
	db = dbpkg.SortRecords(c.Query("sort"), db)
	db = dbpkg.FilterFields(c, controller.model, db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(controller.model, fields)

	result, err := controller.logic.GetMulti(db, queryFields)
	if err != nil {
		controller.OutputError(c, http.StatusBadRequest, err)
		return
	}

	controller.outputter.OutputGetMulti(c, http.StatusOK, result, fields)
}

func (controller *BaseController) Create(c *gin.Context) {
	vs := reflect.ValueOf(controller.model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		controller.OutputError(c, http.StatusBadRequest, errors.New("Invalid model."))
		return
	}
	if !vs.CanInterface() {
		controller.OutputError(c, http.StatusBadRequest, errors.New("Invalid model."))
		return
	}
	container := reflect.New(reflect.TypeOf(vs.Interface())).Interface()

	if err := controller.bind(c, container); err != nil {
		controller.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := controller.logic.Create(db, container)
	if err != nil {
		db.Rollback()
		controller.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	controller.outputter.OutputCreate(c, http.StatusCreated, result)
}

func (controller *BaseController) Update(c *gin.Context) {
	vs := reflect.ValueOf(controller.model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		controller.OutputError(c, http.StatusBadRequest, errors.New("Invalid model."))
		return
	}
	if !vs.CanInterface() {
		controller.OutputError(c, http.StatusBadRequest, errors.New("Invalid model."))
		return
	}
	container := reflect.New(reflect.TypeOf(vs.Interface())).Interface()

	id := c.Params.ByName("id")

	if err := controller.bind(c, container); err != nil {
		controller.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := controller.logic.Update(db, id, container)
	if err != nil {
		db.Rollback()
		controller.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	controller.outputter.OutputUpdate(c, http.StatusOK, result)
}

func (controller *BaseController) Delete(c *gin.Context) {
	id := c.Params.ByName("id")

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	err := controller.logic.Delete(db, id)
	if err != nil {
		db.Rollback()
		controller.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	controller.outputter.OutputDelete(c, http.StatusNoContent)
}

func (controller *BaseController) Patch(c *gin.Context) {
	id := c.Params.ByName("id")

	db := dbpkg.DBInstance(c)

	db = db.Begin()
	result, err := controller.logic.Patch(db, id)
	if err != nil {
		db.Rollback()
		controller.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	controller.outputter.OutputPatch(c, http.StatusOK, result)
}

func (controller *BaseController) Options(c *gin.Context) {
	db := dbpkg.DBInstance(c)

	db = db.Begin()
	err := controller.logic.Options(db)
	if err != nil {
		db.Rollback()
		controller.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db.Commit()

	controller.outputter.OutputOptions(c, http.StatusNoContent)
}

func (controller *BaseController) SetOutputter(outputter extensions.Outputter) {
	controller.outputter = outputter
}
