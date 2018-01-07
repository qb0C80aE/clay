// +build integration

package integration

import (
	"database/sql"
	"fmt"
	"github.com/qb0C80aE/clay/models"
	"net/http"
	"strconv"
	"testing"
)

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
		"preloads": "TemplatePersistentParameters",
	}

	template1 := &models.Template{
		Name:            "test1",
		TemplateContent: "TestTemplate1",
		Description:     "test1desc",
	}

	template2 := &models.Template{
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		Description:     "test2desc",
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

	template3 := &models.Template{
		ID:              100,
		Name:            "test100",
		TemplateContent: "TestTemplate100",
		Description:     "test100desc",
		TemplatePersistentParameters: []*models.TemplatePersistentParameter{
			{
				Name: "testParameter100",
				ValueString: sql.NullString{
					String: "TestParameter100",
					Valid:  true,
				},
				Description: "testParameter100desc",
			},
			{
				ID:   10,
				Name: "testParameter110",
				ValueString: sql.NullString{
					String: "TestParameter110",
					Valid:  true,
				},
				Description: "testParameter110desc",
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
		"preloads": "TemplatePersistentParameters",
	}

	id1 := 101
	template1 := &models.Template{
		ID:              id1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
		Description:     "test1desc",
	}

	id2 := 102
	template2 := &models.Template{
		ID:              id2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		Description:     "test2desc",
		TemplatePersistentParameters: []*models.TemplatePersistentParameter{
			{
				Name: "testParameter21",
				ValueString: sql.NullString{
					String: "TestParameter21",
					Valid:  true,
				},
				Description: "testParameter21desc",
			},
		},
	}

	id3 := 103
	template3 := &models.Template{
		ID:              id3,
		Name:            "test3",
		TemplateContent: "TestTemplate3",
		Description:     "test3desc",
		TemplatePersistentParameters: []*models.TemplatePersistentParameter{
			{
				ID:   1000,
				Name: "testParameter31",
				ValueString: sql.NullString{
					String: "TestParameter31",
					Valid:  true,
				},
				Description: "testParameter31desc",
			},
			{
				ID:   1001,
				Name: "testParameter32",
				ValueString: sql.NullString{
					String: "TestParameter32",
					Valid:  true,
				},
				Description: "testParameter31desc",
			},
		},
	}

	id4 := 104
	template4 := &models.Template{
		ID:              id4,
		Name:            "test4",
		TemplateContent: "TestTemplate4",
		Description:     "test4desc",
		TemplatePersistentParameters: []*models.TemplatePersistentParameter{
			{
				ID:   2000,
				Name: "testParameter41",
				ValueString: sql.NullString{
					String: "TestParameter41",
					Valid:  true,
				},
				Description: "testParameter41desc",
			},
			{
				ID:   2001,
				Name: "testParameter42",
				ValueString: sql.NullString{
					String: "TestParameter42",
					Valid:  true,
				},
				Description: "testParameter42desc",
			},
			{
				ID:   2002,
				Name: "testParameter43",
				ValueString: sql.NullString{
					String: "TestParameter43",
					Valid:  true,
				},
				Description: "testParameter43desc",
			},
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template3)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template4)

	template1.Name = "test1Updated"
	template1.TemplateContent = "TestTemplate1Updated"

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id1), nil), template1)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_1.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id1), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_2.json"), &models.Template{})

	template2.TemplatePersistentParameters = nil

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id2), nil), template2)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_3.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id2), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_4.json"), &models.Template{})

	template3.TemplatePersistentParameters[1].Name = "testParameter32Updated"
	template3.TemplatePersistentParameters[1].ValueString = sql.NullString{
		String: "TestParameter32Updated",
		Valid:  true,
	}
	template3.TemplatePersistentParameters[1].Description = "testParameter32descUpdated"

	template3.TemplatePersistentParameters = append(
		template3.TemplatePersistentParameters,
		&models.TemplatePersistentParameter{
			ID:   1002,
			Name: "testParameter34",
			ValueString: sql.NullString{
				String: "TestParameter34",
				Valid:  true,
			},
		},
	)

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id3), nil), template3)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_5.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id3), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_6.json"), &models.Template{})

	template3.TemplatePersistentParameters[0].ToBeDeleted = true
	template4.TemplatePersistentParameters[0].ToBeDeleted = true
	template4.TemplatePersistentParameters[2].ToBeDeleted = true
	template4.TemplatePersistentParameters = append(template4.TemplatePersistentParameters,
		[]*models.TemplatePersistentParameter{
			{
				ID:   2003,
				Name: "testParameter44",
				ValueString: sql.NullString{
					String: "TestParameter44",
					Valid:  true,
				},
				Description: "testParameter44desc",
			},
			{
				Name: "testParameter45",
				ValueString: sql.NullString{
					String: "TestParameter45",
					Valid:  true,
				},
				Description: "testParameter45desc",
			},
		}...,
	)

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id3), nil), template3)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_7.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id4), nil), template4)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_8.json"), &models.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "templates", parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_9.json"), []*models.Template{})

}

func TestDeleteTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &models.Template{
		ID:              id,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	responseText, code := Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplate_1.json"), &ErrorResponseText{})
}

func TestGetTemplatePersistentParameters_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, EmptyArrayString, []*models.TemplatePersistentParameter{})
}

func TestCreateTemplatePersistentParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	templatePersistentParameter1 := &models.TemplatePersistentParameter{
		TemplateID: 1,
		Name:       "testParameter1",
		ValueInt: sql.NullInt64{
			Int64: 100,
			Valid: true,
		},
		ValueFloat: sql.NullFloat64{
			Float64: 123.456,
			Valid:   true,
		},
		ValueBool: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		ValueString: sql.NullString{
			String: "TestParameter1",
			Valid:  true,
		},
		Description: "testParameter1desc",
	}
	templatePersistentParameter2 := &models.TemplatePersistentParameter{
		TemplateID: 1,
		Name:       "testParameter2",
		ValueInt: sql.NullInt64{
			Int64: 200,
			Valid: true,
		},
		ValueFloat: sql.NullFloat64{
			Float64: 456.789,
			Valid:   true,
		},
		ValueBool: sql.NullBool{
			Bool:  false,
			Valid: true,
		},
		ValueString: sql.NullString{
			String: "TestParameter2",
			Valid:  true,
		},
		Description: "testParameter2desc",
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter1)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplatePersistentParameter_1.json"), &models.TemplatePersistentParameter{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter2)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplatePersistentParameter_2.json"), &models.TemplatePersistentParameter{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplatePersistentParameter_3.json"), []*models.TemplatePersistentParameter{})
}

func TestUpdateTemplatePersistentParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templatePersistentParameter1 := &models.TemplatePersistentParameter{
		TemplateID: id,
		Name:       "testParameter1",
		ValueInt: sql.NullInt64{
			Int64: 100,
			Valid: true,
		},
		ValueFloat: sql.NullFloat64{
			Float64: 123.456,
			Valid:   true,
		},
		ValueBool: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		ValueString: sql.NullString{
			String: "TestParameter1",
			Valid:  true,
		},
		Description: "testParameter1desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter1)

	templatePersistentParameter1.Name = "templatePersistentParameter1Updated"
	templatePersistentParameter1.ValueInt = sql.NullInt64{
		Int64: 999,
		Valid: true,
	}
	templatePersistentParameter1.ValueString = sql.NullString{
		String: "TestParameter1Updated",
		Valid:  true,
	}
	templatePersistentParameter1.Description = "testParameter1descUpdated"

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "template_persistent_parameters", strconv.Itoa(id), nil), templatePersistentParameter1)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplatePersistentParameter_1.json"), &models.TemplatePersistentParameter{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_persistent_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplatePersistentParameter_2.json"), &models.TemplatePersistentParameter{})
}

func TestDeleteTemplatePersistentParameters(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &models.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templatePersistentParameter1 := &models.TemplatePersistentParameter{
		TemplateID: id,
		Name:       "testParameter1",
		ValueString: sql.NullString{
			String: "TestParameter1",
			Valid:  true,
		},
		Description: "testParameter1desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter1)

	templatePersistentParameter1.Name = "templatePersistentParameter1Updated"
	templatePersistentParameter1.ValueString = sql.NullString{
		String: "TestParameter1Updated",
		Valid:  true,
	}
	templatePersistentParameter1.Description = "TestParameter1descUpdated"

	responseText, code := Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "template_persistent_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_persistent_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplatePersistentParameter_1.json"), &ErrorResponseText{})
}

func TestDeleteTemplatePersistentParameters_Cascade(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	templateID := 1
	template := &models.Template{
		ID:              templateID,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templatePersistentParameter1 := &models.TemplatePersistentParameter{
		TemplateID: id,
		Name:       "testParameter1",
		ValueString: sql.NullString{
			String: "TestParameter1",
			Valid:  true,
		},
		Description: "testParameter1desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter1)

	Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "templates", strconv.Itoa(templateID), nil), nil)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_persistent_parameters", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplatePersistentParameter_Cascade_1.json"), &ErrorResponseText{})
}

func TestTemplate_ExtractFromDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template1 := &models.Template{
		ID:              1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
		Description:     "test1desc",
	}
	template2 := &models.Template{
		ID:              2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		Description:     "test2desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)

	templatePersistentParameter11 := &models.TemplatePersistentParameter{
		TemplateID: 1,
		Name:       "testParameter11",
		ValueString: sql.NullString{
			String: "TestParameter11",
			Valid:  true,
		},
		Description: "testParameter11desc",
	}
	templatePersistentParameter12 := &models.TemplatePersistentParameter{
		TemplateID: 1,
		Name:       "testParameter12",
		ValueString: sql.NullString{
			String: "TestParameter12",
			Valid:  true,
		},
		Description: "testParameter12desc",
	}
	templatePersistentParameter21 := &models.TemplatePersistentParameter{
		TemplateID: 2,
		Name:       "testParameter21",
		ValueString: sql.NullString{
			String: "TestParameter21",
			Valid:  true,
		},
		Description: "testParameter21desc",
	}
	templatePersistentParameter22 := &models.TemplatePersistentParameter{
		TemplateID: 2,
		Name:       "testParameter22",
		ValueString: sql.NullString{
			String: "TestParameter22",
			Valid:  true,
		},
		Description: "testParameter22desc",
	}
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter11)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter12)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter21)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter22)
}

func TestTemplate_Generate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &models.Template{
		ID:   id,
		Name: "test",
		TemplateContent: "" +
			`{{.testParameter1.ValueString.String}} is TestParameter1, {{.testParameter2.ValueString.String}} is TestParameter2
urlParam1 is {{.urlParam1}}
urlParam2 is {{.urlParam2}}
urlParam3 is {{.urlParam3}}
{{index .urlParam1 0}}
{{index .urlParam2 0}}
{{index .urlParam2 1}}
{{index .urlParam3 0}}
{{index .urlParam3 1}}
{{index .urlParam3 2}}`,
		TemplatePersistentParameters: []*models.TemplatePersistentParameter{
			{
				Name: "testParameter1",
				ValueString: sql.NullString{
					String: "TestParameter1",
					Valid:  true,
				},
			},
			{
				Name: "testParameter2",
				ValueString: sql.NullString{
					String: "TestParameter2",
					Valid:  true,
				},
			},
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	parameters := map[string]string{
		"urlParam1": "100",
		"urlParam2": "200&urlParam2=201",
		"urlParam3": "300&urlParam3=301&&urlParam3=302",
	}
	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "generation", parameters), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_1.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "raw", parameters), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_2.txt"))
}

func TestTemplate_FuncMaps(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 1
	template1 := &models.Template{
		ID:   id,
		Name: "test1",
		TemplateContent: `--- calc ---
i = 100
{{- $i := 100}}

i = i + 2
{{- $i := add $i 2}}
{{$i}}

i = i - 4
{{- $i := sub $i 4}}
{{$i}}

i = i * 6
{{- $i := mul $i 6}}
{{$i}}

i = i / 2
{{- $i := div $i 2}}
{{$i}}

i = i mod 5
{{- $i := mod $i 5}}
{{$i}}

--- conversion ---
{{- $vint := int "100" }}
{{- $vint8 := int8 "101" }}
{{- $vint16 := int16 "102" }}
{{- $vint32 := int32 "103" }}
{{- $vint64 := int64 "104" }}
{{- $vuint := uint "200" }}
{{- $vuint8 := uint8 "201" }}
{{- $vuint16 := uint16 "202" }}
{{- $vuint32 := uint32 "203" }}
{{- $vuint64 := uint64 "204" }}
{{- $vfloat32 := float32 "300.1" }}
{{- $vfloat64 := float64 "300.2" }}
{{- $vboolean := boolean "false" }}
{{- $vmap := map "key" "value" }}
{{ $vint }}
{{ $vint8 }}
{{ $vint16 }}
{{ $vint32 }}
{{ $vint64 }}
{{ $vuint }}
{{ $vuint8 }}
{{ $vuint16 }}
{{ $vuint32 }}
{{ $vuint64 }}
{{ $vfloat32 }}
{{ $vfloat64 }}
{{ $vboolean }}
{{ $vmap }}
{{- $sint := string $vint }}
{{- $sint8 := string $vint8 }}
{{- $sint16 := string $vint16 }}
{{- $sint32 := string $vint32 }}
{{- $sint64 := string $vint64 }}
{{- $suint := string $vuint }}
{{- $suint8 := string $vuint8 }}
{{- $suint16 := string $vuint16 }}
{{- $suint32 := string $vuint32 }}
{{- $suint64 := string $vuint64 }}
{{- $sfloat32:= string $vfloat32 }}
{{- $sfloat64 := string $vfloat64 }}
{{- $sboolean := string $vboolean }}
{{- $sobject := string $vmap }}
{{ $sint }}
{{ $sint8 }}
{{ $sint16 }}
{{ $sint32 }}
{{ $sint64 }}
{{ $suint }}
{{ $suint8 }}
{{ $suint16 }}
{{ $suint32 }}
{{ $suint64 }}
{{ $sfloat32 }}
{{ $sfloat64 }}
{{ $sboolean }}
{{ $sobject }}
{{ int $vint }}
{{ int $vint8 }}
{{ int $vint16 }}
{{ int $vint32 }}
{{ int $vint64 }}
{{ int $vuint }}
{{ int $vuint8 }}
{{ int $vuint16 }}
{{ int $vuint32 }}
{{ int $vuint64 }}
{{ int $vfloat32 }}
{{ int $vfloat64 }}
{{ int8 $vint }}
{{ int8 $vint8 }}
{{ int8 $vint16 }}
{{ int8 $vint32 }}
{{ int8 $vint64 }}
{{ int8 $vuint }}
{{ int8 $vuint8 }}
{{ int8 $vuint16 }}
{{ int8 $vuint32 }}
{{ int8 $vuint64 }}
{{ int8 $vfloat32 }}
{{ int8 $vfloat64 }}
{{ int16 $vint }}
{{ int16 $vint8 }}
{{ int16 $vint16 }}
{{ int16 $vint32 }}
{{ int16 $vint64 }}
{{ int16 $vuint }}
{{ int16 $vuint8 }}
{{ int16 $vuint16 }}
{{ int16 $vuint32 }}
{{ int16 $vuint64 }}
{{ int16 $vfloat32 }}
{{ int16 $vfloat64 }}
{{ int32 $vint }}
{{ int32 $vint8 }}
{{ int32 $vint16 }}
{{ int32 $vint32 }}
{{ int32 $vint64 }}
{{ int32 $vuint }}
{{ int32 $vuint8 }}
{{ int32 $vuint16 }}
{{ int32 $vuint32 }}
{{ int32 $vuint64 }}
{{ int32 $vfloat32 }}
{{ int32 $vfloat64 }}
{{ int64 $vint }}
{{ int64 $vint8 }}
{{ int64 $vint16 }}
{{ int64 $vint32 }}
{{ int64 $vint64 }}
{{ int64 $vuint }}
{{ int64 $vuint8 }}
{{ int64 $vuint16 }}
{{ int64 $vuint32 }}
{{ int64 $vuint64 }}
{{ int64 $vfloat32 }}
{{ int64 $vfloat64 }}
{{ uint $vint }}
{{ uint $vint8 }}
{{ uint $vint16 }}
{{ uint $vint32 }}
{{ uint $vint64 }}
{{ uint $vuint }}
{{ uint $vuint8 }}
{{ uint $vuint16 }}
{{ uint $vuint32 }}
{{ uint $vuint64 }}
{{ uint $vfloat32 }}
{{ uint $vfloat64 }}
{{ uint8 $vint }}
{{ uint8 $vint8 }}
{{ uint8 $vint16 }}
{{ uint8 $vint32 }}
{{ uint8 $vint64 }}
{{ uint8 $vuint }}
{{ uint8 $vuint8 }}
{{ uint8 $vuint16 }}
{{ uint8 $vuint32 }}
{{ uint8 $vuint64 }}
{{ uint8 $vfloat32 }}
{{ uint8 $vfloat64 }}
{{ uint16 $vint }}
{{ uint16 $vint8 }}
{{ uint16 $vint16 }}
{{ uint16 $vint32 }}
{{ uint16 $vint64 }}
{{ uint16 $vuint }}
{{ uint16 $vuint8 }}
{{ uint16 $vuint16 }}
{{ uint16 $vuint32 }}
{{ uint16 $vuint64 }}
{{ uint16 $vfloat32 }}
{{ uint16 $vfloat64 }}
{{ uint32 $vint }}
{{ uint32 $vint8 }}
{{ uint32 $vint16 }}
{{ uint32 $vint32 }}
{{ uint32 $vint64 }}
{{ uint32 $vuint }}
{{ uint32 $vuint8 }}
{{ uint32 $vuint16 }}
{{ uint32 $vuint32 }}
{{ uint32 $vuint64 }}
{{ uint32 $vfloat32 }}
{{ uint32 $vfloat64 }}
{{ uint64 $vint }}
{{ uint64 $vint8 }}
{{ uint64 $vint16 }}
{{ uint64 $vint32 }}
{{ uint64 $vint64 }}
{{ uint64 $vuint }}
{{ uint64 $vuint8 }}
{{ uint64 $vuint16 }}
{{ uint64 $vuint32 }}
{{ uint64 $vuint64 }}
{{ uint64 $vfloat32 }}
{{ uint64 $vfloat64 }}
{{ float32 $vint }}
{{ float32 $vint8 }}
{{ float32 $vint16 }}
{{ float32 $vint32 }}
{{ float32 $vint64 }}
{{ float32 $vuint }}
{{ float32 $vuint8 }}
{{ float32 $vuint16 }}
{{ float32 $vuint32 }}
{{ float32 $vuint64 }}
{{ float32 $vfloat32 }}
{{ float32 $vfloat64 }}
{{ float64 $vint }}
{{ float64 $vint8 }}
{{ float64 $vint16 }}
{{ float64 $vint32 }}
{{ float64 $vint64 }}
{{ float64 $vuint }}
{{ float64 $vuint8 }}
{{ float64 $vuint16 }}
{{ float64 $vuint32 }}
{{ float64 $vuint64 }}
{{ float64 $vfloat32 }}
{{ float64 $vfloat64 }}

--- slice ---
slice init
{{- $slice1 := slice}}
slice1: {{$slice1}}

slice init
{{- $slice2 := slice 1 2 3 4 5}}
slice2: {{$slice2}}

slice init
{{- $slice3 := slice 3 4 5 6 7}}
slice3: {{$slice3}}

subslice
{{- $subslice := subslice $slice3 -1 3}}
{{$subslice}}

subslice
{{subslice $slice3 1 -1}}

subslice
{{subslice $slice3 -1 -1}}

subslice
{{subslice $slice3 1 3}}

append
{{- $slice3 := append $slice3 8 9 10}}
{{$slice3}}

concatenate
{{concatenate $slice2 $slice3}}

fieldslice
{{- $tpp := multi .ModelStore "template_persistent_parameters" ""}}
{{- $idfields := fieldslice $tpp "ID"}}
{{- $namefields := fieldslice $tpp "Name"}}
{{- $idfields := sort $idfields "asc"}}
{{- $namefields := sort $namefields "asc"}}
{{$idfields}}
{{$namefields}}

sort
{{- $sliceint := slice}}
{{- $v1 := int 3}}
{{- $v2 := int 1}}
{{- $v3 := int -5}}
{{- $v4 := int 2}}
{{- $v5 := int 4}}
{{- $sliceint := append $sliceint $v1 $v2 $v3 $v4 $v5}}
{{- $sliceint := sort $sliceint "asc" }}
sliceint asc: {{$sliceint}}
{{- $sliceint := sort $sliceint "desc" }}
sliceint desc: {{$sliceint}}
{{- $sliceint8 := slice}}
{{- $v1 := int8 3}}
{{- $v2 := int8 1}}
{{- $v3 := int8 -5}}
{{- $v4 := int8 2}}
{{- $v5 := int8 4}}
{{- $sliceint8 := append $sliceint8 $v1 $v2 $v3 $v4 $v5}}
{{- $sliceint8 := sort $sliceint8 "asc" }}
sliceint8 asc: {{$sliceint8}}
{{- $sliceint8 := sort $sliceint8 "desc" }}
sliceint8 desc: {{$sliceint8}}
{{- $sliceint16 := slice}}
{{- $v1 := int16 3}}
{{- $v2 := int16 1}}
{{- $v3 := int16 -5}}
{{- $v4 := int16 2}}
{{- $v5 := int16 4}}
{{- $sliceint16 := append $sliceint16 $v1 $v2 $v3 $v4 $v5}}
{{- $sliceint16 := sort $sliceint16 "asc" }}
sliceint16 asc: {{$sliceint16}}
{{- $sliceint16 := sort $sliceint16 "desc" }}
sliceint16 desc: {{$sliceint16}}
{{- $sliceint32 := slice}}
{{- $v1 := int32 3}}
{{- $v2 := int32 1}}
{{- $v3 := int32 -5}}
{{- $v4 := int32 2}}
{{- $v5 := int32 4}}
{{- $sliceint32 := append $sliceint32 $v1 $v2 $v3 $v4 $v5}}
{{- $sliceint32 := sort $sliceint32 "asc" }}
sliceint32 asc: {{$sliceint32}}
{{- $sliceint32 := sort $sliceint32 "desc" }}
sliceint32 desc: {{$sliceint32}}
{{- $sliceint64 := slice}}
{{- $v1 := int64 3}}
{{- $v2 := int64 1}}
{{- $v3 := int64 -5}}
{{- $v4 := int64 2}}
{{- $v5 := int64 4}}
{{- $sliceint64 := append $sliceint64 $v1 $v2 $v3 $v4 $v5}}
{{- $sliceint64 := sort $sliceint64 "asc" }}
sliceint64 asc: {{$sliceint64}}
{{- $sliceint64 := sort $sliceint64 "desc" }}
sliceint64 desc: {{$sliceint64}}
{{- $sliceuint := slice}}
{{- $v1 := uint 3}}
{{- $v2 := uint 1}}
{{- $v3 := uint 5}}
{{- $v4 := uint 2}}
{{- $v5 := uint 4}}
{{- $sliceuint := append $sliceuint $v1 $v2 $v3 $v4 $v5}}
{{- $sliceuint := sort $sliceuint "asc" }}
sliceuint asc: {{$sliceuint}}
{{- $sliceuint := sort $sliceuint "desc" }}
sliceuint desc: {{$sliceuint}}
{{- $sliceuint8 := slice}}
{{- $v1 := uint8 3}}
{{- $v2 := uint8 1}}
{{- $v3 := uint8 5}}
{{- $v4 := uint8 2}}
{{- $v5 := uint8 4}}
{{- $sliceuint8 := append $sliceuint8 $v1 $v2 $v3 $v4 $v5}}
{{- $sliceuint8 := sort $sliceuint8 "asc" }}
sliceuint8 asc: {{$sliceuint8}}
{{- $sliceuint8 := sort $sliceuint8 "desc" }}
sliceuint8 desc: {{$sliceuint8}}
{{- $sliceuint16 := slice}}
{{- $v1 := uint16 3}}
{{- $v2 := uint16 1}}
{{- $v3 := uint16 5}}
{{- $v4 := uint16 2}}
{{- $v5 := uint16 4}}
{{- $sliceuint16 := append $sliceuint16 $v1 $v2 $v3 $v4 $v5}}
{{- $sliceuint16 := sort $sliceuint16 "asc" }}
sliceuint16 asc: {{$sliceuint16}}
{{- $sliceuint16 := sort $sliceuint16 "desc" }}
sliceuint16 desc: {{$sliceuint16}}
{{- $sliceuint32 := slice}}
{{- $v1 := uint32 3}}
{{- $v2 := uint32 1}}
{{- $v3 := uint32 5}}
{{- $v4 := uint32 2}}
{{- $v5 := uint32 4}}
{{- $sliceuint32 := append $sliceuint32 $v1 $v2 $v3 $v4 $v5}}
{{- $sliceuint32 := sort $sliceuint32 "asc" }}
sliceuint32 asc: {{$sliceuint32}}
{{- $sliceuint32 := sort $sliceuint32 "desc" }}
sliceuint32 desc: {{$sliceuint32}}
{{- $sliceuint64 := slice}}
{{- $v1 := uint64 3}}
{{- $v2 := uint64 1}}
{{- $v3 := uint64 5}}
{{- $v4 := uint64 2}}
{{- $v5 := uint64 4}}
{{- $sliceuint64 := append $sliceuint64 $v1 $v2 $v3 $v4 $v5}}
{{- $sliceuint64 := sort $sliceuint64 "asc" }}
sliceuint64 asc: {{$sliceuint64}}
{{- $sliceuint64 := sort $sliceuint64 "desc" }}
sliceuint64 desc: {{$sliceuint64}}
{{- $slicefloat32 := slice}}
{{- $v1 := float32 3.3}}
{{- $v2 := float32 1}}
{{- $v3 := float32 -5.1}}
{{- $v4 := float32 2.2}}
{{- $v5 := float32 4}}
{{- $slicefloat32 := append $slicefloat32 $v1 $v2 $v3 $v4 $v5}}
{{- $slicefloat32 := sort $slicefloat32 "asc" }}
slicefloat32 asc: {{$slicefloat32}}
{{- $slicefloat32 := sort $slicefloat32 "desc" }}
slicefloat32 desc: {{$slicefloat32}}
{{- $slicefloat64 := slice}}
{{- $v1 := float64 3.3}}
{{- $v2 := float64 1}}
{{- $v3 := float64 -5.1}}
{{- $v4 := float64 2.2}}
{{- $v5 := float64 4}}
{{- $slicefloat64 := append $slicefloat64 $v1 $v2 $v3 $v4 $v5}}
{{- $slicefloat64 := sort $slicefloat64 "asc" }}
slicefloat64 asc: {{$slicefloat64}}
{{- $slicefloat64 := sort $slicefloat64 "desc" }}
slicefloat64 desc: {{$slicefloat64}}
{{- $slicestring := slice}}
{{- $v1 := "3.3"}}
{{- $v2 := "ABC"}}
{{- $v3 := "-5.1"}}
{{- $v4 := "012"}}
{{- $v5 := "def"}}
{{- $slicestring := append $slicestring $v1 $v2 $v3 $v4 $v5}}
{{- $slicestring := sort $slicestring "asc" }}
slicestring asc: {{$slicestring}}
{{- $slicestring := sort $slicestring "desc" }}
slicestring desc: {{$slicestring}}

--- map ---
map
{{- $map1 := map}}
map1: {{$map1}}

map init get
{{- $map2 := map 1 "A" 2 "B"}}
map2[1]: {{get $map2 1}}
map2[2]: {{get $map2 2}}

map init get
{{- $map3 := map 1 "C" 3 "D"}}
map3[1]: {{get $map3 1}}
map3[3]: {{get $map3 3}}

map exists
{{exists $map2 0}}
{{- $e := exists $map2 1}}
{{if eq $e true}}TRUE!!{{else}}FALSE!!{{end}}

map put
{{- $null := put $map1 4 "E"}}
{{- $null := put $map1 5 "F"}}
{{- $null := put $map1 6 "G"}}
map1[4]: {{get $map1 4}}
map1[5]: {{get $map1 5}}
map1[6]: {{get $map1 6}}

map delete
{{- $null := delete $map1 5}}
map1[3]: {{get $map1 3}}
map1[4]: {{get $map1 4}}
map1[5]: {{get $map1 5}}
map1[6]: {{get $map1 6}}
map1[7]: {{get $map1 7}}

map merge
{{- $null := merge $map2 $map1}}
{{- $null := merge $map3 $map1}}
map1[0]: {{get $map1 0}}
map1[1]: {{get $map1 1}}
map1[2]: {{get $map1 2}}
map1[3]: {{get $map1 3}}
map1[4]: {{get $map1 4}}
map1[5]: {{get $map1 5}}
map1[6]: {{get $map1 6}}
map1[7]: {{get $map1 7}}

map keys
{{- $keys := keys $map1}}
{{- $keys := sort $keys "asc"}}
keys of map1: {{$keys}}

--- model store ---
multi
{{- $m := multi .ModelStore "templates" "preloads=TemplatePersistentParameters"}}
{{- $t := index $m 0}}
{{$t.Name}}
{{- $p1 := index $t.TemplatePersistentParameters 0}}
{{$p1.Name}}={{$p1.ValueString.String}}
{{- $p2 := index $t.TemplatePersistentParameters 1}}
{{$p2.Name}}={{$p2.ValueString.String}}

single
{{- $path := printf "/templates/%d" .testParameter11.ValueInt.Int64}}
{{- $s := single .ModelStore $path "preloads=TemplatePersistentParameters"}}
{{$s.Name}}
{{- $p1 := index $s.TemplatePersistentParameters 0}}
{{$p1.Name}}={{$p1.ValueString.String}}
{{- $p2 := index $s.TemplatePersistentParameters 1}}
{{$p2.Name}}={{$p2.ValueString.String}}

total
{{- $t := total .ModelStore "/template_persistent_parameters"}}
{{$t}}

--- hash ---
hash
{{- $h := hash $s.TemplatePersistentParameters "Name"}}
hash[testParameter11]={{get $h "testParameter11"}}
hash[testParameter12]={{get $h "testParameter12"}}

--- slicemap ---
slicemap
{{- $p := multi .ModelStore "template_persistent_parameters" ""}}
{{- $z := slicemap $p "Name"}}
{{- $z1 := get $z "testParameter11"}}
{{- $z2 := get $z "testParameter12"}}
{{- $z3 := get $z "testParameter1X"}}
{{- range $i, $v := $z1 }}
slicemap[testParameter1][{{$i}}]={{$v}}
{{- end}}
{{- range $i, $v := $z2 }}
slicemap[testParameter2][{{$i}}]={{$v}}
{{- end}}
{{- range $i, $v := $z3 }}
slicemap[testParameter1X][{{$i}}]={{$v}}
{{- end}}

--- sequence ---
sequence
{{- $s := sequence 1 10}}
{{- range $i, $v := $s}}
sequence[{{$i}}]={{$v}}
{{- end}}
`,
		Description: "test1desc",
	}

	templatePersistentParameter11 := &models.TemplatePersistentParameter{
		TemplateID: id,
		Name:       "testParameter11",
		ValueString: sql.NullString{
			String: "TestParameter11",
			Valid:  true,
		},
		ValueInt: sql.NullInt64{
			Int64: 1,
			Valid: true,
		},
		Description: "testParameter11desc",
	}
	templatePersistentParameter12 := &models.TemplatePersistentParameter{
		TemplateID: id,
		Name:       "testParameter12",
		ValueString: sql.NullString{
			String: "TestParameter12",
			Valid:  true,
		},
		Description: "testParameter12desc",
	}

	templatePersistentParameter13 := &models.TemplatePersistentParameter{
		TemplateID: id,
		Name:       "testParameter1X",
		ValueInt: sql.NullInt64{
			Int64: 100,
			Valid: true,
		},
		Description: "testParameter1Xdesc",
	}

	id2 := 2
	template2 := &models.Template{
		ID:              id2,
		Name:            "test12",
		TemplateContent: `{{.testParameter1X}}`,
		Description:     "test12desc",
	}
	templatePersistentParameter21 := &models.TemplatePersistentParameter{
		TemplateID: id2,
		Name:       "testParameter1X",
		ValueInt: sql.NullInt64{
			Int64: 200,
			Valid: true,
		},
		Description: "testParameter1Xdesc",
	}

	id3 := 3
	template3 := &models.Template{
		ID:              id3,
		Name:            "test13",
		TemplateContent: `{{.testParameter1X}}`,
		Description:     "test13desc",
	}
	templatePersistentParameter31 := &models.TemplatePersistentParameter{
		TemplateID: id3,
		Name:       "testParameter1X",
		ValueInt: sql.NullInt64{
			Int64: 300,
			Valid: true,
		},
		Description: "testParameter1Xdesc",
	}

	id4 := 4
	template4 := &models.Template{
		ID:              id4,
		Name:            "test14",
		TemplateContent: `volatile1 = {{.volatile1}}, volatile2 = {{.volatile2}}`,
		Description:     "test14desc",
	}

	id5 := 5
	template5 := &models.Template{
		ID:   id5,
		Name: "test15",
		TemplateContent: `include test
{{- $volatile := map }}
{{- $volatile := put $volatile "volatile1" "abc" }}
{{- $volatile := put $volatile "volatile2" 123 }}
{{ include .ModelStore "test12" nil }}
{{ include .ModelStore "test13" nil }}
{{ include .ModelStore "test14" $volatile }}
`,
		Description: "test15desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter11)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter12)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter13)

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter21)

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template3)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_persistent_parameters", nil), templatePersistentParameter31)

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template4)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template5)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_FuncMaps_1.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id2), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_FuncMaps_2.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id3), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_FuncMaps_3.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id5), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_FuncMaps_4.txt"))
}
