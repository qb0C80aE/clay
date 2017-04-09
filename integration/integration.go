package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/server"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

const version = "v1"
const timeout = 15

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
		url = fmt.Sprintf("%s/%s/%s?%s", ts.URL, version, resource, generateQueryParameter(parameters))
	} else {
		url = fmt.Sprintf("%s/%s/%s", ts.URL, version, resource)
	}
	return url
}

// GenerateSingleResourceURL generates a resource url what represents a single resource based on the arguments
func GenerateSingleResourceURL(ts *httptest.Server, resource string, id string, parameters map[string]string) string {
	var url string
	if 0 < len(parameters) {
		url = fmt.Sprintf("%s/%s/%s/%s?%s", ts.URL, version, resource, id, generateQueryParameter(parameters))
	} else {
		url = fmt.Sprintf("%s/%s/%s/%s", ts.URL, version, resource, id)
	}
	return url
}

// SetupServer setups server for integration tests
func SetupServer() *httptest.Server {
	database := db.Connect()
	s := server.Setup(database)
	return httptest.NewServer(s)
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
		error(t, "code is expected as %d, but %d", expectedCode, code)
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
		error(t, "code is expected as %d, but %d", expectedCode, code)
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
