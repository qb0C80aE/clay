// +build integration

package integration

import (
	"fmt"
	"github.com/qb0C80aE/clay/model"
	"net/http"
	"strconv"
	"testing"
)

func TestGetTemplates_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "templates", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, EmptyArrayString, []*model.Template{})
}

func TestCreateTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	parameters := map[string]string{
		"preloads": "TemplateArguments",
	}

	template1 := &model.Template{
		Name:            "test1",
		TemplateContent: "TestTemplate1",
		Description:     "test1desc",
	}

	template2 := &model.Template{
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		Description:     "test2desc",
		TemplateArguments: []*model.TemplateArgument{
			{
				Name:         "testParameter1",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter1",
				Description:  "testParameter1desc",
			},
		},
	}

	template3 := &model.Template{
		ID:              100,
		Name:            "test100",
		TemplateContent: "TestTemplate100",
		Description:     "test100desc",
		TemplateArguments: []*model.TemplateArgument{
			{
				Name:         "testParameter100",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter100",
				Description:  "testParameter100desc",
			},
			{
				ID:           10,
				Name:         "testParameter110",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter110",
				Description:  "testParameter110desc",
			},
		},
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplate_1.json"), &model.Template{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplate_2.json"), &model.Template{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template3)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplate_3.json"), &model.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "templates", parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplate_4.json"), []*model.Template{})
}

func TestUpdateTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	parameters := map[string]string{
		"preloads": "TemplateArguments",
	}

	id1 := 101
	template1 := &model.Template{
		ID:              id1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
		Description:     "test1desc",
	}

	id2 := 102
	template2 := &model.Template{
		ID:              id2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		Description:     "test2desc",
		TemplateArguments: []*model.TemplateArgument{
			{
				Name:         "testParameter21",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter21",
				Description:  "testParameter21desc",
			},
		},
	}

	id3 := 103
	template3 := &model.Template{
		ID:              id3,
		Name:            "test3",
		TemplateContent: "TestTemplate3",
		Description:     "test3desc",
		TemplateArguments: []*model.TemplateArgument{
			{
				ID:           1000,
				Name:         "testParameter31",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter31",
				Description:  "testParameter31desc",
			},
			{
				ID:           1001,
				Name:         "testParameter32",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter32",
				Description:  "testParameter31desc",
			},
		},
	}

	id4 := 104
	template4 := &model.Template{
		ID:              id4,
		Name:            "test4",
		TemplateContent: "TestTemplate4",
		Description:     "test4desc",
		TemplateArguments: []*model.TemplateArgument{
			{
				ID:           2000,
				Name:         "testParameter41",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter41",
				Description:  "testParameter41desc",
			},
			{
				ID:           2001,
				Name:         "testParameter42",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter42",
				Description:  "testParameter42desc",
			},
			{
				ID:           2002,
				Name:         "testParameter43",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter43",
				Description:  "testParameter43desc",
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
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_1.json"), &model.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id1), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_2.json"), &model.Template{})

	template2.TemplateArguments = nil

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id2), nil), template2)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_3.json"), &model.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id2), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_4.json"), &model.Template{})

	template3.TemplateArguments[1].Name = "testParameter32Updated"
	template3.TemplateArguments[1].Type = model.TemplateArgumentTypeInt
	template3.TemplateArguments[1].DefaultValue = "99999"
	template3.TemplateArguments[1].Description = "testParameter32descUpdated"

	template3.TemplateArguments = append(
		template3.TemplateArguments,
		&model.TemplateArgument{
			ID:           1002,
			Name:         "testParameter34",
			Type:         model.TemplateArgumentTypeString,
			DefaultValue: "TestParameter34",
		},
	)

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id3), nil), template3)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_5.json"), &model.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id3), parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_6.json"), &model.Template{})

	template3.TemplateArguments[0].ToBeDeleted = true
	template4.TemplateArguments[0].ToBeDeleted = true
	template4.TemplateArguments[2].ToBeDeleted = true
	template4.TemplateArguments = append(template4.TemplateArguments,
		[]*model.TemplateArgument{
			{
				ID:           2003,
				Name:         "testParameter44",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "TestParameter44",
				Description:  "testParameter44desc",
			},
			{
				Type:         model.TemplateArgumentTypeString,
				Name:         "testParameter45",
				DefaultValue: "TestParameter45",
				Description:  "testParameter45desc",
			},
		}...,
	)

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id3), nil), template3)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_7.json"), &model.Template{})

	responseText, code = Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "templates", strconv.Itoa(id4), nil), template4)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_8.json"), &model.Template{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "templates", parameters), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplate_9.json"), []*model.Template{})

}

func TestDeleteTemplate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &model.Template{
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

func TestGetTemplateArguments_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "template_arguments", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, EmptyArrayString, []*model.TemplateArgument{})
}

func TestCreateTemplateArguments(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &model.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	templateArgument1 := &model.TemplateArgument{
		TemplateID:   1,
		Name:         "testParameter1",
		Type:         model.TemplateArgumentTypeInt,
		DefaultValue: "123",
		Description:  "testParameter1desc",
	}
	templateArgument2 := &model.TemplateArgument{
		TemplateID:   1,
		Name:         "testParameter2",
		Type:         model.TemplateArgumentTypeFloat,
		DefaultValue: "456.789",
		Description:  "testParameter2desc",
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument1)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplateArgument_1.json"), &model.TemplateArgument{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument2)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "template/TestCreateTemplateArgument_2.json"), &model.TemplateArgument{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "template_arguments", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestCreateTemplateArgument_3.json"), []*model.TemplateArgument{})
}

func TestUpdateTemplateArguments(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &model.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templateArgument1 := &model.TemplateArgument{
		TemplateID:   id,
		Name:         "testParameter1",
		Type:         model.TemplateArgumentTypeBool,
		DefaultValue: "true",
		Description:  "testParameter1desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument1)

	templateArgument1.Name = "templateArgument1Updated"
	templateArgument1.Type = model.TemplateArgumentTypeInt
	templateArgument1.DefaultValue = "999"
	templateArgument1.Description = "testParameter1descUpdated"

	responseText, code := Execute(t, http.MethodPut, GenerateSingleResourceURL(server, "template_arguments", strconv.Itoa(id), nil), templateArgument1)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplateArgument_1.json"), &model.TemplateArgument{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_arguments", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestUpdateTemplateArgument_2.json"), &model.TemplateArgument{})
}

func TestDeleteTemplateArguments(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template := &model.Template{
		ID:              1,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templateArgument1 := &model.TemplateArgument{
		TemplateID:   id,
		Name:         "testParameter1",
		Type:         model.TemplateArgumentTypeString,
		DefaultValue: "TestParameter1",
		Description:  "testParameter1desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument1)

	templateArgument1.Name = "templateArgument1Updated"
	templateArgument1.DefaultValue = "TestParameter1Updated"
	templateArgument1.Description = "TestParameter1descUpdated"

	responseText, code := Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "template_arguments", strconv.Itoa(id), nil), nil)
	CheckResponseText(t, code, http.StatusNoContent, responseText, []byte{})

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_arguments", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplateArgument_1.json"), &ErrorResponseText{})
}

func TestDeleteTemplateArguments_Cascade(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	templateID := 1
	template := &model.Template{
		ID:              templateID,
		Name:            "test",
		TemplateContent: "TestTemplate",
		Description:     "testdesc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	id := 1
	templateArgument1 := &model.TemplateArgument{
		TemplateID:   id,
		Name:         "testParameter1",
		Type:         model.TemplateArgumentTypeString,
		DefaultValue: "TestParameter1",
		Description:  "testParameter1desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument1)

	Execute(t, http.MethodDelete, GenerateSingleResourceURL(server, "templates", strconv.Itoa(templateID), nil), nil)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_arguments", strconv.Itoa(id), nil), nil)
	CheckResponseJSON(t, code, http.StatusNotFound, responseText, LoadExpectation(t, "template/TestDeleteTemplateArgument_Cascade_1.json"), &ErrorResponseText{})
}

func TestTemplate_ExtractFromDesign(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	template1 := &model.Template{
		ID:              1,
		Name:            "test1",
		TemplateContent: "TestTemplate1",
		Description:     "test1desc",
	}
	template2 := &model.Template{
		ID:              2,
		Name:            "test2",
		TemplateContent: "TestTemplate2",
		Description:     "test2desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)

	templateArgument11 := &model.TemplateArgument{
		TemplateID:   1,
		Name:         "testParameter11",
		Type:         model.TemplateArgumentTypeString,
		DefaultValue: "TestParameter11",
		Description:  "testParameter11desc",
	}
	templateArgument12 := &model.TemplateArgument{
		TemplateID:   1,
		Name:         "testParameter12",
		Type:         model.TemplateArgumentTypeString,
		DefaultValue: "TestParameter12",
		Description:  "testParameter12desc",
	}
	templateArgument21 := &model.TemplateArgument{
		TemplateID:   2,
		Name:         "testParameter21",
		Type:         model.TemplateArgumentTypeString,
		DefaultValue: "TestParameter21",
		Description:  "testParameter21desc",
	}
	templateArgument22 := &model.TemplateArgument{
		TemplateID:   2,
		Name:         "testParameter22",
		Type:         model.TemplateArgumentTypeString,
		DefaultValue: "TestParameter22",
		Description:  "testParameter22desc",
	}
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument11)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument12)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument21)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument22)
}

func TestTemplate_Generate(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 100
	template := &model.Template{
		ID:   id,
		Name: "test",
		TemplateContent: "" +
			`testParameterInt = {{ .Parameter.testParameterInt }}
testParameterInt8 = {{ .Parameter.testParameterInt8 }}
testParameterInt16 = {{ .Parameter.testParameterInt16 }}
testParameterInt32 = {{ .Parameter.testParameterInt32 }}
testParameterInt64 = {{ .Parameter.testParameterInt64 }}
testParameterUint = {{ .Parameter.testParameterUint }}
testParameterUint8 = {{ .Parameter.testParameterUint8 }}
testParameterUint16 = {{ .Parameter.testParameterUint16 }}
testParameterUint32 = {{ .Parameter.testParameterUint32 }}
testParameterUint64 = {{ .Parameter.testParameterUint64 }}
testParameterFloat32 = {{ .Parameter.testParameterFloat32 }}
testParameterFloat64 = {{ .Parameter.testParameterFloat64 }}
testParameterBool = {{ .Parameter.testParameterBool }}
testParameterString = {{ .Parameter.testParameterString }}
testParameterIntOverride = {{ .Parameter.testParameterIntOverride }}
testParameterFloatOverride = {{ .Parameter.testParameterFloatOverride }}
testParameterBoolOverride = {{ .Parameter.testParameterBoolOverride }}
testParameterStringOverride = {{ .Parameter.testParameterStringOverride }}`,
		TemplateArguments: []*model.TemplateArgument{
			{
				Name:         "testParameterInt",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "123",
			},
			{
				Name:         "testParameterInt8",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "124",
			},
			{
				Name:         "testParameterInt16",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "125",
			},
			{
				Name:         "testParameterInt32",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "126",
			},
			{
				Name:         "testParameterInt64",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "127",
			},
			{
				Name:         "testParameterUint",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "1123",
			},
			{
				Name:         "testParameterUint8",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "1124",
			},
			{
				Name:         "testParameterUint16",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "1125",
			},
			{
				Name:         "testParameterUint32",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "1126",
			},
			{
				Name:         "testParameterUint64",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "1127",
			},
			{
				Name:         "testParameterFloat32",
				Type:         model.TemplateArgumentTypeFloat,
				DefaultValue: "123.1",
			},
			{
				Name:         "testParameterFloat64",
				Type:         model.TemplateArgumentTypeFloat,
				DefaultValue: "123.2",
			},
			{
				Name:         "testParameterBool",
				Type:         model.TemplateArgumentTypeBool,
				DefaultValue: "true",
			},
			{
				Name:         "testParameterString",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "ABCDE",
			},
			{
				Name:         "testParameterIntOverride",
				Type:         model.TemplateArgumentTypeInt,
				DefaultValue: "0",
			},
			{
				Name:         "testParameterFloatOverride",
				Type:         model.TemplateArgumentTypeFloat,
				DefaultValue: "0",
			},
			{
				Name:         "testParameterBoolOverride",
				Type:         model.TemplateArgumentTypeBool,
				DefaultValue: "false",
			},
			{
				Name:         "testParameterStringOverride",
				Type:         model.TemplateArgumentTypeString,
				DefaultValue: "nothing",
			},
		},
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template)

	parameters := map[string]string{
		"p[testParameterIntOverride]":    "100",
		"p[testParameterFloatOverride]":  "200.123",
		"p[testParameterBoolOverride]":   "true",
		"p[testParameterStringOverride]": "QWERTY",
	}
	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "generation", parameters), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_1.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_generations_by_name", template.Name, parameters), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_1.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "raw", parameters), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_2.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_raws_by_name", template.Name, parameters), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_2.txt"))

}

func TestTemplate_Functions(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 1
	template1 := &model.Template{
		ID:   id,
		Name: "test1",
		TemplateContent: `--- calc ---
i = 100
{{- $i := 100 }}

i = i + 2
{{- $i := .Core.Add $i 2 }}
{{ $i }}

i = i - 4
{{- $i := .Core.Sub $i 4 }}
{{ $i }}

i = i * 6
{{- $i := .Core.Mul $i 6 }}
{{ $i }}

i = i / 2
{{- $i := .Core.Div $i 2 }}
{{ $i }}

i = i mod 5
{{- $i := .Core.Mod $i 5 }}
{{ $i }}

--- conversion ---
{{- $vint := .Core.Int "100" }}
{{- $vint8 := .Core.Int8 "101" }}
{{- $vint16 := .Core.Int16 "102" }}
{{- $vint32 := .Core.Int32 "103" }}
{{- $vint64 := .Core.Int64 "104" }}
{{- $vuint := .Core.Uint "200" }}
{{- $vuint8 := .Core.Uint8 "201" }}
{{- $vuint16 := .Core.Uint16 "202" }}
{{- $vuint32 := .Core.Uint32 "203" }}
{{- $vuint64 := .Core.Uint64 "204" }}
{{- $vfloat32 := .Core.Float32 "300.1" }}
{{- $vfloat64 := .Core.Float64 "300.2" }}
{{- $vboolean := .Core.Boolean "false" }}
{{- $vmap := .Core.Map "key" "value" }}
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
{{- $sint := .Core.String $vint }}
{{- $sint8 := .Core.String $vint8 }}
{{- $sint16 := .Core.String $vint16 }}
{{- $sint32 := .Core.String $vint32 }}
{{- $sint64 := .Core.String $vint64 }}
{{- $suint := .Core.String $vuint }}
{{- $suint8 := .Core.String $vuint8 }}
{{- $suint16 := .Core.String $vuint16 }}
{{- $suint32 := .Core.String $vuint32 }}
{{- $suint64 := .Core.String $vuint64 }}
{{- $sfloat32:= .Core.String $vfloat32 }}
{{- $sfloat64 := .Core.String $vfloat64 }}
{{- $sboolean := .Core.String $vboolean }}
{{- $sobject := .Core.String $vmap }}
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
{{ .Core.Int $vint }}
{{ .Core.Int $vint8 }}
{{ .Core.Int $vint16 }}
{{ .Core.Int $vint32 }}
{{ .Core.Int $vint64 }}
{{ .Core.Int $vuint }}
{{ .Core.Int $vuint8 }}
{{ .Core.Int $vuint16 }}
{{ .Core.Int $vuint32 }}
{{ .Core.Int $vuint64 }}
{{ .Core.Int $vfloat32 }}
{{ .Core.Int $vfloat64 }}
{{ .Core.Int8 $vint }}
{{ .Core.Int8 $vint8 }}
{{ .Core.Int8 $vint16 }}
{{ .Core.Int8 $vint32 }}
{{ .Core.Int8 $vint64 }}
{{ .Core.Int8 $vuint }}
{{ .Core.Int8 $vuint8 }}
{{ .Core.Int8 $vuint16 }}
{{ .Core.Int8 $vuint32 }}
{{ .Core.Int8 $vuint64 }}
{{ .Core.Int8 $vfloat32 }}
{{ .Core.Int8 $vfloat64 }}
{{ .Core.Int16 $vint }}
{{ .Core.Int16 $vint8 }}
{{ .Core.Int16 $vint16 }}
{{ .Core.Int16 $vint32 }}
{{ .Core.Int16 $vint64 }}
{{ .Core.Int16 $vuint }}
{{ .Core.Int16 $vuint8 }}
{{ .Core.Int16 $vuint16 }}
{{ .Core.Int16 $vuint32 }}
{{ .Core.Int16 $vuint64 }}
{{ .Core.Int16 $vfloat32 }}
{{ .Core.Int16 $vfloat64 }}
{{ .Core.Int32 $vint }}
{{ .Core.Int32 $vint8 }}
{{ .Core.Int32 $vint16 }}
{{ .Core.Int32 $vint32 }}
{{ .Core.Int32 $vint64 }}
{{ .Core.Int32 $vuint }}
{{ .Core.Int32 $vuint8 }}
{{ .Core.Int32 $vuint16 }}
{{ .Core.Int32 $vuint32 }}
{{ .Core.Int32 $vuint64 }}
{{ .Core.Int32 $vfloat32 }}
{{ .Core.Int32 $vfloat64 }}
{{ .Core.Int64 $vint }}
{{ .Core.Int64 $vint8 }}
{{ .Core.Int64 $vint16 }}
{{ .Core.Int64 $vint32 }}
{{ .Core.Int64 $vint64 }}
{{ .Core.Int64 $vuint }}
{{ .Core.Int64 $vuint8 }}
{{ .Core.Int64 $vuint16 }}
{{ .Core.Int64 $vuint32 }}
{{ .Core.Int64 $vuint64 }}
{{ .Core.Int64 $vfloat32 }}
{{ .Core.Int64 $vfloat64 }}
{{ .Core.Uint $vint }}
{{ .Core.Uint $vint8 }}
{{ .Core.Uint $vint16 }}
{{ .Core.Uint $vint32 }}
{{ .Core.Uint $vint64 }}
{{ .Core.Uint $vuint }}
{{ .Core.Uint $vuint8 }}
{{ .Core.Uint $vuint16 }}
{{ .Core.Uint $vuint32 }}
{{ .Core.Uint $vuint64 }}
{{ .Core.Uint $vfloat32 }}
{{ .Core.Uint $vfloat64 }}
{{ .Core.Uint8 $vint }}
{{ .Core.Uint8 $vint8 }}
{{ .Core.Uint8 $vint16 }}
{{ .Core.Uint8 $vint32 }}
{{ .Core.Uint8 $vint64 }}
{{ .Core.Uint8 $vuint }}
{{ .Core.Uint8 $vuint8 }}
{{ .Core.Uint8 $vuint16 }}
{{ .Core.Uint8 $vuint32 }}
{{ .Core.Uint8 $vuint64 }}
{{ .Core.Uint8 $vfloat32 }}
{{ .Core.Uint8 $vfloat64 }}
{{ .Core.Uint16 $vint }}
{{ .Core.Uint16 $vint8 }}
{{ .Core.Uint16 $vint16 }}
{{ .Core.Uint16 $vint32 }}
{{ .Core.Uint16 $vint64 }}
{{ .Core.Uint16 $vuint }}
{{ .Core.Uint16 $vuint8 }}
{{ .Core.Uint16 $vuint16 }}
{{ .Core.Uint16 $vuint32 }}
{{ .Core.Uint16 $vuint64 }}
{{ .Core.Uint16 $vfloat32 }}
{{ .Core.Uint16 $vfloat64 }}
{{ .Core.Uint32 $vint }}
{{ .Core.Uint32 $vint8 }}
{{ .Core.Uint32 $vint16 }}
{{ .Core.Uint32 $vint32 }}
{{ .Core.Uint32 $vint64 }}
{{ .Core.Uint32 $vuint }}
{{ .Core.Uint32 $vuint8 }}
{{ .Core.Uint32 $vuint16 }}
{{ .Core.Uint32 $vuint32 }}
{{ .Core.Uint32 $vuint64 }}
{{ .Core.Uint32 $vfloat32 }}
{{ .Core.Uint32 $vfloat64 }}
{{ .Core.Uint64 $vint }}
{{ .Core.Uint64 $vint8 }}
{{ .Core.Uint64 $vint16 }}
{{ .Core.Uint64 $vint32 }}
{{ .Core.Uint64 $vint64 }}
{{ .Core.Uint64 $vuint }}
{{ .Core.Uint64 $vuint8 }}
{{ .Core.Uint64 $vuint16 }}
{{ .Core.Uint64 $vuint32 }}
{{ .Core.Uint64 $vuint64 }}
{{ .Core.Uint64 $vfloat32 }}
{{ .Core.Uint64 $vfloat64 }}
{{ .Core.Float32 $vint }}
{{ .Core.Float32 $vint8 }}
{{ .Core.Float32 $vint16 }}
{{ .Core.Float32 $vint32 }}
{{ .Core.Float32 $vint64 }}
{{ .Core.Float32 $vuint }}
{{ .Core.Float32 $vuint8 }}
{{ .Core.Float32 $vuint16 }}
{{ .Core.Float32 $vuint32 }}
{{ .Core.Float32 $vuint64 }}
{{ .Core.Float32 $vfloat32 }}
{{ .Core.Float32 $vfloat64 }}
{{ .Core.Float64 $vint }}
{{ .Core.Float64 $vint8 }}
{{ .Core.Float64 $vint16 }}
{{ .Core.Float64 $vint32 }}
{{ .Core.Float64 $vint64 }}
{{ .Core.Float64 $vuint }}
{{ .Core.Float64 $vuint8 }}
{{ .Core.Float64 $vuint16 }}
{{ .Core.Float64 $vuint32 }}
{{ .Core.Float64 $vuint64 }}
{{ .Core.Float64 $vfloat32 }}
{{ .Core.Float64 $vfloat64 }}

--- string ---
join
{{- $data := .Core.Slice 1 99.99 "a" false }}
{{- $sj := .Core.Join $data "," }}
join string: {{ $sj}}

split
{{- $ss := .Core.Split $sj "," }}
split slice: {{ $ss }}

--- slice ---
slice init
{{- $slice1 := .Core.Slice }}
slice1: {{ $slice1 }}

slice init
{{- $slice2 := .Core.Slice 1 2 3 4 5 }}
slice2: {{ $slice2 }}

slice init
{{- $slice3 := .Core.Slice 3 4 5 6 7 }}
slice3: {{ $slice3 }}

subslice
{{- $subslice := .Core.SubSlice $slice3 -1 3 }}
{{ $subslice }}

subslice
{{ .Core.SubSlice $slice3 1 -1 }}

subslice
{{ .Core.SubSlice $slice3 -1 -1 }}

subslice
{{ .Core.SubSlice $slice3 1 3 }}

append
{{- $slice3 := .Core.Append $slice3 8 9 10 }}
{{ $slice3 }}

concatenate
{{ .Core.Concatenate $slice2 $slice3 }}

fieldslice
{{- $tpp := .ModelStore.Multi "template_arguments" "" }}
{{- $idfields := .Core.FieldSlice $tpp "ID" }}
{{- $namefields := .Core.FieldSlice $tpp "Name" }}
{{- $idfields := .Core.Sort $idfields "asc" }}
{{- $namefields := .Core.Sort $namefields "asc" }}
{{ $idfields }}
{{ $namefields }}

sort
{{- $sliceint := .Core.Slice }}
{{- $v1 := .Core.Int 3 }}
{{- $v2 := .Core.Int 1 }}
{{- $v3 := .Core.Int -5 }}
{{- $v4 := .Core.Int 2 }}
{{- $v5 := .Core.Int 4 }}
{{- $sliceint := .Core.Append $sliceint $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint := .Core.Sort $sliceint "asc" }}
sliceint asc: {{ $sliceint }}
{{- $sliceint := .Core.Sort $sliceint "desc" }}
sliceint desc: {{ $sliceint }}
{{- $sliceint8 := .Core.Slice }}
{{- $v1 := .Core.Int8 3 }}
{{- $v2 := .Core.Int8 1 }}
{{- $v3 := .Core.Int8 -5 }}
{{- $v4 := .Core.Int8 2 }}
{{- $v5 := .Core.Int8 4 }}
{{- $sliceint8 := .Core.Append $sliceint8 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint8 := .Core.Sort $sliceint8 "asc" }}
sliceint8 asc: {{ $sliceint8}}
{{- $sliceint8 := .Core.Sort $sliceint8 "desc" }}
sliceint8 desc: {{ $sliceint8}}
{{- $sliceint16 := .Core.Slice }}
{{- $v1 := .Core.Int16 3 }}
{{- $v2 := .Core.Int16 1 }}
{{- $v3 := .Core.Int16 -5 }}
{{- $v4 := .Core.Int16 2 }}
{{- $v5 := .Core.Int16 4 }}
{{- $sliceint16 := .Core.Append $sliceint16 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint16 := .Core.Sort $sliceint16 "asc" }}
sliceint16 asc: {{ $sliceint16 }}
{{- $sliceint16 := .Core.Sort $sliceint16 "desc" }}
sliceint16 desc: {{ $sliceint16 }}
{{- $sliceint32 := .Core.Slice }}
{{- $v1 := .Core.Int32 3 }}
{{- $v2 := .Core.Int32 1 }}
{{- $v3 := .Core.Int32 -5 }}
{{- $v4 := .Core.Int32 2 }}
{{- $v5 := .Core.Int32 4 }}
{{- $sliceint32 := .Core.Append $sliceint32 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint32 := .Core.Sort $sliceint32 "asc" }}
sliceint32 asc: {{ $sliceint32 }}
{{- $sliceint32 := .Core.Sort $sliceint32 "desc" }}
sliceint32 desc: {{ $sliceint32 }}
{{- $sliceint64 := .Core.Slice }}
{{- $v1 := .Core.Int64 3 }}
{{- $v2 := .Core.Int64 1 }}
{{- $v3 := .Core.Int64 -5 }}
{{- $v4 := .Core.Int64 2 }}
{{- $v5 := .Core.Int64 4 }}
{{- $sliceint64 := .Core.Append $sliceint64 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint64 := .Core.Sort $sliceint64 "asc" }}
sliceint64 asc: {{ $sliceint64 }}
{{- $sliceint64 := .Core.Sort $sliceint64 "desc" }}
sliceint64 desc: {{ $sliceint64 }}
{{- $sliceuint := .Core.Slice }}
{{- $v1 := .Core.Uint 3 }}
{{- $v2 := .Core.Uint 1 }}
{{- $v3 := .Core.Uint 5 }}
{{- $v4 := .Core.Uint 2 }}
{{- $v5 := .Core.Uint 4 }}
{{- $sliceuint := .Core.Append $sliceuint $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint := .Core.Sort $sliceuint "asc" }}
sliceuint asc: {{ $sliceuint }}
{{- $sliceuint := .Core.Sort $sliceuint "desc" }}
sliceuint desc: {{ $sliceuint }}
{{- $sliceuint8 := .Core.Slice }}
{{- $v1 := .Core.Uint8 3 }}
{{- $v2 := .Core.Uint8 1 }}
{{- $v3 := .Core.Uint8 5 }}
{{- $v4 := .Core.Uint8 2 }}
{{- $v5 := .Core.Uint8 4 }}
{{- $sliceuint8 := .Core.Append $sliceuint8 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint8 := .Core.Sort $sliceuint8 "asc" }}
sliceuint8 asc: {{ $sliceuint8}}
{{- $sliceuint8 := .Core.Sort $sliceuint8 "desc" }}
sliceuint8 desc: {{ $sliceuint8}}
{{- $sliceuint16 := .Core.Slice }}
{{- $v1 := .Core.Uint16 3 }}
{{- $v2 := .Core.Uint16 1 }}
{{- $v3 := .Core.Uint16 5 }}
{{- $v4 := .Core.Uint16 2 }}
{{- $v5 := .Core.Uint16 4 }}
{{- $sliceuint16 := .Core.Append $sliceuint16 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint16 := .Core.Sort $sliceuint16 "asc" }}
sliceuint16 asc: {{ $sliceuint16 }}
{{- $sliceuint16 := .Core.Sort $sliceuint16 "desc" }}
sliceuint16 desc: {{ $sliceuint16 }}
{{- $sliceuint32 := .Core.Slice }}
{{- $v1 := .Core.Uint32 3 }}
{{- $v2 := .Core.Uint32 1 }}
{{- $v3 := .Core.Uint32 5 }}
{{- $v4 := .Core.Uint32 2 }}
{{- $v5 := .Core.Uint32 4 }}
{{- $sliceuint32 := .Core.Append $sliceuint32 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint32 := .Core.Sort $sliceuint32 "asc" }}
sliceuint32 asc: {{ $sliceuint32 }}
{{- $sliceuint32 := .Core.Sort $sliceuint32 "desc" }}
sliceuint32 desc: {{ $sliceuint32 }}
{{- $sliceuint64 := .Core.Slice }}
{{- $v1 := .Core.Uint64 3 }}
{{- $v2 := .Core.Uint64 1 }}
{{- $v3 := .Core.Uint64 5 }}
{{- $v4 := .Core.Uint64 2 }}
{{- $v5 := .Core.Uint64 4 }}
{{- $sliceuint64 := .Core.Append $sliceuint64 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint64 := .Core.Sort $sliceuint64 "asc" }}
sliceuint64 asc: {{ $sliceuint64 }}
{{- $sliceuint64 := .Core.Sort $sliceuint64 "desc" }}
sliceuint64 desc: {{ $sliceuint64 }}
{{- $slicefloat64 := .Core.Slice }}
{{- $v1 := .Core.Float32 3.3 }}
{{- $v2 := .Core.Float32 1 }}
{{- $v3 := .Core.Float32 -5.1 }}
{{- $v4 := .Core.Float32 2.2 }}
{{- $v5 := .Core.Float32 4 }}
{{- $slicefloat64 := .Core.Append $slicefloat64 $v1 $v2 $v3 $v4 $v5 }}
{{- $slicefloat64 := .Core.Sort $slicefloat64 "asc" }}
slicefloat32 asc: {{ $slicefloat64 }}
{{- $slicefloat64 := .Core.Sort $slicefloat64 "desc" }}
slicefloat32 desc: {{ $slicefloat64 }}
{{- $slicefloat64 := .Core.Slice }}
{{- $v1 := .Core.Float64 3.3 }}
{{- $v2 := .Core.Float64 1 }}
{{- $v3 := .Core.Float64 -5.1 }}
{{- $v4 := .Core.Float64 2.2 }}
{{- $v5 := .Core.Float64 4 }}
{{- $slicefloat64 := .Core.Append $slicefloat64 $v1 $v2 $v3 $v4 $v5 }}
{{- $slicefloat64 := .Core.Sort $slicefloat64 "asc" }}
slicefloat64 asc: {{ $slicefloat64 }}
{{- $slicefloat64 := .Core.Sort $slicefloat64 "desc" }}
slicefloat64 desc: {{ $slicefloat64 }}
{{- $slicestring := .Core.Slice }}
{{- $v1 := "3.3" }}
{{- $v2 := "ABC" }}
{{- $v3 := "-5.1" }}
{{- $v4 := "012" }}
{{- $v5 := "def" }}
{{- $slicestring := .Core.Append $slicestring $v1 $v2 $v3 $v4 $v5 }}
{{- $slicestring := .Core.Sort $slicestring "asc" }}
slicestring asc: {{ $slicestring}}
{{- $slicestring := .Core.Sort $slicestring "desc" }}
slicestring desc: {{ $slicestring}}

--- map ---
map
{{- $map1 := .Core.Map }}
map1: {{ $map1 }}

map init get
{{- $map2 := .Core.Map 1 "A" 2 "B" }}
map2[1]: {{ .Core.Get $map2 1 }}
map2[2]: {{ .Core.Get $map2 2 }}

map init get
{{- $map3 := .Core.Map 1 "C" 3 "D" }}
map3[1]: {{ .Core.Get $map3 1 }}
map3[3]: {{ .Core.Get $map3 3 }}

map exists
{{ .Core.Exists $map2 0 }}
{{- $e := .Core.Exists $map2 1 }}
{{ if eq $e true }}TRUE!!{{ else }}FALSE!!{{ end }}

map put
{{- $null := .Core.Put $map1 4 "E" }}
{{- $null := .Core.Put $map1 5 "F" }}
{{- $null := .Core.Put $map1 6 "G" }}
map1[4]: {{ .Core.Get $map1 4 }}
map1[5]: {{ .Core.Get $map1 5 }}
map1[6]: {{ .Core.Get $map1 6 }}

map delete
{{- $null := .Core.Delete $map1 5 }}
map1[3]: {{ .Core.Get $map1 3 }}
map1[4]: {{ .Core.Get $map1 4 }}
map1[5]: {{ .Core.Get $map1 5 }}
map1[6]: {{ .Core.Get $map1 6 }}
map1[7]: {{ .Core.Get $map1 7 }}

map merge
{{- $null :=  .Core.Merge $map2 $map1 }}
{{- $null :=  .Core.Merge $map3 $map1 }}
map1[0]: {{ .Core.Get $map1 0 }}
map1[1]: {{ .Core.Get $map1 1 }}
map1[2]: {{ .Core.Get $map1 2 }}
map1[3]: {{ .Core.Get $map1 3 }}
map1[4]: {{ .Core.Get $map1 4 }}
map1[5]: {{ .Core.Get $map1 5 }}
map1[6]: {{ .Core.Get $map1 6 }}
map1[7]: {{ .Core.Get $map1 7 }}

map keys
{{- $keys := .Core.Keys $map1 }}
{{- $keys := .Core.Sort $keys "asc" }}
keys of map1: {{ $keys }}

--- model store ---
multi
{{- $m := .ModelStore.Multi "templates" "preloads=TemplateArguments" }}
{{- $t := index $m 0 }}
{{ $t.Name }}
{{- $p1 := index $t.TemplateArguments 0 }}
{{ $p1.Name }}={{ $p1.DefaultValue }}
{{- $p2 := index $t.TemplateArguments 1 }}
{{ $p2.Name }}={{ $p2.DefaultValue }}

single
{{- $path := printf "/templates/%d" .Parameter.testParameter11 }}
{{- $s := .ModelStore.Single $path "preloads=TemplateArguments" }}
{{ $s.Name }}
{{- $p1 := index $s.TemplateArguments 0 }}
{{ $p1.Name }}={{ $p1.DefaultValue }}
{{- $p2 := index $s.TemplateArguments 1 }}
{{ $p2.Name }}={{ $p2.DefaultValue }}

first
{{- $f := .ModelStore.First "templates" "q[name]=test1&preloads=TemplateArguments" }}
{{ $f.Name }}
{{- $p1 := index $t.TemplateArguments 0 }}
{{ $p1.Name }}={{ $p1.DefaultValue }}
{{- $p2 := index $t.TemplateArguments 1 }}
{{ $p2.Name }}={{ $p2.DefaultValue }}

total
{{- $t := .ModelStore.Total "/template_arguments" }}
{{ $t }}

--- hash ---
hash
{{- $h := .Core.Hash $s.TemplateArguments "Name" }}
hash[testParameter11]={{ .Core.Get $h "testParameter11" }}
hash[testParameter12]={{ .Core.Get $h "testParameter12" }}

--- slicemap ---
slicemap
{{- $p := .ModelStore.Multi "template_arguments" "" }}
{{- $z := .Core.SliceMap $p "Name" }}
{{- $z1 := .Core.Get $z "testParameter11" }}
{{- $z2 := .Core.Get $z "testParameter12" }}
{{- $z3 := .Core.Get $z "testParameter1X" }}
{{- range $i, $v := $z1 }}
slicemap[testParameter1][{{ $i }}]={{ $v }}
{{- end }}
{{- range $i, $v := $z2 }}
slicemap[testParameter2][{{ $i }}]={{ $v }}
{{- end }}
{{- range $i, $v := $z3 }}
slicemap[testParameter1X][{{ $i }}]={{ $v }}
{{- end }}

--- sequence ---
sequence
{{- $s := .Core.Sequence 1 10 }}
{{- range $i, $v := $s }}
sequence[{{ $i }}]={{ $v }}
{{- end }}

--- net ---
{{- $a := .Network.ParseCIDR "192.168.0.100/24" }}
{{ $a.String }}
{{ $a.NetMask }}
{{ $a.CIDR }}
{{ $a.IncreaseHostAddress }}
{{ $a.DecreaseHostAddress }}
{{ $a.IncreaseNetworkAddress }}
{{ $a.DecreaseNetworkAddress }}
{{ $a.IncreaseIPAddress }}
{{ $a.DecreaseIPAddress }}
{{ $a.LimitedBroadcastAddress }}
{{ $a.NetworkAddress }}
{{ $a.MaxHostAddress }}
{{ $a.MinimumHostAddress }}
{{ $a.IsBroadcastAddress }}
{{ $a.IsNetworkAddress }}
{{- $b := .Network.ParseCIDR "192.168.0.200/24" }}
{{ $a.IsIncluding $b }}
{{- $c := .Network.ParseCIDR "192.168.1.200/24" }}
{{ $a.IsIncluding $c }}
`,
		Description: "test1desc",
	}

	templateArgument11 := &model.TemplateArgument{
		TemplateID:   id,
		Name:         "testParameter11",
		Type:         model.TemplateArgumentTypeInt,
		DefaultValue: "1",
		Description:  "testParameter11desc",
	}
	templateArgument12 := &model.TemplateArgument{
		TemplateID:   id,
		Name:         "testParameter12",
		Type:         model.TemplateArgumentTypeString,
		DefaultValue: "TestParameter12",
		Description:  "testParameter12desc",
	}

	templateArgument13 := &model.TemplateArgument{
		TemplateID:   id,
		Name:         "testParameter1X",
		Type:         model.TemplateArgumentTypeInt,
		DefaultValue: "100",
		Description:  "testParameter1Xdesc",
	}

	id2 := 2
	template2 := &model.Template{
		ID:              id2,
		Name:            "test12",
		TemplateContent: `{{ .Parameter.testParameter1X }}`,
		Description:     "test12desc",
	}
	templateArgument21 := &model.TemplateArgument{
		TemplateID:   id2,
		Name:         "testParameter1X",
		Type:         model.TemplateArgumentTypeInt,
		DefaultValue: "200",
		Description:  "testParameter1Xdesc",
	}

	id3 := 3
	template3 := &model.Template{
		ID:              id3,
		Name:            "test13",
		TemplateContent: `{{ .Parameter.testParameter1X }}`,
		Description:     "test13desc",
	}
	templateArgument31 := &model.TemplateArgument{
		TemplateID:   id3,
		Name:         "testParameter1X",
		Type:         model.TemplateArgumentTypeInt,
		DefaultValue: "300",
		Description:  "testParameter1Xdesc",
	}

	id4 := 4
	template4 := &model.Template{
		ID:   id4,
		Name: "test14",
		TemplateContent: `include test
{{ .Template.Include "test12" "" }}
{{ .Template.Include "test13" "p[testParameter1X]=999" }}
`,
		Description: "test15desc",
	}

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template1)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument11)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument12)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument13)

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template2)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument21)

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template3)
	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "template_arguments", nil), templateArgument31)

	Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "templates", nil), template4)

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_1.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_generations_by_name", template1.Name, nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_1.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id2), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_2.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_generations_by_name", template2.Name, nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_2.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id3), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_3.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_generations_by_name", template3.Name, nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_3.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id4), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_4.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, "template_generations_by_name", template4.Name, nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_4.txt"))
}
