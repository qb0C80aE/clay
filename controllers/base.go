package controllers

import (
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/helper"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/version"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

// BaseController is the base class that all controller classes inherit
type BaseController struct {
	model           interface{}
	logic           extensions.Logic
	outputter       extensions.Outputter
	queryCustomizer extensions.QueryCustomizer
}

// NewBaseController creates a new instance of BaseController
func NewBaseController(model interface{}, logic extensions.Logic) *BaseController {
	controller := &BaseController{
		model: model,
		logic: logic,
	}
	controller.outputter = controller
	controller.queryCustomizer = controller
	return controller
}

func (controller *BaseController) deleteMarkedItemsInSlices(db *gorm.DB, data interface{}) error {
	valueOfData := reflect.ValueOf(data)
	for valueOfData.Kind() == reflect.Ptr {
		valueOfData = valueOfData.Elem()
	}

	typeOfData := valueOfData.Type()

	for i := 0; i < typeOfData.NumField(); i++ {
		structField := typeOfData.Field(i)
		fieldValue := valueOfData.FieldByName(structField.Name)

		for fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}

		if fieldValue.Kind() == reflect.Slice {
			processed := reflect.New(fieldValue.Type()).Elem()
			for j := 0; j < fieldValue.Len(); j++ {
				itemValue := fieldValue.Index(j)
				for itemValue.Kind() == reflect.Ptr {
					itemValue = itemValue.Elem()
				}

				toBeDeleted := false
				toBeDeletedFieldValue := itemValue.FieldByName("ToBeDeleted")
				if toBeDeletedFieldValue.IsValid() {
					toBeDeleted = toBeDeletedFieldValue.Bool()
				}

				if toBeDeleted {
					if err := controller.deleteMarkedItemsInSlices(db, fieldValue.Index(j).Interface()); err != nil {
						return err
					}

					logic, err := extensions.RegisteredLogic(itemValue.Interface())
					if err != nil {
						return err
					}

					modelKey, err := extensions.RegisteredModelKey(itemValue.Interface())
					if err != nil {
						return err
					}

					keyFieldValue := itemValue.FieldByName(modelKey.KeyField)
					keyParameterValue := ""

					switch keyFieldValue.Kind() {
					case reflect.String:
						keyParameterValue = keyFieldValue.Interface().(string)
					case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
						keyParameterValue = strconv.Itoa(int(keyFieldValue.Int()))
					case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
						keyParameterValue = strconv.Itoa(int(keyFieldValue.Int()))
					default:
						return fmt.Errorf("the field %s does not exist, or is neither int nor string", modelKey.KeyField)
					}

					parameters := gin.Params{
						{
							Key:   modelKey.KeyParameter,
							Value: keyParameterValue,
						},
					}
					if err := logic.Delete(db, parameters, nil); err != nil {
						return err
					}
				} else {
					processed = reflect.Append(processed, itemValue.Addr())
				}
			}
			fieldValue.Set(processed)
		} else if fieldValue.Kind() == reflect.Struct {
			if err := controller.deleteMarkedItemsInSlices(db, valueOfData.FieldByName(structField.Name).Interface()); err != nil {
				return err
			}
		}
	}
	return nil
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
				if jsonTag == name || formTag == name {
					if field.Type.Kind() == reflect.String {
						vs.FieldByName(field.Name).SetString(buffer.String())
						break
					} else if field.Type.Kind() == reflect.Slice {
						vs.FieldByName(field.Name).SetBytes(buffer.Bytes())
						break
					} else {
						return errors.New("invalid field definition, the field must be string or slice of bytes")
					}
				}
			}

		}
	}
	return nil
}

// ResourceName returns its resource name in REST
func (controller *BaseController) ResourceName() string {
	return extensions.RegisteredResourceName(controller.model)
}

// ResourceSingleURL builds a resource url what represents a single resource based on the argument
func (controller *BaseController) ResourceSingleURL() string {
	return fmt.Sprintf("%s/:id", controller.ResourceName())
}

// ResourceMultiURL builds a resource url what represents multi resources based on the argument
func (controller *BaseController) ResourceMultiURL() string {
	return fmt.Sprintf("%s", controller.ResourceName())
}

// OutputError handles an error output
func (controller *BaseController) OutputError(c *gin.Context, code int, err error) {
	c.JSON(code, gin.H{"error": err.Error()})
}

// OutputGetSingle corresponds HTTP GET message and handles the output of a single result from logic classes
func (controller *BaseController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	if fields == nil {
		c.JSON(code, result)
	} else {
		fieldMap, err := helper.FieldToMap(result, fields)
		if err != nil {
			controller.outputter.OutputError(c, http.StatusBadRequest, err)
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
func (controller *BaseController) OutputGetMulti(c *gin.Context, code int, result interface{}, total int, fields map[string]interface{}) {
	c.Header("Total", strconv.Itoa(total))
	if fields == nil {
		c.JSON(code, result)
	} else {
		v := reflect.ValueOf(result)

		if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
			controller.outputter.OutputError(c, http.StatusBadRequest, errors.New("given argument is neither a slice nor an array"))
			return
		}

		size := v.Len()

		if _, ok := c.GetQuery("stream"); ok {
			enc := json.NewEncoder(c.Writer)
			c.Status(code)

			for i := 0; i < size; i++ {
				item := v.Index(i)

				if !item.CanInterface() {
					controller.outputter.OutputError(c, http.StatusBadRequest, fmt.Errorf("the original item indexed %d in given slice cannot interface", i))
					return
				}

				fieldMap, err := helper.FieldToMap(item.Interface(), fields)

				if err != nil {
					controller.outputter.OutputError(c, http.StatusBadRequest, err)
					return
				}

				if err := enc.Encode(fieldMap); err != nil {
					controller.outputter.OutputError(c, http.StatusBadRequest, err)
					return
				}
			}
		} else {
			fieldMaps := []map[string]interface{}{}

			for i := 0; i < size; i++ {
				item := v.Index(i)

				if !item.CanInterface() {
					controller.outputter.OutputError(c, http.StatusBadRequest, fmt.Errorf("the original item indexed %d in given slice cannot interface", i))
					return
				}

				fieldMap, err := helper.FieldToMap(item.Interface(), fields)

				if err != nil {
					controller.outputter.OutputError(c, http.StatusBadRequest, err)
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
func (controller *BaseController) OutputCreate(c *gin.Context, code int, result interface{}) {
	controller.outputter.OutputGetSingle(c, code, result, nil)
}

// OutputUpdate corresponds HTTP PUT message and handles the output of a single result from logic classes
func (controller *BaseController) OutputUpdate(c *gin.Context, code int, result interface{}) {
	controller.outputter.OutputGetSingle(c, code, result, nil)
}

// OutputDelete corresponds HTTP DELETE message and handles the code result from logic classes
func (controller *BaseController) OutputDelete(c *gin.Context, code int) {
	c.Writer.WriteHeader(code)
}

// OutputPatch corresponds HTTP PATCH message and handles the output of a single result from logic classes
func (controller *BaseController) OutputPatch(c *gin.Context, code int, result interface{}) {
	controller.outputter.OutputGetSingle(c, code, result, nil)
}

// OutputOptions corresponds HTTP DELETE message and handles the code result from logic classes, as well as OutputDelete
func (controller *BaseController) OutputOptions(c *gin.Context, code int) {
	controller.outputter.OutputDelete(c, code)
}

// Queries returns query parameters
func (controller *BaseController) Queries(c *gin.Context) url.Values {
	return c.Request.URL.Query()
}

func (controller *BaseController) getSingle(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (result interface{}, err error) {
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	result, err = controller.logic.GetSingle(db, parameters, urlValues, queryFields)
	return result, err
}

func (controller *BaseController) getMulti(db *gorm.DB, parameters gin.Params, urlValues url.Values, queryFields string) (result interface{}, err error) {
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	result, err = controller.logic.GetMulti(db, parameters, urlValues, queryFields)
	return result, err
}

func (controller *BaseController) create(db *gorm.DB, parameters gin.Params, urlValues url.Values, data interface{}) (result interface{}, err error) {
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	result, err = controller.logic.Create(db, parameters, urlValues, data)
	return result, err
}

func (controller *BaseController) update(db *gorm.DB, parameters gin.Params, urlValues url.Values, data interface{}) (result interface{}, err error) {
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	if err := controller.deleteMarkedItemsInSlices(db, data); err != nil {
		return nil, err
	}

	result, err = controller.logic.Update(db, parameters, urlValues, data)
	return result, err
}

func (controller *BaseController) delete(db *gorm.DB, parameters gin.Params, urlValues url.Values) (err error) {
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	err = controller.logic.Delete(db, parameters, urlValues)
	return err
}

func (controller *BaseController) patch(db *gorm.DB, parameters gin.Params, urlValues url.Values) (result interface{}, err error) {
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	result, err = controller.logic.Patch(db, parameters, urlValues)
	return result, err
}

func (controller *BaseController) options(db *gorm.DB, parameters gin.Params, urlValues url.Values) (err error) {
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	err = controller.logic.Options(db, parameters, urlValues)
	return err
}

func (controller *BaseController) total(db *gorm.DB, model interface{}) (total int, err error) {
	defer func() {
		if recoverResult := recover(); recoverResult != nil {
			err = fmt.Errorf("%v", recoverResult)
		}
	}()

	db = db.New()
	total, err = controller.logic.Total(db, model)
	return total, err
}

// GetSingle corresponds HTTP GET message and handles a request for a single resource to get the information
func (controller *BaseController) GetSingle(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)
	parameter, err := dbpkg.NewParameter(controller.queryCustomizer.Queries(c))
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db = parameter.SetPreloads(db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(controller.model, fields)

	result, err := controller.getSingle(db, c.Params, c.Request.URL.Query(), queryFields)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusNotFound, err)
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	controller.outputter.OutputGetSingle(c, http.StatusOK, result, fields)
}

// GetMulti corresponds HTTP GET message and handles a request for multi resource to get the list of information
func (controller *BaseController) GetMulti(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)
	parameter, err := dbpkg.NewParameter(controller.queryCustomizer.Queries(c))
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db, err = parameter.Paginate(db)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(controller.model, fields)

	result, err := controller.getMulti(db, c.Params, c.Request.URL.Query(), queryFields)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	total, err := controller.total(db, controller.model)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	controller.outputter.OutputGetMulti(c, http.StatusOK, result, total, fields)
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (controller *BaseController) Create(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	vs := reflect.ValueOf(controller.model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		controller.outputter.OutputError(c, http.StatusBadRequest, errors.New("Invalid model"))
		return
	}
	if !vs.CanInterface() {
		controller.outputter.OutputError(c, http.StatusBadRequest, errors.New("Invalid model"))
		return
	}
	container := reflect.New(reflect.TypeOf(vs.Interface())).Interface()

	if err := controller.bind(c, container); err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	result, err := controller.create(tx, c.Params, c.Request.URL.Query(), container)
	if err != nil {
		tx.Rollback()
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	controller.outputter.OutputCreate(c, http.StatusCreated, result)
}

// Update corresponds HTTP PUT message and handles a request for a single resource to update the specific information
func (controller *BaseController) Update(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	vs := reflect.ValueOf(controller.model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		controller.outputter.OutputError(c, http.StatusBadRequest, errors.New("Invalid model"))
		return
	}
	if !vs.CanInterface() {
		controller.outputter.OutputError(c, http.StatusBadRequest, errors.New("Invalid model"))
		return
	}
	container := reflect.New(reflect.TypeOf(vs.Interface())).Interface()

	if err := controller.bind(c, container); err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	result, err := controller.update(tx, c.Params, c.Request.URL.Query(), container)
	if err != nil {
		tx.Rollback()
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	controller.outputter.OutputUpdate(c, http.StatusOK, result)
}

// Delete corresponds HTTP DELETE message and handles a request for a single resource to delete the specific information
func (controller *BaseController) Delete(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	err = controller.delete(tx, c.Params, c.Request.URL.Query())
	if err != nil {
		tx.Rollback()
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	controller.outputter.OutputDelete(c, http.StatusNoContent)
}

// Patch corresponds HTTP PATCH message and handles a request for a single resource to update partially the specific information
func (controller *BaseController) Patch(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	result, err := controller.patch(tx, c.Params, c.Request.URL.Query())
	if err != nil {
		tx.Rollback()
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	controller.outputter.OutputPatch(c, http.StatusOK, result)
}

// Options corresponds HTTP OPTIONS message and handles a request for multi resources to retrieve its supported options
func (controller *BaseController) Options(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	err = controller.options(tx, c.Params, c.Request.URL.Query())
	if err != nil {
		tx.Rollback()
		controller.outputter.OutputError(c, http.StatusBadRequest, err)
		return
	}

	tx.Commit()

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	controller.outputter.OutputOptions(c, http.StatusNoContent)
}

// SetOutputter sets an outputter for this controller
func (controller *BaseController) SetOutputter(outputter extensions.Outputter) {
	controller.outputter = outputter
}

// SetQueryCustomizer sets a query customizer for this controller
func (controller *BaseController) SetQueryCustomizer(queryCustomizer extensions.QueryCustomizer) {
	controller.queryCustomizer = queryCustomizer
}
