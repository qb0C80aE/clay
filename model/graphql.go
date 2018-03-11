package model

import (
	"github.com/qb0C80aE/clay/extension"
)

// Graphql is the model class what represents graphql query and mutation
type Graphql struct {
	Base
	Query    string
	Mutation string
}

// NewGraphql creates a graphql model instance
func NewGraphql() *Graphql {
	return ConvertContainerToModel(&Graphql{}).(*Graphql)
}

func init() {
	extension.RegisterModel(NewGraphql(), false)
}
