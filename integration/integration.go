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

var EmptyArrayString = []byte("[]")

type ErrorResponseText struct {
	Error string `json:"error"`
}

func generateQueryParameter(data map[string]string) string {
	var buffer *bytes.Buffer = &bytes.Buffer{}
	for key, value := range data {
		buffer.WriteString(key)
		buffer.WriteByte('=')
		buffer.WriteString(value)
		buffer.WriteByte('&')
	}
	queryParameter := buffer.String()
	return queryParameter[0 : len(queryParameter)-1]
}

func GenerateMultiResourceUrl(ts *httptest.Server, resource string, parameters map[string]string) string {
	var url string
	if 0 < len(parameters) {
		url = fmt.Sprintf("%s/%s/%s?%s", ts.URL, version, resource, generateQueryParameter(parameters))
	} else {
		url = fmt.Sprintf("%s/%s/%s", ts.URL, version, resource)
	}
	return url
}

func GenerateSingleResourceUrl(ts *httptest.Server, resource string, id string, parameters map[string]string) string {
	var url string
	if 0 < len(parameters) {
		url = fmt.Sprintf("%s/%s/%s/%s?%s", ts.URL, version, resource, id, generateQueryParameter(parameters))
	} else {
		url = fmt.Sprintf("%s/%s/%s/%s", ts.URL, version, resource, id)
	}
	return url
}

func SetupServer() *httptest.Server {
	database := db.Connect()
	s := server.Setup(database)
	return httptest.NewServer(s)
}

func Execute(t *testing.T, method string, resourceUrl string, data interface{}) ([]byte, int) {
	byteArray, err := json.Marshal(data)

	if err != nil {
		t.Fatalf("error Occured %v", err)
	}

	request, err := http.NewRequest(
		method,
		resourceUrl,
		bytes.NewBuffer(byteArray),
	)
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		t.Fatalf("error Occured %v", err)
	}

	client := &http.Client{Timeout: time.Duration(timeout * time.Second)}
	response, err := client.Do(request)

	if err != nil {
		t.Fatalf("%s", err)
	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatalf("%s", err)
	}
	return contents, response.StatusCode
}

func CheckResponseJson(t *testing.T, code int, expectedCode int, responseText []byte, expectedResponseText []byte, model interface{}) {
	if code != expectedCode {
		t.Fatalf("code is expected as %d, but %d", expectedCode, code)
	}

	vs := reflect.ValueOf(model)
	for vs.Kind() == reflect.Ptr {
		vs = vs.Elem()
	}
	if !vs.IsValid() {
		t.Fatalf("invalid model")

	}
	if !vs.CanInterface() {
		t.Fatalf("model cannot interface()")
	}
	responseContainer := reflect.New(reflect.TypeOf(vs.Interface())).Interface()
	expectationContainer := reflect.New(reflect.TypeOf(vs.Interface())).Interface()

	err := json.Unmarshal(responseText, &responseContainer)
	if err != nil {
		t.Fatalf("couldn't marshal the responseText: %s", string(responseText))
	}
	jsonByteArray, err := json.Marshal(responseContainer)
	if err != nil {
		t.Fatalf("couldn't unmarshal the responseContainer: %v", responseContainer)
	}
	response := string(jsonByteArray)

	err = json.Unmarshal(expectedResponseText, &expectationContainer)
	if err != nil {
		t.Fatalf("couldn't marshal the expectedResponseText: %s", string(expectedResponseText))
	}
	jsonByteArray, err = json.Marshal(expectationContainer)
	if err != nil {
		t.Fatalf("couldn't unmarshal the expectationContainer: %v", expectationContainer)
	}
	expectation := string(jsonByteArray)

	if response != expectation {
		t.Fatalf("response is expected as '%s', but '%s'", expectation, response)
	}
}

func CheckResponseText(t *testing.T, code int, expectedCode int, responseText []byte, expectedResponseText []byte) {
	if code != expectedCode {
		t.Fatalf("code is expected as %d, but %d", expectedCode, code)
	}

	if string(responseText) != string(expectedResponseText) {
		t.Fatalf("response is expected as '%s', but '%s'", expectedResponseText, responseText)
	}
}

func LoadExpectation(t *testing.T, testCaseName string) []byte {
	expectationFile := fmt.Sprintf("expectations/%s", testCaseName)
	data, err := ioutil.ReadFile(expectationFile)
	if err != nil {
		t.Fatalf("couldn't load an expectation file %s", expectationFile)
	}
	return data
}
