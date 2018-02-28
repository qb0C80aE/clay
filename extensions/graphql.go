package extensions

import (
	"fmt"
	graphqlGo "github.com/graphql-go/graphql"
	"github.com/qb0C80aE/clay/logging"
	"reflect"
)

var graphqlTypes = []GraphqlType{}
var graphqlTypeMap = map[reflect.Type]GraphqlType{}

type GraphqlType interface {
	BuildTypeFieldsEarly() graphqlGo.Fields
	BuildTypeFieldsLate() graphqlGo.Fields
	BuildInputTypeFieldsEarly() graphqlGo.InputObjectConfigFieldMap
	BuildInputTypeFieldsLate() graphqlGo.InputObjectConfigFieldMap
	BuildTypeArguments() graphqlGo.FieldConfigArgument
	BuildTypeObject(fieldsEarly graphqlGo.Fields, inputObjectConfigFieldMapEarly graphqlGo.InputObjectConfigFieldMap, fieldConfigArgument graphqlGo.FieldConfigArgument)
	UpdateTypeObject(fieldsLate graphqlGo.Fields, inputObjectConfigFieldMapLate graphqlGo.InputObjectConfigFieldMap)
	BuiltTypeFields() graphqlGo.Fields
	BuiltInputTypeFields() graphqlGo.InputObjectConfigFieldMap
	BuiltTypeArguments() graphqlGo.FieldConfigArgument
	BuiltTypeObject() *graphqlGo.Object
	BuiltInputTypeObject() *graphqlGo.InputObject
	ResolveForQuery(p graphqlGo.ResolveParams) (interface{}, error)
	ResolveForMutation(p graphqlGo.ResolveParams) (interface{}, error)
}

func RegisterGraphqlType(model interface{}, graphqlType GraphqlType) {
	modelType := ModelType(model)
	graphqlTypes = append(graphqlTypes, graphqlType)
	graphqlTypeMap[modelType] = graphqlType
}

func RegisteredGraphqlTypes() []GraphqlType {
	result := make([]GraphqlType, len(graphqlTypeMap))
	for i, graphqlType := range graphqlTypes {
		result[i] = graphqlType
	}
	return result
}

func RegisteredGraphqlType(model interface{}) (GraphqlType, error) {
	modelType := ModelType(model)
	result, exist := graphqlTypeMap[modelType]
	if !exist {
		logging.Logger().Debugf("the graphql type related to given name %s does not exist", modelType.Name())
		return nil, fmt.Errorf("the graphql type related to given name %s does not exist", modelType.Name())
	}
	return result, nil
}
