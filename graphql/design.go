package graphql

import (
	graphqlGo "github.com/graphql-go/graphql"
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
)

type designGraphqlType struct {
	*BaseGraphqlType
}

func newDesignGraphqlType() *designGraphqlType {
	graphqlType := &designGraphqlType{
		BaseGraphqlType: NewBaseGraphqlType(
			models.SharedDesignModel(),
		),
	}
	return graphqlType
}

var uniqueDesignGraphqlType = newDesignGraphqlType()

// UniqueDesignGraphqlType returns the unique design graphql type instance
func UniqueDesignGraphqlType() extensions.GraphqlType {
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
	extensions.RegisterGraphqlType(models.SharedDesignModel(), UniqueDesignGraphqlType())
}
