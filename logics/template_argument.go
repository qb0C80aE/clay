package logics

import (
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
)

type templateArgumentLogic struct {
	*BaseLogic
}

func newTemplateArgumentLogic() *templateArgumentLogic {
	logic := &templateArgumentLogic{
		BaseLogic: NewBaseLogic(
			models.SharedTemplateArgumentModel(),
		),
	}
	return logic
}

var uniqueTemplateArgumentLogic = newTemplateArgumentLogic()

// UniqueTemplateArgumentLogic returns the unique template argument logic instance
func UniqueTemplateArgumentLogic() extensions.Logic {
	return uniqueTemplateArgumentLogic
}

func init() {
	extensions.RegisterLogic(models.SharedTemplateArgumentModel(), UniqueTemplateArgumentLogic())
	extensions.RegisterDesignAccessor(uniqueTemplateArgumentLogic)
}
