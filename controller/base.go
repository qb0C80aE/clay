package controller

import (
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/helper"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"github.com/qb0C80aE/clay/version"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

// BaseController is the base class that all controller classes inherit
type BaseController struct {
	actualController extension.Controller
	model            extension.Model
	binder           extension.Binder
	outputHandler    extension.OutputHandler
	queryCustomizer  extension.QueryCustomizer
}

// CreateController creates a new instance of actual controller with BaseController
func CreateController(actualController extension.Controller, model extension.Model) extension.Controller {
	actualControllerValue := reflect.ValueOf(actualController).Elem()
	baseController := BaseController{
		actualController: actualController,
		model:            model,
		binder:           actualController.(extension.Binder),
		outputHandler:    actualController.(extension.OutputHandler),
		queryCustomizer:  actualController.(extension.QueryCustomizer),
	}
	baseControllerValue := reflect.ValueOf(baseController)
	actualControllerValue.FieldByName("BaseController").Set(baseControllerValue)
	return actualController
}

// Bind binds input data to model instance
func (receiver *BaseController) Bind(c *gin.Context, container interface{}) error {
	if err := c.Bind(container); err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}
	if c.Request.ParseMultipartForm(1024*1024) == nil {
		files := c.Request.MultipartForm.File
		for name := range files {
			buffer := &bytes.Buffer{}
			file, _, err := c.Request.FormFile(name)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return err
			}

			data, err := ioutil.ReadAll(file)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return err
			}

			_, err = buffer.Write(data)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return err
			}

			vs := reflect.ValueOf(container)
			for vs.Kind() == reflect.Ptr {
				vs = vs.Elem()
			}
			if !vs.IsValid() {
				logging.Logger().Debug("invalid model")
				return errors.New("invalid model")
			}
			if !vs.CanInterface() {
				logging.Logger().Debug("invalid model")
				return errors.New("invalid model")
			}
			value := vs.Interface()

			t := reflect.TypeOf(value)
			for i := 0; i < t.NumField(); i++ {
				field := t.Field(i)
				jsonTag := field.Tag.Get("json")
				formTag := field.Tag.Get("form")
				if jsonTag == name || formTag == name {
					if field.Type.Kind() == reflect.String {
						vs.FieldByName(field.Name).SetString(buffer.String())
						break
					} else if field.Type.Kind() == reflect.Slice {
						vs.FieldByName(field.Name).SetBytes(buffer.Bytes())
						break
					} else {
						logging.Logger().Debug("invalid field definition, the field must be string or slice of bytes")
						return errors.New("invalid field definition, the field must be string or slice of bytes")
					}
				}
			}

		}
	}
	return nil
}

// GetModel returns its model
func (receiver *BaseController) GetModel() extension.Model {
	return receiver.model
}

// GetResourceName returns its resource name in REST
func (receiver *BaseController) GetResourceName() string {
	return extension.GetAssociateResourceNameWithModel(receiver.model)
}

// GetResourceSingleURL builds a resource url what represents a single resource based on the argument
func (receiver *BaseController) GetResourceSingleURL() string {
	return fmt.Sprintf("%s/:id", receiver.GetResourceName())
}

// GetResourceMultiURL builds a resource url what represents multi resources based on the argument
func (receiver *BaseController) GetResourceMultiURL() string {
	return fmt.Sprintf("%s", receiver.GetResourceName())
}

// OutputError handles an error output
func (receiver *BaseController) OutputError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
}

// OutputGetSingle corresponds HTTP GET message and handles the output of a single result from logic classes
func (receiver *BaseController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	_, allFieldExists := fields["*"]
	if (fields == nil) || ((len(fields) == 1) && allFieldExists) {
		c.JSON(code, result)
	} else {
		fieldMap, err := helper.FieldToMap(result, fields)
		if err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(code, fieldMap)
		} else {
			c.JSON(code, fieldMap)
		}
	}
}

// OutputGetMulti corresponds HTTP GET message and handles the output of multiple result from logic classes
func (receiver *BaseController) OutputGetMulti(c *gin.Context, code int, result interface{}, total int, fields map[string]interface{}) {
	c.Header("Total", strconv.Itoa(total))
	_, allFieldExists := fields["*"]
	if (fields == nil) || ((len(fields) == 1) && allFieldExists) {
		c.JSON(code, result)
	} else {
		v := reflect.ValueOf(result)

		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			logging.Logger().Debug("given argument is neither a slice nor an array")
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("given argument is neither a slice nor an array"))
			return
		}

		size := v.Len()

		if _, ok := c.GetQuery("stream"); ok {
			enc := json.NewEncoder(c.Writer)
			c.Status(code)

			for i := 0; i < size; i++ {
				item := v.Index(i)

				if !item.CanInterface() {
					logging.Logger().Debugf("the original item indexed %d in given slice cannot interface", i)
					receiver.outputHandler.OutputError(c, http.StatusBadRequest, fmt.Errorf("the original item indexed %d in given slice cannot interface", i))
					return
				}

				fieldMap, err := helper.FieldToMap(item.Interface(), fields)

				if err != nil {
					logging.Logger().Debug(err.Error())
					receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
					return
				}

				if err := enc.Encode(fieldMap); err != nil {
					logging.Logger().Debug(err.Error())
					receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
					return
				}
			}
		} else {
			fieldMaps := []map[string]interface{}{}

			for i := 0; i < size; i++ {
				item := v.Index(i)

				if !item.CanInterface() {
					logging.Logger().Debugf("the original item indexed %d in given slice cannot interface", i)
					receiver.outputHandler.OutputError(c, http.StatusBadRequest, fmt.Errorf("the original item indexed %d in given slice cannot interface", i))
					return
				}

				fieldMap, err := helper.FieldToMap(item.Interface(), fields)

				if err != nil {
					logging.Logger().Debug(err.Error())
					receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
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

// OutputCreate corresponds HTTP POST message and handles the output of a single result from logic classes
func (receiver *BaseController) OutputCreate(c *gin.Context, code int, result interface{}) {
	receiver.outputHandler.OutputGetSingle(c, code, result, nil)
}

// OutputUpdate corresponds HTTP PUT message and handles the output of a single result from logic classes
func (receiver *BaseController) OutputUpdate(c *gin.Context, code int, result interface{}) {
	receiver.outputHandler.OutputGetSingle(c, code, result, nil)
}

// OutputDelete corresponds HTTP DELETE message and handles the code result from logic classes
func (receiver *BaseController) OutputDelete(c *gin.Context, code int) {
	c.Writer.WriteHeader(code)
}

// OutputPatch corresponds HTTP PATCH message and handles the output of a single result from logic classes
func (receiver *BaseController) OutputPatch(c *gin.Context, code int, result interface{}) {
	receiver.outputHandler.OutputGetSingle(c, code, result, nil)
}

// OutputGetOptions corresponds HTTP OPTIONS message and handles the code result from logic classes, as well as OutputDelete
func (receiver *BaseController) OutputGetOptions(c *gin.Context, code int) {
	receiver.outputHandler.OutputDelete(c, code)
}

// GetQueries returns query parameters
func (receiver *BaseController) GetQueries(c *gin.Context) url.Values {
	return c.Request.URL.Query()
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (receiver *BaseController) GetSingle(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)
	parameter, err := dbpkg.NewParameter(receiver.queryCustomizer.GetQueries(c))
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db = parameter.SetPreloads(db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(receiver.model, fields)

	result, err := receiver.model.GetSingle(db, c.Params, c.Request.URL.Query(), queryFields)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusNotFound, err)
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputGetSingle(c, http.StatusOK, result, fields)
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (receiver *BaseController) GetMulti(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)
	parameter, err := dbpkg.NewParameter(receiver.queryCustomizer.GetQueries(c))
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db, err = parameter.Paginate(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(receiver.model, fields)

	result, err := receiver.model.GetMulti(db, c.Params, c.Request.URL.Query(), queryFields)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	// reset conditions
	db = db.New()

	total, err := receiver.model.GetTotal(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputGetMulti(c, http.StatusOK, result, total, fields)
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *BaseController) Create(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	container := receiver.model.New()

	if err := receiver.binder.Bind(c, container); err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	model := model.ConvertContainerToModel(container)

	db := dbpkg.Instance(c)

	tx := db.Begin()
	result, err := receiver.model.Create(tx, c.Params, c.Request.URL.Query(), model)
	if err != nil {
		tx.Rollback()
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputCreate(c, http.StatusCreated, result)
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (receiver *BaseController) Update(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	container := receiver.model.New()

	if err := receiver.binder.Bind(c, container); err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	model := model.ConvertContainerToModel(container)

	db := dbpkg.Instance(c)

	tx := db.Begin()
	result, err := receiver.model.Update(tx, c.Params, c.Request.URL.Query(), model)
	if err != nil {
		tx.Rollback()
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputUpdate(c, http.StatusOK, result)
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (receiver *BaseController) Delete(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	err = receiver.model.Delete(tx, c.Params, c.Request.URL.Query())
	if err != nil {
		tx.Rollback()
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputDelete(c, http.StatusNoContent)
}

// Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
func (receiver *BaseController) Patch(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	model := receiver.model.New()

	db := dbpkg.Instance(c)

	tx := db.Begin()
	result, err := receiver.model.Patch(tx, c.Params, c.Request.URL.Query(), model)
	if err != nil {
		tx.Rollback()
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputUpdate(c, http.StatusOK, result)
}

// GetOptions corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
func (receiver *BaseController) GetOptions(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	err = receiver.model.GetOptions(tx, c.Params, c.Request.URL.Query())
	if err != nil {
		tx.Rollback()
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	receiver.outputHandler.OutputGetOptions(c, http.StatusNoContent)
}

// DoAfterDBMigration execute initialization process after DB migration
func (receiver *BaseController) DoAfterDBMigration(db *gorm.DB) error {
	return nil
}

// DoBeforeRouterSetup execute initialization process before Router initialization
func (receiver *BaseController) DoBeforeRouterSetup(r *gin.Engine) error {
	return nil
}

// DoAfterRouterSetup execute initialization process after Router initialization
func (receiver *BaseController) DoAfterRouterSetup(r *gin.Engine) error {
	return nil
}
