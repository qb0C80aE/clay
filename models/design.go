package models

import "github.com/qb0C80aE/clay/extensions"

// Design is the model class what represents the whole object model store
type Design struct {
	ID      int                    `json:"-,omitempty" gorm:"primary_key"`
	Content map[string]interface{} `json:"content" gorm:"-"`
}

// NewDesignModel creates a Design model instance
func NewDesignModel() *Design {
	return &Design{}
}

var sharedDesignModel = NewDesignModel()

// SharedDesignModel returns the design model instance used as a model prototype and type analysis
func SharedDesignModel() *Design {
	return sharedDesignModel
}

func init() {
	extensions.RegisterModel(sharedDesignModel, false)
}
