// +build integration

package integration

import (
	"github.com/bouk/monkey"
	_ "github.com/qb0C80aE/clay/buildtime"
	"github.com/qb0C80aE/clay/model"
	"net/http"
	"testing"
	"time"
)

func TestGetDesign_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestGetDesign_Empty_1.json"), &model.Design{})
}

func TestGetDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template1 := &model.Template{
		ID:              1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
		Description:     "tedst1desc",
	}

	template2 := &model.Template{
		ID:              2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		Description:     "tedst2desc",
		TemplateArguments: []*model.TemplateArgument{
			{
				Name:         "testParameter1",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter1",
				Description:  "testParameter1desc",
			},
		},
	}

	templateArgument22 := &model.TemplateArgument{
		TemplateID:   2,
		Name:         "testParameter2",
		Type:         model.TemplateArgumentTypeString,
		DefaultValue: "TestParameter2",
		Description:  "testParameter2desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templateArgument22)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestGetDesign_1.json"), &model.Design{})

	parameters := map[string]string{
		"timestamp": "",
	}
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestGetDesign_2.json"), &model.Design{})
}

func TestUpdateDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	design := &model.Design{
		Content: map[string]interface{}{
			"template_persistent_parameters": []*model.TemplateArgument{
				{
					ID:           1,
					TemplateID:   1,
					Name:         "testParameter11",
					Type:         model.TemplateArgumentTypeString,
					DefaultValue: "TestParameter11",
					Description:  "testParameter11desc",
				},
				{
					ID:           2,
					TemplateID:   1,
					Name:         "testParameter12",
					Type:         model.TemplateArgumentTypeString,
					DefaultValue: "TestParameter12",
					Description:  "testParameter12desc",
				},
			},
			"templates": []*model.Template{
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
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestUpdateDesign_1.json"), &model.Design{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestUpdateDesign_2.json"), &model.Design{})
}

func TestDeleteDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	design := &model.Design{
		Content: map[string]interface{}{
			"template_persistent_parameters": []*model.TemplateArgument{
				{
					ID:           1,
					TemplateID:   1,
					Name:         "testParameter11",
					Type:         model.TemplateArgumentTypeString,
					DefaultValue: "TestParameter11",
					Description:  "testParameter11desc",
				},
				{
					ID:           2,
					TemplateID:   1,
					Name:         "testParameter12",
					Type:         model.TemplateArgumentTypeString,
					DefaultValue: "TestParameter12",
					Description:  "testParameter12desc",
				},
			},
			"templates": []*model.Template{
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
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestDeleteDesign_1.json"), &model.Design{})

	responseText, code = Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "designs", "present", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "design/TestDeleteDesign_2.json"), &model.Design{})
}

func init() {
	wayback := time.Date(1974, time.May, 19, 1, 2, 3, 4, time.UTC)
	monkey.Patch(time.Now, func() time.Time { return wayback })
}
