package controller

import (
	"strings"

	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/helper"

	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/util/conversion"
	"github.com/qb0C80aE/clay/util/mapstruct"
	"github.com/qb0C80aE/clay/version"
	validatorpkg "gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

// BaseController is the base class that all controller classes inherit
type BaseController struct {
	actualController extension.Controller
	model            extension.Model
	binder           extension.Binder
	outputHandler    extension.OutputHandler
	queryCustomizer  extension.QueryCustomizer
}

var validator = validatorpkg.New()

var intTypeString = reflect.TypeOf(int(0)).String()
var int8TypeString = reflect.TypeOf(int8(0)).String()
var int16TypeString = reflect.TypeOf(int16(0)).String()
var int32TypeString = reflect.TypeOf(int32(0)).String()
var int64TypeString = reflect.TypeOf(int64(0)).String()
var uintTypeString = reflect.TypeOf(uint(0)).String()
var uint8TypeString = reflect.TypeOf(uint8(0)).String()
var uint16TypeString = reflect.TypeOf(uint16(0)).String()
var uint32TypeString = reflect.TypeOf(uint32(0)).String()
var uint64TypeString = reflect.TypeOf(uint64(0)).String()
var float32TypeString = reflect.TypeOf(float32(0)).String()
var float64TypeString = reflect.TypeOf(float64(0)).String()
var booleanTypeString = reflect.TypeOf(true).String()
var stringTypeString = reflect.TypeOf("").String()
var bytesTypeString = reflect.TypeOf([]byte{}).String()
var stringInterfaceMapTypeString = reflect.TypeOf(map[string]interface{}{}).String()

var formKeyHashRegExp, _ = regexp.Compile("^([a-zA-Z_][a-zA-Z0-9_]*)\\[([a-zA-Z_][a-zA-Z0-9_]*)\\]$")

func getMapingTagValue(structField *reflect.StructField, tagName string) string {
	tag := structField.Name
	tagValueList := strings.Split(structField.Tag.Get(tagName), ",")
	for _, tagValue := range tagValueList {
		switch tagValue {
		case "omitempty":
			continue
		default:
			tag = tagValue
			break
		}
	}
	return tag
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

func executeValidation(c *gin.Context, resourceName string, inputContainer interface{}, keyParameterSpecifier string) error {
	model, err := extension.GetAssociatedModelWithResourceName(resourceName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	// for update method, set/override key parameter to container
	modelKey, err := extension.GetModelKey(model, keyParameterSpecifier)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}
	keyParameterValue := c.Param(modelKey.KeyParameter)
	if len(keyParameterValue) > 0 {
		containerValue := reflect.ValueOf(inputContainer)
		for containerValue.Kind() == reflect.Ptr {
			containerValue = containerValue.Elem()
		}
		keyField := containerValue.FieldByName(modelKey.KeyField)
		switch keyField.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := strconv.Atoi(keyParameterValue)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return err
			}
			keyField.SetInt(int64(v))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := strconv.Atoi(keyParameterValue)
			if err != nil {
				logging.Logger().Debug(err.Error())
				return err
			}
			keyField.SetUint(uint64(v))
		default:
			keyField.SetString(fmt.Sprintf("%s", keyParameterValue))
		}
	}

	// validate again with validator.v9, with not "binding" tag, but "validate" tag
	validator.SetTagName("validate")
	if err := validator.Struct(inputContainer); err != nil {
		logging.Logger().Debug(err.Error())
		return err
	}

	return nil
}

func (receiver *BaseController) determineResponseContentTypeFromAccept(c *gin.Context) string {
	acceptList := strings.Split(c.Request.Header.Get("Accept"), ",")

	if len(acceptList) == 0 {
		return extension.AcceptJSON
	}

	result := strings.Trim(acceptList[0], " ")
	switch result {
	case extension.AcceptAll:
		return extension.AcceptJSON
	default:
		return result
	}
}

func (receiver *BaseController) determineResponseCharsetTypeFromAcceptCharset(c *gin.Context) string {
	acceptCharsetList := strings.Split(c.Request.Header.Get("Accept-Charset"), ",")

	switch len(acceptCharsetList) {
	case 0:
		return extension.AcceptCharsetUTF8
	default:
		firstAcceptCharset := strings.Trim(acceptCharsetList[0], " ")
		if len(firstAcceptCharset) == 0 {
			return extension.AcceptCharsetUTF8
		}
		return firstAcceptCharset
	}
}

func (receiver *BaseController) determineResponseMappingTagFromAccept(c *gin.Context) string {
	switch receiver.determineResponseContentTypeFromAccept(c) {
	case extension.AcceptXYAML, extension.AcceptTextYAML:
		return extension.TagYAML
	default:
		return extension.TagJSON
	}
}

func (receiver *BaseController) outputTextWithContentType(c *gin.Context, code int, result interface{}) {
	text := result.(string)

	accept := receiver.determineResponseContentTypeFromAccept(c)

	switch accept {
	case "", extension.AcceptAll:
		c.String(code, text)
	default:
		acceptCharset := receiver.determineResponseCharsetTypeFromAcceptCharset(c)
		var contentType string
		if len(strings.Trim(acceptCharset, " ")) > 0 {
			contentType = fmt.Sprintf("%s; charset=%s", accept, acceptCharset)
		} else {
			contentType = accept
		}
		c.Data(code, contentType, []byte(text))
	}
}

// Bind binds input data to a container instance
func (receiver *BaseController) Bind(c *gin.Context, resourceName string) (interface{}, error) {
	preloadedBody := c.MustGet("PreloadedBody").([]byte)

	switch c.ContentType() {
	case extension.ContentTypeJSON, extension.ContentTypeXYAML, extension.ContentTypeTextYAML:
		inputMap := map[string]interface{}{}
		var data interface{}
		switch c.ContentType() {
		case extension.ContentTypeJSON:
			if err := json.Unmarshal(preloadedBody, &inputMap); err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			data = inputMap
		case extension.ContentTypeXYAML, extension.ContentTypeTextYAML:
			if err := yaml.Unmarshal(preloadedBody, &inputMap); err != nil {
				logging.Logger().Debug(err.Error())
				return nil, err
			}
			// yaml.Unmarshal builds map[interface{}]interface{} what cannot be used in json.Marshal
			// so convert it to string-keyed map
			data = conversion.GetUtility().ConvertToStringKeyMap(inputMap)
		}

		container, err := extension.CreateInputContainerByResourceName(resourceName, data)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}

		// remap inputMap to container
		// instead of c.Bind (c.Bind does not bind yaml correctly)
		if err := mapstruct.GetUtility().MapToStruct(data, container); err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}

		if err := executeValidation(c, resourceName, container, c.Request.URL.Query().Get("key_parameter")); err != nil {
			return nil, err
		}

		return container, nil
	case extension.ContentTypeMultipartFormData:
		inputMap := map[string]interface{}{}
		container, err := extension.CreateInputContainerByResourceName(resourceName, inputMap)
		if err != nil {
			logging.Logger().Debug(err.Error())
			return nil, err
		}

		vs := reflect.ValueOf(container)
		for vs.Kind() == reflect.Ptr {
			vs = vs.Elem()
		}
		if !vs.IsValid() {
			logging.Logger().Debug("invalid model")
			return nil, errors.New("invalid model")
		}
		if !vs.CanInterface() {
			logging.Logger().Debug("invalid model")
			return nil, errors.New("invalid model")
		}
		value := vs.Interface()
		t := reflect.TypeOf(value)

		// bind once except for form fields which are unable to bind
		// in order to do that, ignore error
		c.Bind(container)

		// bind the multipart form manually
		if c.Request.ParseMultipartForm(1024*1024) == nil {

			// form fields binding process (in case of Map)
			for formKey := range c.Request.Form {
				formValue := c.Request.Form.Get(formKey)
				submatchList := formKeyHashRegExp.FindStringSubmatch(formKey)

				if len(submatchList) == 3 {
					formMainKey := submatchList[1]
					formSubKey := submatchList[2]
					for i := 0; i < t.NumField(); i++ {
						field := t.Field(i)
						formTag := getMapingTagValue(&field, "form")

						if formTag == formMainKey {
							switch field.Type.String() {
							case stringInterfaceMapTypeString:
								if vs.FieldByName(field.Name).IsNil() {
									vs.FieldByName(field.Name).Set(reflect.ValueOf(map[string]interface{}{}))
								}
								vs.FieldByName(field.Name).SetMapIndex(reflect.ValueOf(formSubKey), reflect.ValueOf(formValue))
								break
							default:
								logging.Logger().Debug("invalid field definition, the field must be map[string]interface{}")
								return nil, errors.New("invalid field definition, the field must be map[string]interface{}")
							}
						}
					}
				} else {
					for i := 0; i < t.NumField(); i++ {
						field := t.Field(i)
						formTag := getMapingTagValue(&field, "form")

						if formTag == formKey {
							switch field.Type.String() {
							case intTypeString, int8TypeString, int16TypeString, int32TypeString, int64TypeString:
								intValue, err := strconv.ParseInt(formValue, 10, 64)
								if err != nil {
									logging.Logger().Debug(err.Error())
									return nil, err
								}
								vs.FieldByName(field.Name).SetInt(intValue)
								break
							case uintTypeString, uint8TypeString, uint16TypeString, uint32TypeString, uint64TypeString:
								uintValue, err := strconv.ParseUint(formValue, 10, 64)
								if err != nil {
									logging.Logger().Debug(err.Error())
									return nil, err
								}
								vs.FieldByName(field.Name).SetUint(uintValue)
								break
							case float32TypeString, float64TypeString:
								floatValue, err := strconv.ParseFloat(formValue, 64)
								if err != nil {
									logging.Logger().Debug(err.Error())
									return nil, err
								}
								vs.FieldByName(field.Name).SetFloat(floatValue)
								break
							case booleanTypeString:
								booleanValue, err := strconv.ParseBool(formValue)
								if err != nil {
									logging.Logger().Debug(err.Error())
									return nil, err
								}
								vs.FieldByName(field.Name).SetBool(booleanValue)
								break
							case stringTypeString:
								vs.FieldByName(field.Name).SetString(formValue)
								break
							case bytesTypeString:
								vs.FieldByName(field.Name).SetBytes([]byte(formValue))
								break
							default:
								logging.Logger().Debug("invalid field definition, the field must be string or slice of bytes")
								return nil, errors.New("invalid field definition, the field must be string or slice of bytes")
							}
						}
					}
				}
			}

			// uploaded files binding process
			files := c.Request.MultipartForm.File
			for formKey := range files {
				buffer := &bytes.Buffer{}
				file, _, err := c.Request.FormFile(formKey)
				if err != nil {
					logging.Logger().Debug(err.Error())
					return nil, err
				}
				defer file.Close()

				data, err := ioutil.ReadAll(file)

				if err != nil {
					logging.Logger().Debug(err.Error())
					return nil, err
				}

				_, err = buffer.Write(data)
				if err != nil {
					logging.Logger().Debug(err.Error())
					return nil, err
				}

				submatchList := formKeyHashRegExp.FindStringSubmatch(formKey)
				if len(submatchList) == 3 {
					formMainKey := submatchList[1]
					formSubKey := submatchList[2]

					for i := 0; i < t.NumField(); i++ {
						field := t.Field(i)
						formTag := getMapingTagValue(&field, "form")

						if formTag == formMainKey {
							switch field.Type.String() {
							case stringInterfaceMapTypeString:
								if vs.FieldByName(field.Name).IsNil() {
									vs.FieldByName(field.Name).Set(reflect.ValueOf(map[string]interface{}{}))
								}
								vs.FieldByName(field.Name).SetMapIndex(reflect.ValueOf(formSubKey), reflect.ValueOf(buffer.Bytes()))
								break
							default:
								logging.Logger().Debug("invalid field definition, the field must be map[string]interface{}")
								return nil, errors.New("invalid field definition, the field must be map[string]interface{}")
							}
						}
					}
				} else {
					for i := 0; i < t.NumField(); i++ {
						field := t.Field(i)
						formTag := getMapingTagValue(&field, "form")

						if formTag == formKey {
							switch field.Type.String() {
							case stringTypeString:
								vs.FieldByName(field.Name).SetString(buffer.String())
								break
							case bytesTypeString:
								vs.FieldByName(field.Name).SetBytes(buffer.Bytes())
								break
							default:
								logging.Logger().Debug("invalid field definition, the field must be string or slice of bytes")
								return nil, errors.New("invalid field definition, the field must be string or slice of bytes")
							}
						}
					}
				}
			}
		}

		if err := executeValidation(c, resourceName, container, c.Request.URL.Query().Get("key_parameter")); err != nil {
			return nil, err
		}

		return container, nil
	}

	return nil, fmt.Errorf("Content-Type %s is not supported", c.ContentType())
}

// GetModel returns its model
func (receiver *BaseController) GetModel() extension.Model {
	return receiver.model
}

// GetResourceName returns its resource/table name in REST/DB
func (receiver *BaseController) GetResourceName() (string, error) {
	return extension.GetAssociatedResourceNameWithModel(receiver.model)
}

// GetResourceSingleURL builds a resource url what represents a single resource based on the argument
func (receiver *BaseController) GetResourceSingleURL() (string, error) {
	resourceName, err := receiver.GetResourceName()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s/:key_parameter", resourceName), nil
}

// GetResourceMultiURL builds a resource url what represents multi resources based on the argument
func (receiver *BaseController) GetResourceMultiURL() (string, error) {
	resourceName, err := receiver.GetResourceName()
	if err != nil {
		logging.Logger().Debug(err.Error())
		return "", err
	}

	return fmt.Sprintf("%s", resourceName), nil
}

func (receiver *BaseController) outputDataWithType(c *gin.Context, code int, obj interface{}) {
	switch receiver.determineResponseContentTypeFromAccept(c) {
	case extension.AcceptXYAML, extension.AcceptTextYAML:
		c.YAML(code, obj)
	default:
		// default is json
		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(code, obj)
		} else {
			c.JSON(code, obj)
		}
	}
}

// OutputError handles an error output
func (receiver *BaseController) OutputError(c *gin.Context, code int, err error) {
	receiver.outputDataWithType(c, code, gin.H{"error": err.Error()})
}

// OutputGetSingle corresponds HTTP GET message and handles the output of a single result from logic classes
func (receiver *BaseController) OutputGetSingle(c *gin.Context, code int, result interface{}, fields map[string]interface{}) {
	_, allFieldExists := fields["*"]
	if (fields == nil) || ((len(fields) == 1) && allFieldExists) {
		receiver.outputDataWithType(c, code, result)
	} else {
		targetTag := ""
		switch receiver.determineResponseContentTypeFromAccept(c) {
		case extension.AcceptXYAML, extension.AcceptTextYAML:
			targetTag = extension.TagYAML
		default:
			targetTag = extension.TagJSON
		}

		fieldMap, err := helper.FieldToMap(result, fields, targetTag)
		if err != nil {
			logging.Logger().Debug(err.Error())
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
			return
		}

		receiver.outputDataWithType(c, code, fieldMap)
	}
}

// OutputGetMulti corresponds HTTP GET message and handles the output of multiple result from logic classes
func (receiver *BaseController) OutputGetMulti(c *gin.Context, code int, result interface{}, total int, countBeforePagination int, fields map[string]interface{}) {
	c.Header("Total", strconv.Itoa(total))
	c.Header("Count-Before-Pagination", strconv.Itoa(countBeforePagination))
	_, allFieldExists := fields["*"]
	if (fields == nil) || ((len(fields) == 1) && allFieldExists) {
		receiver.outputDataWithType(c, code, result)
	} else {
		targetTag := ""
		switch receiver.determineResponseContentTypeFromAccept(c) {
		case extension.AcceptXYAML, extension.AcceptTextYAML:
			targetTag = extension.TagYAML
		default:
			targetTag = extension.TagJSON
		}

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

				fieldMap, err := helper.FieldToMap(item.Interface(), fields, targetTag)

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

				fieldMap, err := helper.FieldToMap(item.Interface(), fields, targetTag)

				if err != nil {
					logging.Logger().Debug(err.Error())
					receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
					return
				}

				fieldMaps = append(fieldMaps, fieldMap)
			}

			receiver.outputDataWithType(c, code, fieldMaps)
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

	resourceName, err := receiver.model.GetResourceName(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	resultForQueryFields, err := extension.CreateOutputContainerByResourceName(resourceName, c.Request.URL.Query().Get("preloads"))
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db = parameter.SetPreloads(db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(resultForQueryFields, fields, receiver.determineResponseMappingTagFromAccept(c))

	result, err := receiver.model.GetSingle(receiver.model, db, c.Params, c.Request.URL.Query(), queryFields)
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

	resourceName, err := receiver.model.GetResourceName(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	elementOfResultForQueryFields, err := extension.CreateOutputContainerByResourceName(resourceName, c.Request.URL.Query().Get("preloads"))
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(elementOfResultForQueryFields, fields, receiver.determineResponseMappingTagFromAccept(c))

	result, err := receiver.model.GetMulti(receiver.model, db, c.Params, c.Request.URL.Query(), queryFields)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	// reset all conditions in order to get the total number of records
	db = db.New()
	total, err := receiver.model.GetCount(receiver.model, db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	// reset conditions except for limit and offset in order to get the record count before limitation
	db = db.New()
	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	countBeforePagination, err := receiver.model.GetCount(receiver.model, db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	_, first := c.GetQuery("first")
	if first {
		resultValue := reflect.ValueOf(result)
		if resultValue.Len() == 0 {
			receiver.outputHandler.OutputError(c, http.StatusBadRequest, errors.New("no records retrieved"))
		} else {
			receiver.outputHandler.OutputGetSingle(c, http.StatusOK, reflect.ValueOf(result).Index(0).Interface(), fields)
		}
	} else {
		receiver.outputHandler.OutputGetMulti(c, http.StatusOK, result, total, countBeforePagination, fields)
	}
}

// Create corresponds HTTP POST message and handles a request for multi resource to create a new information
func (receiver *BaseController) Create(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	resourceName, err := receiver.model.GetResourceName(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	container, err := receiver.binder.Bind(c, resourceName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()

	result, err := receiver.model.Create(receiver.model, tx, c.Params, c.Request.URL.Query(), container)
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

	resourceName, err := receiver.model.GetResourceName(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	container, err := receiver.binder.Bind(c, resourceName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	result, err := receiver.model.Update(receiver.model, tx, c.Params, c.Request.URL.Query(), container)
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
	err = receiver.model.Delete(receiver.model, tx, c.Params, c.Request.URL.Query())
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

	resourceName, err := receiver.model.GetResourceName(receiver.model)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	container, err := receiver.binder.Bind(c, resourceName)
	if err != nil {
		logging.Logger().Debug(err.Error())
		receiver.outputHandler.OutputError(c, http.StatusBadRequest, err)
		return
	}

	db := dbpkg.Instance(c)

	tx := db.Begin()
	result, err := receiver.model.Patch(receiver.model, tx, c.Params, c.Request.URL.Query(), container)
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
	err = receiver.model.GetOptions(receiver.model, tx, c.Params, c.Request.URL.Query())
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

// DoBeforeDBMigration execute initialization process before DB migration
func (receiver *BaseController) DoBeforeDBMigration(_ *gorm.DB) error {
	return nil
}

// DoAfterDBMigration execute initialization process after DB migration
func (receiver *BaseController) DoAfterDBMigration(_ *gorm.DB) error {
	return nil
}

// DoBeforeRouterSetup execute initialization process before Router initialization
func (receiver *BaseController) DoBeforeRouterSetup(_ *gin.Engine) error {
	return nil
}

// DoAfterRouterSetup execute initialization process after Router initialization
func (receiver *BaseController) DoAfterRouterSetup(_ *gin.Engine) error {
	return nil
}
