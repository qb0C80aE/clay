package graphql

import (
	"fmt"
	"github.com/gin-gonic/gin"
	graphqlGo "github.com/graphql-go/graphql"
	dbpkg "github.com/qb0C80aE/clay/db"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/logging"
	"github.com/qb0C80aE/clay/model"
	"net/url"
)

type templateGraphqlType struct {
	*BaseGraphqlType
}

func newTemplateGraphqlType() *templateGraphqlType {
	graphqlType := &templateGraphqlType{
		BaseGraphqlType: NewBaseGraphqlType(
			model.NewTemplate(),
		),
	}
	return graphqlType
}

var uniqueTemplateGraphqlType = newTemplateGraphqlType()

// UniqueTemplateGraphqlType returns the unique template graphql type instance
func UniqueTemplateGraphqlType() extension.GraphqlType {
	return uniqueTemplateGraphqlType
}

func (graphqlType *templateGraphqlType) BuildTypeFieldsEarly() graphqlGo.Fields {
	fields := graphqlGo.Fields{
		"id": &graphqlGo.Field{
			Type: graphqlGo.Int,
		},
		"name": &graphqlGo.Field{
			Type: graphqlGo.String,
		},
		"description": &graphqlGo.Field{
			Type: graphqlGo.String,
		},
		"template_content": &graphqlGo.Field{
			Type: graphqlGo.String,
		},
	}

	return fields
}

func (graphqlType *templateGraphqlType) BuildTypeFieldsLate() graphqlGo.Fields {
	fields := graphqlGo.Fields{
		"template_arguments": &graphqlGo.Field{
			Type:    graphqlGo.NewList(UniqueTemplateArgumentGraphqlType().BuiltTypeObject()),
			Args:    UniqueTemplateArgumentGraphqlType().BuiltTypeArguments(),
			Resolve: UniqueTemplateArgumentGraphqlType().ResolveForQuery,
		},
	}

	return fields
}

func (graphqlType *templateGraphqlType) BuildInputTypeFieldsEarly() graphqlGo.InputObjectConfigFieldMap {
	inputObjectConfigFieldMap := graphqlGo.InputObjectConfigFieldMap{
		"id": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.Int,
		},
		"name": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.String,
		},
		"description": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.String,
		},
		"template_content": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.String,
		},
	}

	return inputObjectConfigFieldMap
}

func (graphqlType *templateGraphqlType) BuildInputTypeFieldsLate() graphqlGo.InputObjectConfigFieldMap {
	inputObjectConfigFieldMap := graphqlGo.InputObjectConfigFieldMap{
		"template_arguments": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.NewList(UniqueTemplateArgumentGraphqlType().BuiltInputTypeObject()),
		},
	}

	return inputObjectConfigFieldMap
}

func (graphqlType *templateGraphqlType) BuildTypeArguments() graphqlGo.FieldConfigArgument {
	fieldConfigArgument := graphqlGo.FieldConfigArgument{
		"id": &graphqlGo.ArgumentConfig{
			Type: graphqlGo.Int,
		},
		"name": &graphqlGo.ArgumentConfig{
			Type: graphqlGo.String,
		},
		"description": &graphqlGo.ArgumentConfig{
			Type: graphqlGo.String,
		},
		"template_content": &graphqlGo.ArgumentConfig{
			Type: graphqlGo.String,
		},
	}

	return fieldConfigArgument
}

func (graphqlType *templateGraphqlType) ResolveForQuery(p graphqlGo.ResolveParams) (interface{}, error) {
	urlValues := url.Values{}
	typeFields := graphqlType.BuiltTypeFields()

	for key, _ := range typeFields {
		if value, exists := p.Args[key]; exists {
			urlValues.Add(fmt.Sprintf("q[%s]", key), fmt.Sprintf("%v", value))
		}
	}

	for key, _ := range commonArgumentMap {
		if value, exists := p.Args[key]; exists {
			urlValues.Add(fmt.Sprintf("%s", key), fmt.Sprintf("%v", value))
		}
	}

	parameter, err := dbpkg.NewParameter(urlValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	c := p.Context.(*gin.Context)
	db := dbpkg.Instance(c)

	template := model.NewTemplate()

	db = db.New()
	db = parameter.SortRecords(db)

	db, err = parameter.Paginate(db)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	// db = parameter.SetPreloads(db)
	// db = parameter.FilterFields(db)
	// queryFields := helper.QueryFields(graphqlType.GetModel(), "*")

	result, err := template.GetMulti(db, nil, urlValues, "*")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	c.Set("parents", result)

	return result, err
}

func init() {
	extension.RegisterGraphqlType(model.NewTemplate(), UniqueTemplateGraphqlType())
}
