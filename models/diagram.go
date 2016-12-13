package models

type Diagram struct {
	Nodes []*DiagramNode `json:"nodes"`
	Links []*DiagramLink `json:"links"`
}

type DiagramNode struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type DiagramInterface struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

type DiagramMeta struct {
	Interface *DiagramInterface `json:"interface"`
}

type DiagramLink struct {
	Source string       `json:"source"`
	Target string       `json:"target"`
	Meta   *DiagramMeta `json:"meta"`
}

var DiagramModel = &Diagram{}
var DiagramNodeModel = &DiagramNode{}
var DiagramInterfaceModel = &DiagramInterface{}
var DiagramMetaModel = &DiagramMeta{}
var DiagramLinkModel = &DiagramLink{}
