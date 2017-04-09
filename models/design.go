package models

// Design is the model class what represents the whole object model store
type Design struct {
	Content map[string]interface{} `json:"content"`
}

// NewDesignModel creates a Design model instance
func NewDesignModel() *Design {
	return &Design{}
}

var designModel = NewDesignModel()

// SharedDesignModel returns the design model instance used as a model prototype and type analysis
func SharedDesignModel() *Design {
	return designModel
}

func init() {
}
