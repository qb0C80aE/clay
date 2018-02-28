package models

import (
	"github.com/qb0C80aE/clay/extensions"
)

// Graphql is the model class what represents graphql query and mutation
type Graphql struct {
	Query    string
	Mutation string
}

// NewGraphqlModel creates a Graphql model instance
func NewGraphqlModel() *Graphql {
	return &Graphql{}
}

var sharedGraphqlModel = NewGraphqlModel()

// SharedGraphqlModel returns the graphql model instance used as a model prototype and type analysis
func SharedGraphqlModel() *Graphql {
	return sharedGraphqlModel
}

func init() {
	extensions.RegisterModel(sharedGraphqlModel, false)
}
