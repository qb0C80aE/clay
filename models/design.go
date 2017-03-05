package models

type Design struct {
	Content map[string]interface{} `json:"content"`
}

var DesignModel = &Design{}

func init() {
}
