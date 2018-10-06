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
	parametersByName := map[string]string{
		"p[testParameterIntOverride]":    "100",
		"p[testParameterFloatOverride]":  "200.123",
		"p[testParameterBoolOverride]":   "true",
		"p[testParameterStringOverride]": "QWERTY",
		"key_parameter":                  "name",
	}

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "generation", parameters), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_1.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%s", template.Name), "generation", parametersByName), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_1.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "raw", parameters), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_2.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%s", template.Name), "raw", parametersByName), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestGenerateTemplate_2.txt"))

}

func TestTemplate_Functions(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	id := 1
	template1 := &model.Template{
		ID:   id,
		Name: "test1",
		TemplateContent: `--- query ---
len="{{ len .Query }}"
query1="{{ .Query.Get "query1" }}"
query2="{{ .Query.Get "query2" }}"
key_parameter="{{ .Query.Get "key_parameter" }}"

--- calc ---
i = 100
{{- $i := 100 }}

i = i + 2
{{- $i := add $i 2 }}
{{ $i }}

i = i - 4
{{- $i := sub $i 4 }}
{{ $i }}

i = i * 6
{{- $i := mul $i 6 }}
{{ $i }}

i = i / 2
{{- $i := div $i 2 }}
{{ $i }}

i = i mod 5
{{- $i := mod $i 5 }}
{{ $i }}

--- conversion ---
{{- $vint := .Conversion.Int "100" }}
{{- $vint8 := .Conversion.Int8 "101" }}
{{- $vint16 := .Conversion.Int16 "102" }}
{{- $vint32 := .Conversion.Int32 "103" }}
{{- $vint64 := .Conversion.Int64 "104" }}
{{- $vuint := .Conversion.Uint "200" }}
{{- $vuint8 := .Conversion.Uint8 "201" }}
{{- $vuint16 := .Conversion.Uint16 "202" }}
{{- $vuint32 := .Conversion.Uint32 "203" }}
{{- $vuint64 := .Conversion.Uint64 "204" }}
{{- $vfloat32 := .Conversion.Float32 "300.1" }}
{{- $vfloat64 := .Conversion.Float64 "300.2" }}
{{- $vboolean := .Conversion.Boolean "false" }}
{{- $vmap := .Collection.Map "key" "value" }}
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
{{- $sint := .Conversion.String $vint }}
{{- $sint8 := .Conversion.String $vint8 }}
{{- $sint16 := .Conversion.String $vint16 }}
{{- $sint32 := .Conversion.String $vint32 }}
{{- $sint64 := .Conversion.String $vint64 }}
{{- $suint := .Conversion.String $vuint }}
{{- $suint8 := .Conversion.String $vuint8 }}
{{- $suint16 := .Conversion.String $vuint16 }}
{{- $suint32 := .Conversion.String $vuint32 }}
{{- $suint64 := .Conversion.String $vuint64 }}
{{- $sfloat32:= .Conversion.String $vfloat32 }}
{{- $sfloat64 := .Conversion.String $vfloat64 }}
{{- $sboolean := .Conversion.String $vboolean }}
{{- $sobject := .Conversion.String $vmap }}
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
{{ .Conversion.Int $vint }}
{{ .Conversion.Int $vint8 }}
{{ .Conversion.Int $vint16 }}
{{ .Conversion.Int $vint32 }}
{{ .Conversion.Int $vint64 }}
{{ .Conversion.Int $vuint }}
{{ .Conversion.Int $vuint8 }}
{{ .Conversion.Int $vuint16 }}
{{ .Conversion.Int $vuint32 }}
{{ .Conversion.Int $vuint64 }}
{{ .Conversion.Int $vfloat32 }}
{{ .Conversion.Int $vfloat64 }}
{{ .Conversion.Int8 $vint }}
{{ .Conversion.Int8 $vint8 }}
{{ .Conversion.Int8 $vint16 }}
{{ .Conversion.Int8 $vint32 }}
{{ .Conversion.Int8 $vint64 }}
{{ .Conversion.Int8 $vuint }}
{{ .Conversion.Int8 $vuint8 }}
{{ .Conversion.Int8 $vuint16 }}
{{ .Conversion.Int8 $vuint32 }}
{{ .Conversion.Int8 $vuint64 }}
{{ .Conversion.Int8 $vfloat32 }}
{{ .Conversion.Int8 $vfloat64 }}
{{ .Conversion.Int16 $vint }}
{{ .Conversion.Int16 $vint8 }}
{{ .Conversion.Int16 $vint16 }}
{{ .Conversion.Int16 $vint32 }}
{{ .Conversion.Int16 $vint64 }}
{{ .Conversion.Int16 $vuint }}
{{ .Conversion.Int16 $vuint8 }}
{{ .Conversion.Int16 $vuint16 }}
{{ .Conversion.Int16 $vuint32 }}
{{ .Conversion.Int16 $vuint64 }}
{{ .Conversion.Int16 $vfloat32 }}
{{ .Conversion.Int16 $vfloat64 }}
{{ .Conversion.Int32 $vint }}
{{ .Conversion.Int32 $vint8 }}
{{ .Conversion.Int32 $vint16 }}
{{ .Conversion.Int32 $vint32 }}
{{ .Conversion.Int32 $vint64 }}
{{ .Conversion.Int32 $vuint }}
{{ .Conversion.Int32 $vuint8 }}
{{ .Conversion.Int32 $vuint16 }}
{{ .Conversion.Int32 $vuint32 }}
{{ .Conversion.Int32 $vuint64 }}
{{ .Conversion.Int32 $vfloat32 }}
{{ .Conversion.Int32 $vfloat64 }}
{{ .Conversion.Int64 $vint }}
{{ .Conversion.Int64 $vint8 }}
{{ .Conversion.Int64 $vint16 }}
{{ .Conversion.Int64 $vint32 }}
{{ .Conversion.Int64 $vint64 }}
{{ .Conversion.Int64 $vuint }}
{{ .Conversion.Int64 $vuint8 }}
{{ .Conversion.Int64 $vuint16 }}
{{ .Conversion.Int64 $vuint32 }}
{{ .Conversion.Int64 $vuint64 }}
{{ .Conversion.Int64 $vfloat32 }}
{{ .Conversion.Int64 $vfloat64 }}
{{ .Conversion.Uint $vint }}
{{ .Conversion.Uint $vint8 }}
{{ .Conversion.Uint $vint16 }}
{{ .Conversion.Uint $vint32 }}
{{ .Conversion.Uint $vint64 }}
{{ .Conversion.Uint $vuint }}
{{ .Conversion.Uint $vuint8 }}
{{ .Conversion.Uint $vuint16 }}
{{ .Conversion.Uint $vuint32 }}
{{ .Conversion.Uint $vuint64 }}
{{ .Conversion.Uint $vfloat32 }}
{{ .Conversion.Uint $vfloat64 }}
{{ .Conversion.Uint8 $vint }}
{{ .Conversion.Uint8 $vint8 }}
{{ .Conversion.Uint8 $vint16 }}
{{ .Conversion.Uint8 $vint32 }}
{{ .Conversion.Uint8 $vint64 }}
{{ .Conversion.Uint8 $vuint }}
{{ .Conversion.Uint8 $vuint8 }}
{{ .Conversion.Uint8 $vuint16 }}
{{ .Conversion.Uint8 $vuint32 }}
{{ .Conversion.Uint8 $vuint64 }}
{{ .Conversion.Uint8 $vfloat32 }}
{{ .Conversion.Uint8 $vfloat64 }}
{{ .Conversion.Uint16 $vint }}
{{ .Conversion.Uint16 $vint8 }}
{{ .Conversion.Uint16 $vint16 }}
{{ .Conversion.Uint16 $vint32 }}
{{ .Conversion.Uint16 $vint64 }}
{{ .Conversion.Uint16 $vuint }}
{{ .Conversion.Uint16 $vuint8 }}
{{ .Conversion.Uint16 $vuint16 }}
{{ .Conversion.Uint16 $vuint32 }}
{{ .Conversion.Uint16 $vuint64 }}
{{ .Conversion.Uint16 $vfloat32 }}
{{ .Conversion.Uint16 $vfloat64 }}
{{ .Conversion.Uint32 $vint }}
{{ .Conversion.Uint32 $vint8 }}
{{ .Conversion.Uint32 $vint16 }}
{{ .Conversion.Uint32 $vint32 }}
{{ .Conversion.Uint32 $vint64 }}
{{ .Conversion.Uint32 $vuint }}
{{ .Conversion.Uint32 $vuint8 }}
{{ .Conversion.Uint32 $vuint16 }}
{{ .Conversion.Uint32 $vuint32 }}
{{ .Conversion.Uint32 $vuint64 }}
{{ .Conversion.Uint32 $vfloat32 }}
{{ .Conversion.Uint32 $vfloat64 }}
{{ .Conversion.Uint64 $vint }}
{{ .Conversion.Uint64 $vint8 }}
{{ .Conversion.Uint64 $vint16 }}
{{ .Conversion.Uint64 $vint32 }}
{{ .Conversion.Uint64 $vint64 }}
{{ .Conversion.Uint64 $vuint }}
{{ .Conversion.Uint64 $vuint8 }}
{{ .Conversion.Uint64 $vuint16 }}
{{ .Conversion.Uint64 $vuint32 }}
{{ .Conversion.Uint64 $vuint64 }}
{{ .Conversion.Uint64 $vfloat32 }}
{{ .Conversion.Uint64 $vfloat64 }}
{{ .Conversion.Float32 $vint }}
{{ .Conversion.Float32 $vint8 }}
{{ .Conversion.Float32 $vint16 }}
{{ .Conversion.Float32 $vint32 }}
{{ .Conversion.Float32 $vint64 }}
{{ .Conversion.Float32 $vuint }}
{{ .Conversion.Float32 $vuint8 }}
{{ .Conversion.Float32 $vuint16 }}
{{ .Conversion.Float32 $vuint32 }}
{{ .Conversion.Float32 $vuint64 }}
{{ .Conversion.Float32 $vfloat32 }}
{{ .Conversion.Float32 $vfloat64 }}
{{ .Conversion.Float64 $vint }}
{{ .Conversion.Float64 $vint8 }}
{{ .Conversion.Float64 $vint16 }}
{{ .Conversion.Float64 $vint32 }}
{{ .Conversion.Float64 $vint64 }}
{{ .Conversion.Float64 $vuint }}
{{ .Conversion.Float64 $vuint8 }}
{{ .Conversion.Float64 $vuint16 }}
{{ .Conversion.Float64 $vuint32 }}
{{ .Conversion.Float64 $vuint64 }}
{{ .Conversion.Float64 $vfloat32 }}
{{ .Conversion.Float64 $vfloat64 }}

--- string ---
join
{{- $data := .Collection.Slice 1 99.99 "a" false }}
{{- $sj := .String.Join $data "," }}
join string: {{ $sj }}

split
{{- $ss := .String.Split $sj "," }}
split slice: {{ $ss }}

sprintf
{{- $sp := .String.Sprintf " [%s]     " "abcbabcba" }}
sprintf: {{ $sp }}

trim
{{- $st := .String.Trim $sp " []" }}
trim: <{{ $st }}>

replace
{{- $sr := .String.Replace $st "(.)[b-c]+" "${1}X" }}
replace: {{ $sr }}

--- slice ---
slice init
{{- $slice1 := .Collection.Slice }}
slice1: {{ $slice1 }}

slice init
{{- $slice2 := .Collection.Slice 1 2 3 4 5 }}
slice2: {{ $slice2 }}

slice init
{{- $slice3 := .Collection.Slice 3 4 5 6 7 }}
slice3: {{ $slice3 }}

subslice
{{- $subslice := .Collection.SubSlice $slice3 -1 3 }}
{{ $subslice }}

subslice
{{ .Collection.SubSlice $slice3 1 -1 }}

subslice
{{ .Collection.SubSlice $slice3 -1 -1 }}

subslice
{{ .Collection.SubSlice $slice3 1 3 }}

append
{{- $slice3 := .Collection.Append $slice3 8 9 10 }}
{{ $slice3 }}

concatenate
{{ .Collection.Concatenate $slice2 $slice3 }}

fieldslice
{{- $tpp := .ModelStore.Multi "template_arguments" "" }}
{{- $idfields := .Collection.FieldSlice $tpp.Records "ID" }}
{{- $namefields := .Collection.FieldSlice $tpp.Records "Name" }}
{{- $idfields := .Collection.Sort $idfields "asc" }}
{{- $namefields := .Collection.Sort $namefields "asc" }}
{{ $idfields }}
{{ $namefields }}

sort
{{- $sliceint := .Collection.Slice }}
{{- $v1 := .Conversion.Int 3 }}
{{- $v2 := .Conversion.Int 1 }}
{{- $v3 := .Conversion.Int -5 }}
{{- $v4 := .Conversion.Int 2 }}
{{- $v5 := .Conversion.Int 4 }}
{{- $sliceint := .Collection.Append $sliceint $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint := .Collection.Sort $sliceint "asc" }}
sliceint asc: {{ $sliceint }}
{{- $sliceint := .Collection.Sort $sliceint "desc" }}
sliceint desc: {{ $sliceint }}
{{- $sliceint8 := .Collection.Slice }}
{{- $v1 := .Conversion.Int8 3 }}
{{- $v2 := .Conversion.Int8 1 }}
{{- $v3 := .Conversion.Int8 -5 }}
{{- $v4 := .Conversion.Int8 2 }}
{{- $v5 := .Conversion.Int8 4 }}
{{- $sliceint8 := .Collection.Append $sliceint8 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint8 := .Collection.Sort $sliceint8 "asc" }}
sliceint8 asc: {{ $sliceint8}}
{{- $sliceint8 := .Collection.Sort $sliceint8 "desc" }}
sliceint8 desc: {{ $sliceint8}}
{{- $sliceint16 := .Collection.Slice }}
{{- $v1 := .Conversion.Int16 3 }}
{{- $v2 := .Conversion.Int16 1 }}
{{- $v3 := .Conversion.Int16 -5 }}
{{- $v4 := .Conversion.Int16 2 }}
{{- $v5 := .Conversion.Int16 4 }}
{{- $sliceint16 := .Collection.Append $sliceint16 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint16 := .Collection.Sort $sliceint16 "asc" }}
sliceint16 asc: {{ $sliceint16 }}
{{- $sliceint16 := .Collection.Sort $sliceint16 "desc" }}
sliceint16 desc: {{ $sliceint16 }}
{{- $sliceint32 := .Collection.Slice }}
{{- $v1 := .Conversion.Int32 3 }}
{{- $v2 := .Conversion.Int32 1 }}
{{- $v3 := .Conversion.Int32 -5 }}
{{- $v4 := .Conversion.Int32 2 }}
{{- $v5 := .Conversion.Int32 4 }}
{{- $sliceint32 := .Collection.Append $sliceint32 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint32 := .Collection.Sort $sliceint32 "asc" }}
sliceint32 asc: {{ $sliceint32 }}
{{- $sliceint32 := .Collection.Sort $sliceint32 "desc" }}
sliceint32 desc: {{ $sliceint32 }}
{{- $sliceint64 := .Collection.Slice }}
{{- $v1 := .Conversion.Int64 3 }}
{{- $v2 := .Conversion.Int64 1 }}
{{- $v3 := .Conversion.Int64 -5 }}
{{- $v4 := .Conversion.Int64 2 }}
{{- $v5 := .Conversion.Int64 4 }}
{{- $sliceint64 := .Collection.Append $sliceint64 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceint64 := .Collection.Sort $sliceint64 "asc" }}
sliceint64 asc: {{ $sliceint64 }}
{{- $sliceint64 := .Collection.Sort $sliceint64 "desc" }}
sliceint64 desc: {{ $sliceint64 }}
{{- $sliceuint := .Collection.Slice }}
{{- $v1 := .Conversion.Uint 3 }}
{{- $v2 := .Conversion.Uint 1 }}
{{- $v3 := .Conversion.Uint 5 }}
{{- $v4 := .Conversion.Uint 2 }}
{{- $v5 := .Conversion.Uint 4 }}
{{- $sliceuint := .Collection.Append $sliceuint $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint := .Collection.Sort $sliceuint "asc" }}
sliceuint asc: {{ $sliceuint }}
{{- $sliceuint := .Collection.Sort $sliceuint "desc" }}
sliceuint desc: {{ $sliceuint }}
{{- $sliceuint8 := .Collection.Slice }}
{{- $v1 := .Conversion.Uint8 3 }}
{{- $v2 := .Conversion.Uint8 1 }}
{{- $v3 := .Conversion.Uint8 5 }}
{{- $v4 := .Conversion.Uint8 2 }}
{{- $v5 := .Conversion.Uint8 4 }}
{{- $sliceuint8 := .Collection.Append $sliceuint8 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint8 := .Collection.Sort $sliceuint8 "asc" }}
sliceuint8 asc: {{ $sliceuint8}}
{{- $sliceuint8 := .Collection.Sort $sliceuint8 "desc" }}
sliceuint8 desc: {{ $sliceuint8}}
{{- $sliceuint16 := .Collection.Slice }}
{{- $v1 := .Conversion.Uint16 3 }}
{{- $v2 := .Conversion.Uint16 1 }}
{{- $v3 := .Conversion.Uint16 5 }}
{{- $v4 := .Conversion.Uint16 2 }}
{{- $v5 := .Conversion.Uint16 4 }}
{{- $sliceuint16 := .Collection.Append $sliceuint16 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint16 := .Collection.Sort $sliceuint16 "asc" }}
sliceuint16 asc: {{ $sliceuint16 }}
{{- $sliceuint16 := .Collection.Sort $sliceuint16 "desc" }}
sliceuint16 desc: {{ $sliceuint16 }}
{{- $sliceuint32 := .Collection.Slice }}
{{- $v1 := .Conversion.Uint32 3 }}
{{- $v2 := .Conversion.Uint32 1 }}
{{- $v3 := .Conversion.Uint32 5 }}
{{- $v4 := .Conversion.Uint32 2 }}
{{- $v5 := .Conversion.Uint32 4 }}
{{- $sliceuint32 := .Collection.Append $sliceuint32 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint32 := .Collection.Sort $sliceuint32 "asc" }}
sliceuint32 asc: {{ $sliceuint32 }}
{{- $sliceuint32 := .Collection.Sort $sliceuint32 "desc" }}
sliceuint32 desc: {{ $sliceuint32 }}
{{- $sliceuint64 := .Collection.Slice }}
{{- $v1 := .Conversion.Uint64 3 }}
{{- $v2 := .Conversion.Uint64 1 }}
{{- $v3 := .Conversion.Uint64 5 }}
{{- $v4 := .Conversion.Uint64 2 }}
{{- $v5 := .Conversion.Uint64 4 }}
{{- $sliceuint64 := .Collection.Append $sliceuint64 $v1 $v2 $v3 $v4 $v5 }}
{{- $sliceuint64 := .Collection.Sort $sliceuint64 "asc" }}
sliceuint64 asc: {{ $sliceuint64 }}
{{- $sliceuint64 := .Collection.Sort $sliceuint64 "desc" }}
sliceuint64 desc: {{ $sliceuint64 }}
{{- $slicefloat64 := .Collection.Slice }}
{{- $v1 := .Conversion.Float32 3.3 }}
{{- $v2 := .Conversion.Float32 1 }}
{{- $v3 := .Conversion.Float32 -5.1 }}
{{- $v4 := .Conversion.Float32 2.2 }}
{{- $v5 := .Conversion.Float32 4 }}
{{- $slicefloat64 := .Collection.Append $slicefloat64 $v1 $v2 $v3 $v4 $v5 }}
{{- $slicefloat64 := .Collection.Sort $slicefloat64 "asc" }}
slicefloat32 asc: {{ $slicefloat64 }}
{{- $slicefloat64 := .Collection.Sort $slicefloat64 "desc" }}
slicefloat32 desc: {{ $slicefloat64 }}
{{- $slicefloat64 := .Collection.Slice }}
{{- $v1 := .Conversion.Float64 3.3 }}
{{- $v2 := .Conversion.Float64 1 }}
{{- $v3 := .Conversion.Float64 -5.1 }}
{{- $v4 := .Conversion.Float64 2.2 }}
{{- $v5 := .Conversion.Float64 4 }}
{{- $slicefloat64 := .Collection.Append $slicefloat64 $v1 $v2 $v3 $v4 $v5 }}
{{- $slicefloat64 := .Collection.Sort $slicefloat64 "asc" }}
slicefloat64 asc: {{ $slicefloat64 }}
{{- $slicefloat64 := .Collection.Sort $slicefloat64 "desc" }}
slicefloat64 desc: {{ $slicefloat64 }}
{{- $slicestring := .Collection.Slice }}
{{- $v1 := "3.3" }}
{{- $v2 := "ABC" }}
{{- $v3 := "-5.1" }}
{{- $v4 := "012" }}
{{- $v5 := "def" }}
{{- $slicestring := .Collection.Append $slicestring $v1 $v2 $v3 $v4 $v5 }}
{{- $slicestring := .Collection.Sort $slicestring "asc" }}
slicestring asc: {{ $slicestring}}
{{- $slicestring := .Collection.Sort $slicestring "desc" }}
slicestring desc: {{ $slicestring}}

--- map ---
map
{{- $map1 := .Collection.Map }}
map1: {{ $map1 }}

map init get
{{- $map2 := .Collection.Map 1 "A" 2 "B" }}
map2[1]: {{ $map2.Get 1 }}
map2[2]: {{ $map2.Get 2 }}

map init get
{{- $map3 := .Collection.Map 1 "C" 3 "D" }}
map3[1]: {{ $map3.Get 1 }}
map3[3]: {{ $map3.Get 3 }}

map exists
{{ $map2.Exists 0 }}
{{- $e := $map2.Exists 1 }}
{{ if eq $e true }}TRUE!!{{ else }}FALSE!!{{ end }}

map put
{{- $null := $map1.Put 4 "E" }}
{{- $null := $map1.Put 5 "F" }}
{{- $null := $map1.Put 6 "G" }}
map1[4]: {{ $map1.Get 4 }}
map1[5]: {{ $map1.Get 5 }}
map1[6]: {{ $map1.Get 6 }}

map delete
{{- $null := $map1.Delete 5 }}
map1[3]: {{ $map1.Get 3 }}
map1[4]: {{ $map1.Get 4 }}
map1[5]: {{ $map1.Get 5 }}
map1[6]: {{ $map1.Get 6 }}
map1[7]: {{ $map1.Get 7 }}

map merge
{{- $null :=  $map1.Merge $map2 }}
{{- $null :=  $map1.Merge $map3 }}
map1[0]: {{ $map1.Get 0 }}
map1[1]: {{ $map1.Get 1 }}
map1[2]: {{ $map1.Get 2 }}
map1[3]: {{ $map1.Get 3 }}
map1[4]: {{ $map1.Get 4 }}
map1[5]: {{ $map1.Get 5 }}
map1[6]: {{ $map1.Get 6 }}
map1[7]: {{ $map1.Get 7 }}

map keys
{{- $keys := $map1.Keys }}
{{- $keys := .Collection.Sort $keys "asc" }}
keys of map1: {{ $keys }}

map values
{{- $values := $map1.Values }}
{{- $values := .Collection.Sort $values "asc" }}
values of map1: {{ $values }}

--- model store ---
multi
{{- $m:= .ModelStore.Multi "templates" "preloads=TemplateArguments" }}
{{- $t := index $m.Records 0 }}
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

total, paginaiton
{{- $t := .ModelStore.Multi "/template_arguments" "limit=2&q[name]=%25X" }}
{{ len $t.Records }}
{{ $t.Total }}
{{ $t.CountBeforePagination }}

--- hash ---
hash
{{- $h := .Collection.Hash $s.TemplateArguments "Name" }}
hash[testParameter11]={{ $h.Get "testParameter11" }}
hash[testParameter12]={{ $h.Get "testParameter12" }}

--- sequence ---
sequence
{{- $s := .Collection.Sequence 1 10 }}
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

--- json, yaml ---
{{- $x := .ModelStore.Multi "templates" "preloads=template_arguments&fields=id,name" }}
{{- $x2 := .Collection.Map }}
{{- $x2 := $x2.Put "title" "test" }}
{{- $x2 := $x2.Put "abc" 100 }}
{{- $x3 := $x2.Put "records" $x.Records }}
json
{{ .Conversion.JSONMarshal $x2 "	" }}
yaml
{{ .Conversion.YAMLMarshal $x2 }}
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
{{ .ModelStore.Single "templates/test12/generation" "key_parameter=name" }}
{{ .ModelStore.Single "templates/test12/generation" "key_parameter=name&p[testParameter1X]=999" }}
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

	parametersQuery := map[string]string{
		"query1": "123",
		"query2": "abc",
	}

	parametersByName := map[string]string{
		"query1":        "456",
		"query2":        "def",
		"key_parameter": "name",
	}

	responseText, code := Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id), "generation", parametersQuery), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_1.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%s", template1.Name), "generation", parametersByName), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_2.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id2), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_3.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%s", template2.Name), "generation", parametersByName), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_3.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id3), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_4.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%s", template3.Name), "generation", parametersByName), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_4.txt"))

	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%d", id4), "generation", nil), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_5.txt"))
	responseText, code = Execute(t, http.MethodGet, GenerateSingleResourceURL(server, fmt.Sprintf("templates/%s", template4.Name), "generation", parametersByName), nil)
	CheckResponseText(t, code, http.StatusOK, responseText, LoadExpectation(t, "template/TestTemplate_Functions_5.txt"))
}
