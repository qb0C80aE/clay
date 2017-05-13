// +build integration

package integration

import (
	"database/sql"
	"github.com/qb0C80aE/clay/models"
	"net/http"
	"testing"
)

func TestGetDesign_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestGetDesign_Empty_1.json"), &models.Design{})
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
		TemplatePersistentParameters: []*models.TemplatePersistentParameter{
			{
				Name: "testParameter1",
				ValueString: sql.NullString{
					String: "TestParameter1",
					Valid:  true,
				},
			},
		},
	}

	templatePersistentParameter22 := &models.TemplatePersistentParameter{
		TemplateID: 2,
		Name:       "testParameter2",
		ValueString: sql.NullString{
			String: "TestParameter2",
			Valid:  true,
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter22)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestGetDesign_1.json"), &models.Design{})
}

func TestUpdateDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	design := &models.Design{
		Content: map[string]interface{}{
			"template_persistent_parameters": []*models.TemplatePersistentParameter{
				{
					ID:         1,
					TemplateID: 1,
					Name:       "testParameter11",
					ValueString: sql.NullString{
						String: "TestParameter11",
						Valid:  true,
					},
				},
				{
					ID:         2,
					TemplateID: 1,
					Name:       "testParameter12",
					ValueString: sql.NullString{
						String: "TestParameter12",
						Valid:  true,
					},
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

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "designs", "present", nil), design)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestUpdateDesign_1.json"), &models.Design{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestUpdateDesign_2.json"), &models.Design{})
}

func TestDeleteDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	design := &models.Design{
		Content: map[string]interface{}{
			"template_persistent_parameters": []*models.TemplatePersistentParameter{
				{
					ID:         1,
					TemplateID: 1,
					Name:       "testParameter11",
					ValueString: sql.NullString{
						String: "TestParameter11",
						Valid:  true,
					},
				},
				{
					ID:         2,
					TemplateID: 1,
					Name:       "testParameter12",
					ValueString: sql.NullString{
						String: "TestParameter12",
						Valid:  true,
					},
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

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "designs", "present", nil), design)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestDeleteDesign_1.json"), &models.Design{})

	responseText, code = Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestDeleteDesign_2.json"), &models.Design{})
}
