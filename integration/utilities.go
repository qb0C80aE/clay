// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/middleware"
	"github.com/qb0C80aE/clay/router"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

const timeout = 15

var engine *gin.Engine
var database *gorm.DB
var isDBCustomFunctionSetupCompleted = false

// EmptyArrayString represents an empty JSON array in string type
var EmptyArrayString = []byte("[]")

// ErrorResponseText is the struct represents error JSON responses
type ErrorResponseText struct {
	Error string `json:"error"`
}

func error(t *testing.T, message string, args ...interface{}) {
	result := fmt.Sprintf(message, args...)
	t.Fatalf(result)
}

func generateQueryParameter(data map[string]string) string {
	buffer := &bytes.Buffer{}
	for key, value := range data {
		buffer.WriteString(key)
		buffer.WriteByte('=')
		buffer.WriteString(value)
		buffer.WriteByte('&')
	}
	queryParameter := buffer.String()
	return queryParameter[0 : len(queryParameter)-1]
}

// GenerateMultiResourceURL generates a resource url what represents multi resources based on the arguments
func GenerateMultiResourceURL(ts *httptest.Server, resource string, parameters map[string]string) string {
	var url string
	if 0 < len(parameters) {
		url = fmt.Sprintf("%s/%s?%s", ts.URL, resource, generateQueryParameter(parameters))
	} else {
		url = fmt.Sprintf("%s/%s", ts.URL, resource)
	}
	return url
}

// GenerateSingleResourceURL generates a resource url what represents a single resource based on the arguments
func GenerateSingleResourceURL(ts *httptest.Server, resource string, id string, parameters map[string]string) string {
	var url string
	if 0 < len(parameters) {
		url = fmt.Sprintf("%s/%s/%s?%s", ts.URL, resource, id, generateQueryParameter(parameters))
	} else {
		url = fmt.Sprintf("%s/%s/%s", ts.URL, resource, id)
	}
	return url
}

func setDBtoContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", database)
		c.Next()
	}
}

func doAutoMigrationOnly(db *gorm.DB, initializerList []extension.Initializer, modelList []extension.Model) *gorm.DB {
	db.Exec("pragma foreign_keys = on")

	containerToBeMigratedList := []interface{}{}
	for _, model := range modelList {
		container, err := model.GetContainerForMigration()
		if err != nil {
			logging.Logger().Critical(err.Error())
			panic(err)
		}

		if container != nil {
			containerType := extension.InspectActualElementType(container)
			_, exists := containerType.FieldByName("StructMetaInformation")
			if !exists {
				logging.Logger().Critical("the container does not have StructMetaInformation field, it might not be a container")
				panic(errors.New("the container does not have StructMetaInformation field, it might not be a container"))
			}
			containerToBeMigratedList = append(containerToBeMigratedList, container)
		}
	}

	if err := db.AutoMigrate(containerToBeMigratedList...).Error; err != nil {
		logging.Logger().Critical(err.Error())
		panic(err)
	}

	db.Exec("pragma foreign_keys = off;")

	tx := db.Begin()
	for _, initializer := range initializerList {
		err := initializer.DoAfterDBMigration(tx)
		if err != nil {
			tx.Rollback()
			logging.Logger().Critical(err.Error())
			panic(err)
		}
	}
	tx.Commit()

	db.Exec("pragma foreign_keys = on;")

	return db
}

// SetupServer setups server for integration tests
func SetupServer() *httptest.Server {
	// Initialize DB every test goes
	if !isDBCustomFunctionSetupCompleted {
		dbpkg.SetupCustomDBFunctions()
		isDBCustomFunctionSetupCompleted = true
	}

	newDatabase, err := dbpkg.Connect("memory")
	if err != nil {
		panic(err)
	}
	database = newDatabase

	// Setup engine once when the first test goes
	initializerList := extension.GetRegisteredInitializerList()
	modelList := extension.GetRegisteredModelList()
	if engine == nil {
		// run SetupModel once
		database, err = extension.SetupModel(database, initializerList, modelList)
		if err != nil {
			panic(err)
		}

		// setup engine once, but make engine use new DB
		engine = gin.Default()
		engine.Use(setDBtoContext())
		engine.Use(middleware.PreloadBody())
		if err := router.Setup(engine); err != nil {
			panic(err)
		}
	} else {
		// never run SetupModel again, but do auto-migration only
		doAutoMigrationOnly(database, initializerList, modelList)
		if err != nil {
			panic(err)
		}
	}

	return httptest.NewServer(engine)
}

// Execute send a HTTP request to the test server and receive a response
func Execute(t *testing.T, method string, resourceURL string, data interface{}) ([]byte, int) {
	byteArray, err := json.Marshal(data)

	if err != nil {
		error(t, "error Occured %v", err)
	}

	request, err := http.NewRequest(
		method,
		resourceURL,
		bytes.NewBuffer(byteArray),
	)
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		error(t, "error Occured %v", err)
	}

	client := &http.Client{Timeout: time.Duration(timeout * time.Second)}
	response, err := client.Do(request)

	if err != nil {
		error(t, "%s", err)
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		error(t, "%s", err)
	}
	return contents, response.StatusCode
}

// CheckResponseJSON checks given response JSON text
func CheckResponseJSON(t *testing.T, code int, expectedCode int, responseText []byte, expectedResponseText []byte, model interface{}) {
	if code != expectedCode {
		error(t, "code is expected as %d, but %d\nresponseText: %s", expectedCode, code, string(responseText))
	}

	vs := reflect.ValueOf(model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		error(t, "invalid model")

	}
	if !vs.CanInterface() {
		error(t, "model cannot interface()")
	}
	responseContainer := reflect.New(reflect.TypeOf(vs.Interface())).Interface()
	expectationContainer := reflect.New(reflect.TypeOf(vs.Interface())).Interface()

	err := json.Unmarshal(responseText, &responseContainer)
	if err != nil {
		error(t, "couldn't marshal the responseText: %s", string(responseText))
	}
	jsonByteArray, err := json.Marshal(responseContainer)
	if err != nil {
		error(t, "couldn't unmarshal the responseContainer: %v", responseContainer)
	}
	response := string(jsonByteArray)

	err = json.Unmarshal(expectedResponseText, &expectationContainer)
	if err != nil {
		error(t, "couldn't marshal the expectedResponseText: %s", string(expectedResponseText))
	}
	jsonByteArray, err = json.Marshal(expectationContainer)
	if err != nil {
		error(t, "couldn't unmarshal the expectationContainer: %v", expectationContainer)
	}
	expectation := string(jsonByteArray)

	if response != expectation {
		error(t, "response is expected as '%s', but '%s'", expectation, response)
	}
}

// CheckResponseText checks given response text
func CheckResponseText(t *testing.T, code int, expectedCode int, responseText []byte, expectedResponseText []byte) {
	if code != expectedCode {
		error(t, "code is expected as %d, but %d\nresponseText: %s", expectedCode, code, string(responseText))
	}

	if string(responseText) != string(expectedResponseText) {
		error(t, "response is expected as '%s', but '%s'", expectedResponseText, responseText)
	}
}

// LoadExpectation loads expectation files to test automatically
func LoadExpectation(t *testing.T, testCaseName string) []byte {
	expectationFile := fmt.Sprintf("expectations/%s", testCaseName)
	data, err := ioutil.ReadFile(expectationFile)
	if err != nil {
		error(t, "couldn't load an expectation file %s", expectationFile)
	}
	return data
}
