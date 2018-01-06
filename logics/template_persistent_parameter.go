package logics

import (
	"github.com/qb0C80aE/clay/extensions"
	"github.com/qb0C80aE/clay/models"
)

type templatePersistentParameterLogic struct {
	*BaseLogic
}

func newTemplatePersistentParameterLogic() *templatePersistentParameterLogic {
	logic := &templatePersistentParameterLogic{
		BaseLogic: NewBaseLogic(
			models.SharedTemplatePersistentParameterModel(),
		),
	}
	return logic
}

var uniqueTemplatePersistentParameterLogic = newTemplatePersistentParameterLogic()

// UniqueTemplatePersistentParameterLogic returns the unique template persistent parameter logic instance
func UniqueTemplatePersistentParameterLogic() extensions.Logic {
	return uniqueTemplatePersistentParameterLogic
}

func init() {
	extensions.RegisterLogic(models.SharedTemplatePersistentParameterModel(), UniqueTemplatePersistentParameterLogic())
	extensions.RegisterDesignAccessor(uniqueTemplatePersistentParameterLogic)
}
