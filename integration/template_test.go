package integration

import (
	"github.com/qb0C80aE/clay/models"
	"net/http"
	"strconv"
	"testing"
)

// +build integration

func TestGetTemplates_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceUrl(server, "templates", nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, EmptyArrayString, []*models.Template{})
}

func TestGetTemplates(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template1 := &models.Template{
		Name:            "test1",
		TemplateContent: "TestTemplate1",
	}

	template2 := &models.Template{
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name:  "testParameter1",
				Value: "TestParameter1",
			},
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template2)

	parameters := map[string]string{
		"preloads": "TemplateExternalParameters",
	}
	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceUrl(server, "templates", parameters), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGetTemplates_1.json"), []*models.Template{})
}

func TestCreateTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 1
	template := &models.Template{
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)
	CheckResponseJson(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplate_1.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplate_2.json"), &models.Template{})
}

func TestCreateTemplateWithID(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &models.Template{
		ID:              id,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)
	CheckResponseJson(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplateWithID_1.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplateWithID_2.json"), &models.Template{})
}

func TestCreateTemplate_WithID_WithTemplateExternalParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &models.Template{
		ID:              100,
		Name:            "test",
		TemplateContent: "TestTemplate",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name:  "testParameter1",
				Value: "TestParameter1",
			},
			{
				ID:    10,
				Name:  "testParameter10",
				Value: "TestParameter10",
			},
		},
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)
	CheckResponseJson(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplate_WithID_WithTemplateExternalParameters_1.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplate_WithID_WithTemplateExternalParameters_2.json"), &models.Template{})

	parameters := map[string]string{
		"preloads": "TemplateExternalParameters",
	}
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), parameters), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplate_WithID_WithTemplateExternalParameters_3.json"), &models.Template{})
}

func TestUpdateTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &models.Template{
		ID:              id,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)

	template.Name = "testUpdated"
	template.TemplateContent = "TestTemplateUpdated"

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), template)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_1.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_2.json"), &models.Template{})
}

func TestUpdateTemplate_WithTemplateExternalParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &models.Template{
		ID:              id,
		Name:            "test",
		TemplateContent: "TestTemplate",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name:  "testParameter1",
				Value: "TestParameter1",
			},
			{
				Name:  "testParameter2",
				Value: "TestParameter2",
			},
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)

	template.Name = "testUpdated"
	template.TemplateContent = "TestTemplateUpdated"
	template.TemplateExternalParameters = []*models.TemplateExternalParameter{
		{
			ID:    2,
			Name:  "testParameter2Updated",
			Value: "TestParameter2Updated",
		},
		{
			Name:  "testParameterNew",
			Value: "TestParameterNew",
		},
	}

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), template)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_WithTemplateExternalParameters_1.json"), &models.Template{})

	parameters := map[string]string{
		"preloads": "TemplateExternalParameters",
	}
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), parameters), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_WithTemplateExternalParameters_2.json"), &models.Template{})
}

func TestDeleteTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &models.Template{
		ID:              id,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)

	responseText, code := Execute(t, http.MethodDelete, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseJson(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplate_1.json"), &ErrorResponseText{})
}

func TestPatchTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &models.Template{
		ID:              id,
		Name:            "test",
		TemplateContent: "{{.TemplateExternalParameters.testParameter1}} is TestParameter1, {{.TemplateExternalParameters.testParameter2}} is TestParameter2.",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name:  "testParameter1",
				Value: "TestParameter1",
			},
			{
				Name:  "testParameter2",
				Value: "TestParameter2",
			},
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)

	responseText, code := Execute(t, http.MethodPatch, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(id), nil), template)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestPatchTemplate_1.txt"))
}

func TestGetTemplateExternalParameters_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceUrl(server, "template_external_parameters", nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, EmptyArrayString, []*models.TemplateExternalParameter{})
}

func TestCreateTemplateExternalParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)

	templateExternalParameter1 := &models.TemplateExternalParameter{
		TemplateID: 1,
		Name:       "testParameter1",
		Value:      "TestParameter1",
	}
	templateExternalParameter2 := &models.TemplateExternalParameter{
		TemplateID: 1,
		Name:       "testParameter2",
		Value:      "TestParameter2",
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "template_external_parameters", nil), templateExternalParameter1)
	CheckResponseJson(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplateExternalParameter_1.json"), &models.TemplateExternalParameter{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "template_external_parameters", nil), templateExternalParameter2)
	CheckResponseJson(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplateExternalParameter_2.json"), &models.TemplateExternalParameter{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceUrl(server, "template_external_parameters", nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplateExternalParameter_3.json"), []*models.TemplateExternalParameter{})
}

func TestUpdateTemplateExternalParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)

	id := 1
	templateExternalParameter1 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter1",
		Value:      "TestParameter1",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "template_external_parameters", nil), templateExternalParameter1)

	templateExternalParameter1.Name = "templateExternalParameter1Updated"
	templateExternalParameter1.Value = "TestParameter1Updated"

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceUrl(server, "template_external_parameters", strconv.Itoa(id), nil), templateExternalParameter1)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplateExternalParameter_1.json"), &models.TemplateExternalParameter{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "template_external_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplateExternalParameter_2.json"), &models.TemplateExternalParameter{})
}

func TestDeleteTemplateExternalParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)

	id := 1
	templateExternalParameter1 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter1",
		Value:      "TestParameter1",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "template_external_parameters", nil), templateExternalParameter1)

	templateExternalParameter1.Name = "templateExternalParameter1Updated"
	templateExternalParameter1.Value = "TestParameter1Updated"

	responseText, code := Execute(t, http.MethodDelete, GenerateSingleResourceUrl(server, "template_external_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "template_external_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJson(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplateExternalParameter_1.json"), &ErrorResponseText{})
}

func TestDeleteTemplateExternalParameters_Cascade(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	templateId := 1
	template := &models.Template{
		ID:              templateId,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template)

	id := 1
	templateExternalParameter1 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter1",
		Value:      "TestParameter1",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "template_external_parameters", nil), templateExternalParameter1)

	Execute(t, http.MethodDelete, GenerateSingleResourceUrl(server, "templates", strconv.Itoa(templateId), nil), nil)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "template_external_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJson(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplateExternalParameter_Cascade_1.json"), &ErrorResponseText{})
}
