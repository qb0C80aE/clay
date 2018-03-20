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
	"reflect"
)

type templateArgumentGraphqlType struct {
	*BaseGraphqlType
}

func newTemplateArgumentGraphqlType() *templateArgumentGraphqlType {
	graphqlType := &templateArgumentGraphqlType{
		BaseGraphqlType: NewBaseGraphqlType(
			model.NewTemplateArgument(),
		),
	}
	return graphqlType
}

var uniqueTemplateArgumentGraphqlType = newTemplateArgumentGraphqlType()

// UniqueTemplateArgumentGraphqlType returns the unique template argument graphql type instance
func UniqueTemplateArgumentGraphqlType() extension.GraphqlType {
	return uniqueTemplateArgumentGraphqlType
}

func (graphqlType *templateArgumentGraphqlType) BuildTypeFieldsEarly() graphqlGo.Fields {
	fields := graphqlGo.Fields{
		"id": &graphqlGo.Field{
			Type: graphqlGo.Int,
		},
		"name": &graphqlGo.Field{
			Type: graphqlGo.String,
		},
		"template_id": &graphqlGo.Field{
			Type: graphqlGo.Int,
		},
		"description": &graphqlGo.Field{
			Type: graphqlGo.String,
		},
		"type": &graphqlGo.Field{
			Type: graphqlGo.Int,
		},
		"default_value_int": &graphqlGo.Field{
			Type: graphqlGo.Int,
		},
		"default_value_float": &graphqlGo.Field{
			Type: graphqlGo.Float,
		},
		"default_value_bool": &graphqlGo.Field{
			Type: graphqlGo.Boolean,
		},
		"default_value_string": &graphqlGo.Field{
			Type: graphqlGo.String,
		},
	}

	return fields
}

func (graphqlType *templateArgumentGraphqlType) BuildTypeFieldsLate() graphqlGo.Fields {
	fields := graphqlGo.Fields{
		"template": &graphqlGo.Field{
			Type:    UniqueTemplateGraphqlType().BuiltTypeObject(),
			Args:    UniqueTemplateGraphqlType().BuiltTypeArguments(),
			Resolve: UniqueTemplateGraphqlType().ResolveForQuery,
		},
	}

	return fields
}

func (graphqlType *templateArgumentGraphqlType) BuildInputTypeFieldsEarly() graphqlGo.InputObjectConfigFieldMap {
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
		"type": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.Int,
		},
		"default_value_int": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.Int,
		},
		"default_value_float": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.Float,
		},
		"default_value_bool": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.Boolean,
		},
		"default_value_string": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.String,
		},
		"to_be_deleted": &graphqlGo.InputObjectFieldConfig{
			Type: graphqlGo.Boolean,
		},
	}

	return inputObjectConfigFieldMap
}

func (graphqlType *templateArgumentGraphqlType) BuildInputTypeFieldsLate() graphqlGo.InputObjectConfigFieldMap {
	inputObjectConfigFieldMap := graphqlGo.InputObjectConfigFieldMap{
		"template": &graphqlGo.InputObjectFieldConfig{
			Type: UniqueTemplateArgumentGraphqlType().BuiltInputTypeObject(),
		},
	}

	return inputObjectConfigFieldMap
}

func (graphqlType *templateArgumentGraphqlType) BuildTypeArguments() graphqlGo.FieldConfigArgument {
	fieldConfigArgument := graphqlGo.FieldConfigArgument{
		"id": &graphqlGo.ArgumentConfig{
			Type: graphqlGo.Int,
		},
		"name": &graphqlGo.ArgumentConfig{
			Type: graphqlGo.String,
		},
	}

	return fieldConfigArgument
}

func (graphqlType *templateArgumentGraphqlType) ResolveForQuery(p graphqlGo.ResolveParams) (interface{}, error) {
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

	fmt.Println(p.Source)
	fmt.Println(urlValues)

	parameter, err := dbpkg.NewParameter(urlValues)
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}



	c := p.Context.(*gin.Context)
	db := dbpkg.Instance(c)

	parents, _ := c.Get("parents")
	fmt.Printf("parents: %v", parents)

	templateArgument := model.NewTemplateArgument()

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

	result, err := templateArgument.GetMulti(db, nil, urlValues, "*")
	if err != nil {
		logging.Logger().Debug(err.Error())
		return nil, err
	}

	return result, err
}

func init() {
	extension.RegisterGraphqlType(model.NewTemplateArgument(), UniqueTemplateArgumentGraphqlType())
}
