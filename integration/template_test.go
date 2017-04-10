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

	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "templates", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, EmptyArrayString, []*models.Template{})
}

func TestCreateTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	parameters := map[string]string{
		"preloads": "TemplateExternalParameters",
	}

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

	template3 := &models.Template{
		ID:              100,
		Name:            "test100",
		TemplateContent: "TestTemplate100",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name:  "testParameter100",
				Value: "TestParameter100",
			},
			{
				ID:    10,
				Name:  "testParameter110",
				Value: "TestParameter110",
			},
		},
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplate_1.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplate_2.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template3)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplate_3.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "templates", parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplate_4.json"), []*models.Template{})
}

func TestUpdateTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	parameters := map[string]string{
		"preloads": "TemplateExternalParameters",
	}

	id1 := 101
	template1 := &models.Template{
		ID:              id1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
	}

	id2 := 102
	template2 := &models.Template{
		ID:              id2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name:  "testParameter21",
				Value: "TestParameter21",
			},
		},
	}

	id3 := 103
	template3 := &models.Template{
		ID:              id3,
		Name:            "test3",
		TemplateContent: "TestTemplate3",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name:  "testParameter31",
				Value: "TestParameter31",
			},
			{
				Name:  "testParameter32",
				Value: "TestParameter32",
			},
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template3)

	template1.Name = "test1Updated"
	template1.TemplateContent = "TestTemplate1Updated"

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id1), nil), template1)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_1.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id1), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_2.json"), &models.Template{})

	template2.TemplateExternalParameters = nil

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id2), nil), template2)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_3.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id2), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_4.json"), &models.Template{})

	template3.TemplateExternalParameters[1].Name = "testParameter32Updated"
	template3.TemplateExternalParameters[1].Value = "TestParameter32Updated"
	template3.TemplateExternalParameters = append(
		template3.TemplateExternalParameters,
		&models.TemplateExternalParameter{
			Name:  "testParameter33",
			Value: "TestParameter33",
		},
	)

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id3), nil), template3)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_5.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id3), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_6.json"), &models.Template{})
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

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	responseText, code := Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplate_1.json"), &ErrorResponseText{})
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

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	responseText, code := Execute(t, http.MethodPatch, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestPatchTemplate_1.txt"))
}

func TestGetTemplateExternalParameters_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "template_external_parameters", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, EmptyArrayString, []*models.TemplateExternalParameter{})
}

func TestCreateTemplateExternalParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

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

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter1)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplateExternalParameter_1.json"), &models.TemplateExternalParameter{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter2)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplateExternalParameter_2.json"), &models.TemplateExternalParameter{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "template_external_parameters", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplateExternalParameter_3.json"), []*models.TemplateExternalParameter{})
}

func TestUpdateTemplateExternalParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templateExternalParameter1 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter1",
		Value:      "TestParameter1",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter1)

	templateExternalParameter1.Name = "templateExternalParameter1Updated"
	templateExternalParameter1.Value = "TestParameter1Updated"

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "template_external_parameters", strconv.Itoa(id), nil), templateExternalParameter1)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplateExternalParameter_1.json"), &models.TemplateExternalParameter{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_external_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplateExternalParameter_2.json"), &models.TemplateExternalParameter{})
}

func TestDeleteTemplateExternalParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templateExternalParameter1 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter1",
		Value:      "TestParameter1",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter1)

	templateExternalParameter1.Name = "templateExternalParameter1Updated"
	templateExternalParameter1.Value = "TestParameter1Updated"

	responseText, code := Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "template_external_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_external_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplateExternalParameter_1.json"), &ErrorResponseText{})
}

func TestDeleteTemplateExternalParameters_Cascade(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	templateID := 1
	template := &models.Template{
		ID:              templateID,
		Name:            "test",
		TemplateContent: "TestTemplate",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templateExternalParameter1 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter1",
		Value:      "TestParameter1",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter1)

	Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "templates", strconv.Itoa(templateID), nil), nil)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_external_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplateExternalParameter_Cascade_1.json"), &ErrorResponseText{})
}

func TestTemplate_ExtractFromDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template1 := &models.Template{
		ID:              1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
	}
	template2 := &models.Template{
		ID:              2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)

	templateExternalParameter11 := &models.TemplateExternalParameter{
		TemplateID: 1,
		Name:       "testParameter11",
		Value:      "TestParameter11",
	}
	templateExternalParameter12 := &models.TemplateExternalParameter{
		TemplateID: 1,
		Name:       "testParameter12",
		Value:      "TestParameter12",
	}
	templateExternalParameter21 := &models.TemplateExternalParameter{
		TemplateID: 2,
		Name:       "testParameter21",
		Value:      "TestParameter21",
	}
	templateExternalParameter22 := &models.TemplateExternalParameter{
		TemplateID: 2,
		Name:       "testParameter22",
		Value:      "TestParameter22",
	}
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter11)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter12)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter21)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter22)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_ExtractFromDesign_1.json"), &models.Design{})
}

func TestTemplate_LoadToDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template1 := &models.Template{
		ID:              1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
	}
	template2 := &models.Template{
		ID:              2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
	}
	templateExternalParameter11 := &models.TemplateExternalParameter{
		ID:         1,
		TemplateID: 1,
		Name:       "testParameter11",
		Value:      "TestParameter11",
	}
	templateExternalParameter12 := &models.TemplateExternalParameter{
		ID:         2,
		TemplateID: 1,
		Name:       "testParameter12",
		Value:      "TestParameter12",
	}
	templateExternalParameter21 := &models.TemplateExternalParameter{
		ID:         3,
		TemplateID: 2,
		Name:       "testParameter21",
		Value:      "TestParameter21",
	}
	templateExternalParameter22 := &models.TemplateExternalParameter{
		ID:         4,
		TemplateID: 2,
		Name:       "testParameter22",
		Value:      "TestParameter22",
	}

	design := &models.Design{
		Content: map[string]interface{}{
			"templates": []interface{}{
				template1,
				template2,
			},
			"template_external_parameters": []interface{}{
				templateExternalParameter11,
				templateExternalParameter12,
				templateExternalParameter21,
				templateExternalParameter22,
			},
		},
	}

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "designs", "present", nil), design)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_LoadToDesign_1.json"), &models.Design{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_LoadToDesign_2.json"), &models.Design{})
}

func TestTemplate_DeleteFromDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template1 := &models.Template{
		ID:              1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
	}
	template2 := &models.Template{
		ID:              2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
	}
	templateExternalParameter11 := &models.TemplateExternalParameter{
		ID:         1,
		TemplateID: 1,
		Name:       "testParameter11",
		Value:      "TestParameter11",
	}
	templateExternalParameter12 := &models.TemplateExternalParameter{
		ID:         2,
		TemplateID: 1,
		Name:       "testParameter12",
		Value:      "TestParameter12",
	}
	templateExternalParameter21 := &models.TemplateExternalParameter{
		ID:         3,
		TemplateID: 2,
		Name:       "testParameter21",
		Value:      "TestParameter21",
	}
	templateExternalParameter22 := &models.TemplateExternalParameter{
		ID:         4,
		TemplateID: 2,
		Name:       "testParameter22",
		Value:      "TestParameter22",
	}

	design := &models.Design{
		Content: map[string]interface{}{
			"templates": []interface{}{
				template1,
				template2,
			},
			"template_external_parameters": []interface{}{
				templateExternalParameter11,
				templateExternalParameter12,
				templateExternalParameter21,
				templateExternalParameter22,
			},
		},
	}

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "designs", "present", nil), design)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_DeleteFromDesign_1.json"), &models.Design{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_DeleteFromDesign_2.json"), &models.Design{})

	responseText, code = Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_DeleteFromDesign_3.json"), &models.Design{})
}

func TestTemplate_GenerateTemplateParameter(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 1
	template1 := &models.Template{
		ID:              id,
		Name:            "test1",
		TemplateContent: "{{range $key, $value := .TemplateExternalParameters}}{{$key}}={{$value}}\n{{end}}",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)

	templateExternalParameter11 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter11",
		Value:      "TestParameter11",
	}
	templateExternalParameter12 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter12",
		Value:      "TestParameter12",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter11)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter12)

	responseText, code := Execute(t, http.MethodPatch, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_GenerateTemplateParameter_1.txt"))
}
