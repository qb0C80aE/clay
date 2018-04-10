// +build integration

package integration

import (
	"github.com/qb0C80aE/clay/model"
	"net/http"
	"testing"
)

func TestGetUserDefinedModelDefinitions_Empty(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	responseText, code := Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "user_defined_model_definitions", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, EmptyArrayString, []*model.UserDefinedModelDefinition{})
}

func TestCreateUserDefinedModelDefinition(t *testing.T) {
	server := SetupServer()
	defer server.Close()

	userDefinedModelDefinition1 := model.UserDefinedModelDefinition{
		TypeName:            "Node",
		ResourceName:        "nodes",
		ToBeMigrated:        true,
		IsControllerEnabled: true,
		Fields: []*model.UserDefinedModelFieldDefinition{
			{
				Name:     "ID",
				Tag:      `json:"id" gorm:"primary_key;auto_increment"`,
				TypeName: "int",
				IsSlice:  false,
			},
			{
				Name:     "Name",
				Tag:      `json:"name" gorm:"unique"`,
				TypeName: "string",
				IsSlice:  false,
			},
			{
				Name:     "NodeChildren",
				Tag:      `json:"node_children" gorm:"ForeignKey:NodeID"`,
				TypeName: "Node",
				IsSlice:  true,
			},
		},
	}

	userDefinedModelDefinition2 := model.UserDefinedModelDefinition{
		TypeName:            "NodeAttribute",
		ResourceName:        "node_attributes",
		ToBeMigrated:        true,
		IsControllerEnabled: true,
		Fields: []*model.UserDefinedModelFieldDefinition{
			{
				Name:     "ID",
				Tag:      `json:"id" gorm:"primary_key;auto_increment"`,
				TypeName: "int",
				IsSlice:  false,
			},
			{
				Name:     "Name",
				Tag:      `json:"name" gorm:"unique"`,
				TypeName: "string",
				IsSlice:  false,
			},
			{
				Name:     "NodeID",
				Tag:      `json:"node_id" sql:"type:integer references nodes(id)"`,
				TypeName: "int",
				IsSlice:  false,
			},
			{
				Name:     "Node",
				Tag:      `json:"node" gorm:"ForeignKey:NodeID"`,
				TypeName: "Node",
				IsSlice:  false,
			},
			{
				Name:     "TemplateID",
				Tag:      `json:"template_id" sql:"type:integer references templates(id)" binding:"required,min=1,max=100"`,
				TypeName: "int",
				IsSlice:  false,
			},
			{
				Name:     "Template",
				Tag:      `json:"template" gorm:"ForeignKey:TemplateID"`,
				TypeName: "Template",
				IsSlice:  false,
			},
		},
	}

	responseText, code := Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "user_defined_model_definitions", nil), userDefinedModelDefinition1)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "user_defined_model_definition/TestCreateUserDefinedModelDefinition_1.json"), &model.UserDefinedModelDefinition{})

	responseText, code = Execute(t, http.MethodPost, GenerateMultiResourceURL(server, "user_defined_model_definitions", nil), userDefinedModelDefinition2)
	CheckResponseJSON(t, code, http.StatusCreated, responseText, LoadExpectation(t, "user_defined_model_definition/TestCreateUserDefinedModelDefinition_2.json"), &model.UserDefinedModelDefinition{})

	responseText, code = Execute(t, http.MethodGet, GenerateMultiResourceURL(server, "user_defined_model_definitions", nil), nil)
	CheckResponseJSON(t, code, http.StatusOK, responseText, LoadExpectation(t, "user_defined_model_definition/TestCreateUserDefinedModelDefinition_4.json"), []*model.UserDefinedModelDefinition{})
}
