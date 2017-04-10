package models

type Design struct {
	Content map[string]interface{} `json:"content"`
}

func NewDesignModel() *Design {
	return &Design{}
}

var designModel = NewDesignModel()

func SharedDesignModel() *Design {
	return designModel
}

func init() {
}
