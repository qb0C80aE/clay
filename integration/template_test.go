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
				Name: "testParameter1",
				ValueString: sql.NullString{
					String: "TestParameter1",
					Valid:  true,
				},
			},
		},
	}

	template3 := &models.Template{
		ID:              100,
		Name:            "test100",
		TemplateContent: "TestTemplate100",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
			{
				Name: "testParameter100",
				ValueString: sql.NullString{
					String: "TestParameter100",
					Valid:  true,
				},
			},
			{
				ID:   10,
				Name: "testParameter110",
				ValueString: sql.NullString{
					String: "TestParameter110",
					Valid:  true,
				},
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
				Name: "testParameter21",
				ValueString: sql.NullString{
					String: "TestParameter21",
					Valid:  true,
				},
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
				ID:   1000,
				Name: "testParameter31",
				ValueString: sql.NullString{
					String: "TestParameter31",
					Valid:  true,
				},
			},
			{
				ID:   1001,
				Name: "testParameter32",
				ValueString: sql.NullString{
					String: "TestParameter32",
					Valid:  true,
				},
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
	template3.TemplateExternalParameters[1].ValueString = sql.NullString{
		String: "TestParameter32Updated",
		Valid:  true,
	}

	template3.TemplateExternalParameters = append(
		template3.TemplateExternalParameters,
		&models.TemplateExternalParameter{
			ID:   1003,
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
		TemplateContent: "{{.testParameter1.ValueString.String}} is TestParameter1, {{.testParameter2.ValueString.String}} is TestParameter2.",
		TemplateExternalParameters: []*models.TemplateExternalParameter{
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

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "generation", nil), nil)
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
	}
	templateExternalParameter2 := &models.TemplateExternalParameter{
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
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter1)

	templateExternalParameter1.Name = "templateExternalParameter1Updated"
	templateExternalParameter1.ValueInt = sql.NullInt64{
		Int64: 999,
		Valid: true,
	}
	templateExternalParameter1.ValueString = sql.NullString{
		String: "TestParameter1Updated",
		Valid:  true,
	}

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
		ValueString: sql.NullString{
			String: "TestParameter1",
			Valid:  true,
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter1)

	templateExternalParameter1.Name = "templateExternalParameter1Updated"
	templateExternalParameter1.ValueString = sql.NullString{
		String: "TestParameter1Updated",
		Valid:  true,
	}

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
		ValueString: sql.NullString{
			String: "TestParameter1",
			Valid:  true,
		},
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
		ValueString: sql.NullString{
			String: "TestParameter11",
			Valid:  true,
		},
	}
	templateExternalParameter12 := &models.TemplateExternalParameter{
		TemplateID: 1,
		Name:       "testParameter12",
		ValueString: sql.NullString{
			String: "TestParameter12",
			Valid:  true,
		},
	}
	templateExternalParameter21 := &models.TemplateExternalParameter{
		TemplateID: 2,
		Name:       "testParameter21",
		ValueString: sql.NullString{
			String: "TestParameter21",
			Valid:  true,
		},
	}
	templateExternalParameter22 := &models.TemplateExternalParameter{
		TemplateID: 2,
		Name:       "testParameter22",
		ValueString: sql.NullString{
			String: "TestParameter22",
			Valid:  true,
		},
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
		ValueString: sql.NullString{
			String: "TestParameter11",
			Valid:  true,
		},
	}
	templateExternalParameter12 := &models.TemplateExternalParameter{
		ID:         2,
		TemplateID: 1,
		Name:       "testParameter12",
		ValueString: sql.NullString{
			String: "TestParameter12",
			Valid:  true,
		},
	}
	templateExternalParameter21 := &models.TemplateExternalParameter{
		ID:         3,
		TemplateID: 2,
		Name:       "testParameter21",
		ValueString: sql.NullString{
			String: "TestParameter21",
			Valid:  true,
		},
	}
	templateExternalParameter22 := &models.TemplateExternalParameter{
		ID:         4,
		TemplateID: 2,
		Name:       "testParameter22",
		ValueString: sql.NullString{
			String: "TestParameter22",
			Valid:  true,
		},
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
		ValueString: sql.NullString{
			String: "TestParameter11",
			Valid:  true,
		},
	}
	templateExternalParameter12 := &models.TemplateExternalParameter{
		ID:         2,
		TemplateID: 1,
		Name:       "testParameter12",
		ValueString: sql.NullString{
			String: "TestParameter12",
			Valid:  true,
		},
	}
	templateExternalParameter21 := &models.TemplateExternalParameter{
		ID:         3,
		TemplateID: 2,
		Name:       "testParameter21",
		ValueString: sql.NullString{
			String: "TestParameter21",
			Valid:  true,
		},
	}
	templateExternalParameter22 := &models.TemplateExternalParameter{
		ID:         4,
		TemplateID: 2,
		Name:       "testParameter22",
		ValueString: sql.NullString{
			String: "TestParameter22",
			Valid:  true,
		},
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

--- model store ---
multi
{{- $m := multi .ModelStore "templates" "preloads=TemplateExternalParameters"}}
{{- $t := index $m 0}}
{{$t.Name}}
{{- $p1 := index $t.TemplateExternalParameters 0}}
{{$p1.Name}}={{$p1.ValueString.String}}
{{- $p2 := index $t.TemplateExternalParameters 1}}
{{$p2.Name}}={{$p2.ValueString.String}}

single
{{- $s := single .ModelStore "templates" .testParameter11.ValueInt.Int64 "preloads=TemplateExternalParameters"}}
{{$s.Name}}
{{- $p1 := index $s.TemplateExternalParameters 0}}
{{$p1.Name}}={{$p1.ValueString.String}}
{{- $p2 := index $s.TemplateExternalParameters 1}}
{{$p2.Name}}={{$p2.ValueString.String}}

total
{{- $t := total .ModelStore "template_external_parameters"}}
{{$t}}

--- hash ---
hash
{{- $h := hash $s.TemplateExternalParameters "Name"}}
hash[testParameter11]={{get $h "testParameter11"}}
hash[testParameter12]={{get $h "testParameter12"}}

--- slicemap ---
slicemap
{{- $p := multi .ModelStore "template_external_parameters" ""}}
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
	}

	templateExternalParameter11 := &models.TemplateExternalParameter{
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
	}
	templateExternalParameter12 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter12",
		ValueString: sql.NullString{
			String: "TestParameter12",
			Valid:  true,
		},
	}

	templateExternalParameter13 := &models.TemplateExternalParameter{
		TemplateID: id,
		Name:       "testParameter1X",
		ValueInt: sql.NullInt64{
			Int64: 100,
			Valid: true,
		},
	}

	id2 := 2
	template2 := &models.Template{
		ID:              id2,
		Name:            "test12",
		TemplateContent: `{{.testParameter1X}}`,
	}
	templateExternalParameter21 := &models.TemplateExternalParameter{
		TemplateID: id2,
		Name:       "testParameter1X",
		ValueInt: sql.NullInt64{
			Int64: 200,
			Valid: true,
		},
	}

	id3 := 3
	template3 := &models.Template{
		ID:              id3,
		Name:            "test13",
		TemplateContent: `{{.testParameter1X}}`,
	}
	templateExternalParameter31 := &models.TemplateExternalParameter{
		TemplateID: id3,
		Name:       "testParameter1X",
		ValueInt: sql.NullInt64{
			Int64: 300,
			Valid: true,
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter11)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter12)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter13)

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter21)

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template3)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_external_parameters", nil), templateExternalParameter31)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_FuncMaps_1.txt"))
}
