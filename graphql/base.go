package graphql

import (
	"fmt"
	graphqlGo "github.com/graphql-go/graphql"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/serenize/snaker"
)

var commonArgumentMap = graphqlGo.FieldConfigArgument{
	"limit": &graphqlGo.ArgumentConfig{
		Type: graphqlGo.Int,
	},
	"page": &graphqlGo.ArgumentConfig{
		Type: graphqlGo.Int,
	},
	"last_id": &graphqlGo.ArgumentConfig{
		Type: graphqlGo.Int,
	},
	"order": &graphqlGo.ArgumentConfig{
		Type: graphqlGo.String,
	},
	"sort": &graphqlGo.ArgumentConfig{
		Type: graphqlGo.String,
	},
}

type BaseGraphqlType struct {
	model                     interface{}
	fields                    graphqlGo.Fields
	inputObjectConfigFieldMap graphqlGo.InputObjectConfigFieldMap
	fieldConfigArguments      graphqlGo.FieldConfigArgument
	typeObject                *graphqlGo.Object
	inputTypeObject           *graphqlGo.InputObject
}

// NewBaseGraphqlType creates a new instance of BaseGraphqlType
func NewBaseGraphqlType(model interface{}) *BaseGraphqlType {
	graphqlType := &BaseGraphqlType{
		model: model,
	}
	return graphqlType
}

func (graphqlType *BaseGraphqlType) BuildTypeObject(fieldsEarly graphqlGo.Fields,
	inputObjectConfigFieldMapEarly graphqlGo.InputObjectConfigFieldMap,
	fieldConfigArgument graphqlGo.FieldConfigArgument) {

	graphqlType.fields = fieldsEarly
	graphqlType.inputObjectConfigFieldMap = inputObjectConfigFieldMapEarly

	argumentMap := graphqlGo.FieldConfigArgument{}
	for key, argument := range fieldConfigArgument {
		argumentMap[key] = argument
	}
	for key, commonArgument := range commonArgumentMap {
		argumentMap[key] = commonArgument
	}
	graphqlType.fieldConfigArguments = argumentMap

	resourceName := extensions.RegisteredResourceName(graphqlType.model)

	graphqlType.typeObject = graphqlGo.NewObject(
		graphqlGo.ObjectConfig{
			Name:   snaker.SnakeToCamel(resourceName),
			Fields: fieldsEarly,
		},
	)

	graphqlType.inputTypeObject = graphqlGo.NewInputObject(
		graphqlGo.InputObjectConfig{
			Name:   fmt.Sprintf("%sInput", snaker.SnakeToCamel(resourceName)),
			Fields: inputObjectConfigFieldMapEarly,
		},
	)
}

func (graphqlType *BaseGraphqlType) UpdateTypeObject(fieldsLate graphqlGo.Fields, inputObjectConfigFieldMapLate graphqlGo.InputObjectConfigFieldMap) {
	for key, value := range fieldsLate {
		graphqlType.typeObject.AddFieldConfig(key, value)
	}

	inputTypeObjectName := graphqlType.inputTypeObject.Name()
	InputObjectConfigFieldMap := graphqlGo.InputObjectConfigFieldMap{}

	inputTypeObjectFields := graphqlType.inputTypeObject.Fields()

	for key, value := range inputTypeObjectFields {
		InputObjectConfigFieldMap[key] = &graphqlGo.InputObjectFieldConfig{
			Type: value,
		}
	}
	for key, value := range inputObjectConfigFieldMapLate {
		InputObjectConfigFieldMap[key] = value
	}
	graphqlType.inputTypeObject = graphqlGo.NewInputObject(
		graphqlGo.InputObjectConfig{
			Name:   inputTypeObjectName,
			Fields: InputObjectConfigFieldMap,
		},
	)
}

func (graphqlType *BaseGraphqlType) BuildTypeFieldsEarly() graphqlGo.Fields {
	return graphqlGo.Fields{}
}

func (graphqlType *BaseGraphqlType) BuildInputTypeFieldsEarly() graphqlGo.InputObjectConfigFieldMap {
	return graphqlGo.InputObjectConfigFieldMap{}
}

func (graphqlType *BaseGraphqlType) BuildTypeFieldsLate() graphqlGo.Fields {
	return graphqlGo.Fields{}
}

func (graphqlType *BaseGraphqlType) BuildInputTypeFieldsLate() graphqlGo.InputObjectConfigFieldMap {
	return graphqlGo.InputObjectConfigFieldMap{}
}

func (graphqlType *BaseGraphqlType) BuildTypeArguments() graphqlGo.FieldConfigArgument {
	return graphqlGo.FieldConfigArgument{}
}

func (graphqlType *BaseGraphqlType) ResolveForQuery(p graphqlGo.ResolveParams) (interface{}, error) {
	return nil, nil
}

func (graphqlType *BaseGraphqlType) ResolveForMutation(p graphqlGo.ResolveParams) (interface{}, error) {
	return nil, nil
}

func (graphqlType *BaseGraphqlType) BuiltTypeFields() graphqlGo.Fields {
	return graphqlType.fields
}

func (graphqlType *BaseGraphqlType) BuiltInputTypeFields() graphqlGo.InputObjectConfigFieldMap {
	return graphqlType.inputObjectConfigFieldMap
}

func (graphqlType *BaseGraphqlType) BuiltTypeArguments() graphqlGo.FieldConfigArgument {
	return graphqlType.fieldConfigArguments
}

func (graphqlType *BaseGraphqlType) BuiltTypeObject() *graphqlGo.Object {
	return graphqlType.typeObject
}

func (graphqlType *BaseGraphqlType) BuiltInputTypeObject() *graphqlGo.InputObject {
	return graphqlType.inputTypeObject
}
