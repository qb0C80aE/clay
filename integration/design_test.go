// +build integration

package integration

import (
	"database/sql"
	"github.com/bouk/monkey"
	_ "github.com/qb0C80aE/clay/buildtime"
	"github.com/qb0C80aE/clay/models"
	"net/http"
	"testing"
	"time"
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
		Description:     "tedst1desc",
	}

	template2 := &models.Template{
		ID:              2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		Description:     "tedst2desc",
		TemplatePersistentParameters: []*models.TemplatePersistentParameter{
			{
				Name: "testParameter1",
				ValueString: sql.NullString{
					String: "TestParameter1",
					Valid:  true,
				},
				Description: "testParameter1desc",
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
		Description: "testParameter2desc",
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
					Description: "testParameter11desc",
				},
				{
					ID:         2,
					TemplateID: 1,
					Name:       "testParameter12",
					ValueString: sql.NullString{
						String: "TestParameter12",
						Valid:  true,
					},
					Description: "testParameter12desc",
				},
			},
			"templates": []*models.Template{
				{
					ID:              1,
					Name:            "test1",
					TemplateContent: "TestTemplate1",
					Description:     "test1desc",
				},
				{
					ID:              2,
					Name:            "test2",
					TemplateContent: "TestTemplate2",
					Description:     "test2desc",
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
					Description: "testParameter11desc",
				},
				{
					ID:         2,
					TemplateID: 1,
					Name:       "testParameter12",
					ValueString: sql.NullString{
						String: "TestParameter12",
						Valid:  true,
					},
					Description: "testParameter12desc",
				},
			},
			"templates": []*models.Template{
				{
					ID:              1,
					Name:            "test1",
					TemplateContent: "TestTemplate1",
					Description:     "test1desc",
				},
				{
					ID:              2,
					Name:            "test2",
					TemplateContent: "TestTemplate2",
					Description:     "test2desc",
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

func init() {
	wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	monkey.Patch(time.Now, func() time.Time { return wayback })
}
