package models

type Design struct {
	Content map[string]interface{} `json:"content"`
}

func init() {
}

var DesignModel = &Design{}
