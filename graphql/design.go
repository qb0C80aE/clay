package graphql

import (
	graphqlGo "github.com/graphql-go/graphql"
	"github.com/qb0C80aE/clay/extension"
	"github.com/qb0C80aE/clay/model"
)

type designGraphqlType struct {
	*BaseGraphqlType
}

func newDesignGraphqlType() *designGraphqlType {
	graphqlType := &designGraphqlType{
		BaseGraphqlType: NewBaseGraphqlType(
			model.NewDesign(),
		),
	}
	return graphqlType
}

var uniqueDesignGraphqlType = newDesignGraphqlType()

// UniqueDesignGraphqlType returns the unique design graphql type instance
func UniqueDesignGraphqlType() extension.GraphqlType {
	return uniqueDesignGraphqlType
}

func (graphqlType *designGraphqlType) BuildTypeFieldsLate() graphqlGo.Fields {
	fields := graphqlGo.Fields{
		"templates": &graphqlGo.Field{
			Type:    graphqlGo.NewList(UniqueTemplateGraphqlType().BuiltTypeObject()),
			Args:    UniqueTemplateGraphqlType().BuiltTypeArguments(),
			Resolve: UniqueTemplateGraphqlType().ResolveForQuery,
		},
		"template_arguments": &graphqlGo.Field{
			Type:    graphqlGo.NewList(UniqueTemplateArgumentGraphqlType().BuiltTypeObject()),
			Args:    UniqueTemplateArgumentGraphqlType().BuiltTypeArguments(),
			Resolve: UniqueTemplateArgumentGraphqlType().ResolveForQuery,
		},
	}

	return fields
}

func init() {
	extension.RegisterGraphqlType(model.NewDesign(), UniqueDesignGraphqlType())
}
