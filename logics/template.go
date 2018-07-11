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

func (receiver *templateLogic) GetSequenceNumber() int {
	return 1
}

var uniqueTemplateLogic = newTemplateLogic()

// UniqueTemplateLogic returns the unique template logic instance
func UniqueTemplateLogic() extensions.Logic {
	return uniqueTemplateLogic
}

func init() {
	extensions.RegisterLogic(models.SharedTemplateModel(), UniqueTemplateLogic())
	extensions.RegisterDesignAccessor(uniqueTemplateLogic)
}
