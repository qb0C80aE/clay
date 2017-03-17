package integration

import (
	"github.com/qb0C80aE/clay/models"
	"net/http"
	"testing"
)

// +build integration

func TestGetDesign_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "designs", "present", nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestGetDesign_Empty_1.json"), &models.Design{})
}

func TestGetDesign(t *testing.T) {
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
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name:  "testParameter1",
				Value: "TestParameter1",
			},
		},
	}

	templateExternalParameter22 := &models.TemplateExternalParameter{
		TemplateID: 2,
		Name:       "testParameter2",
		Value:      "TestParameter2",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "templates", nil), template2)
	Execute(t, http.MethodPost, GenerateMultiResourceUrl(server, "template_external_parameters", nil), templateExternalParameter22)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "designs", "present", nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestGetDesign_1.json"), &models.Design{})
}

func TestUpdateDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	design := &models.Design{
		Content: map[string]interface{}{
			"template_external_parameters": []*models.TemplateExternalParameter{
				{
					ID:         1,
					TemplateID: 1,
					Name:       "testParameter11",
					Value:      "TestParameter11",
				},
				{
					ID:         2,
					TemplateID: 1,
					Name:       "testParameter12",
					Value:      "TestParameter12",
				},
			},
			"templates": []*models.Template{
				{
					ID:              1,
					Name:            "test1",
					TemplateContent: "TestTemplate1",
				},
				{
					ID:              2,
					Name:            "test2",
					TemplateContent: "TestTemplate2",
				},
			},
		},
	}

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceUrl(server, "designs", "present", nil), design)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestUpdateDesign_1.json"), &models.Design{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "designs", "present", nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestUpdateDesign_2.json"), &models.Design{})
}

func TestDeleteDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	design := &models.Design{
		Content: map[string]interface{}{
			"template_external_parameters": []*models.TemplateExternalParameter{
				{
					ID:         1,
					TemplateID: 1,
					Name:       "testParameter11",
					Value:      "TestParameter11",
				},
				{
					ID:         2,
					TemplateID: 1,
					Name:       "testParameter12",
					Value:      "TestParameter12",
				},
			},
			"templates": []*models.Template{
				{
					ID:              1,
					Name:            "test1",
					TemplateContent: "TestTemplate1",
				},
				{
					ID:              2,
					Name:            "test2",
					TemplateContent: "TestTemplate2",
				},
			},
		},
	}

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceUrl(server, "designs", "present", nil), design)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestDeleteDesign_1.json"), &models.Design{})

	responseText, code = Execute(t, http.MethodDelete, GenerateSingleResourceUrl(server, "designs", "present", nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceUrl(server, "designs", "present", nil), nil)
	CheckResponseJson(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestDeleteDesign_2.json"), &models.Design{})
}
