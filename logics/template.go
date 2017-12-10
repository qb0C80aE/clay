package logics

import (
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
)

type templateLogic struct {
	*BaseLogic
}

func newTemplateLogic() *templateLogic {
	logic := &templateLogic{
		BaseLogic: NewBaseLogic(
			models.SharedTemplateModel(),
		),
	}
	return logic
}

var uniqueTemplateLogic = newTemplateLogic()

// UniqueTemplateLogic returns the unique template logic instance
func UniqueTemplateLogic() extensions.Logic {
	return uniqueTemplateLogic
}

func init() {
	extensions.RegisterDesignAccessor(uniqueTemplateLogic)
	extensions.RegisterTemplateParameterGenerator(models.SharedTemplateModel(), uniqueTemplateLogic)
}
