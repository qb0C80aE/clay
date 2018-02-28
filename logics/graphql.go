package logics

import (
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
)

type graphqlLogic struct {
	*BaseLogic
}

func newGraphqlLogic() *graphqlLogic {
	logic := &graphqlLogic{
		BaseLogic: NewBaseLogic(
			models.SharedGraphqlModel(),
		),
	}
	return logic
}

var uniqueGraphqlLogic = newGraphqlLogic()

// UniqueGraphqlLogic returns the unique graphql logic instance
func UniqueGraphqlLogic() extensions.Logic {
	return uniqueGraphqlLogic
}

func init() {
	extensions.RegisterLogic(models.SharedGraphqlModel(), UniqueGraphqlLogic())
	extensions.RegisterDesignAccessor(uniqueGraphqlLogic)
}
