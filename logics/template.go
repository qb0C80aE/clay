package logics

import (
	"bytes"
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/models"
	"strconv"
	tplpkg "text/template"
)

func GetTemplateExternalParameters(db *gorm.DB, queryFields string) ([]interface{}, error) {

	templateExternalParameters := []*models.TemplateExternalParameter{}

	if err := db.Select(queryFields).Find(&templateExternalParameters).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(templateExternalParameters))
	for i, data := range templateExternalParameters {
		result[i] = data
	}

	return result, nil

}

func GetTemplateExternalParameter(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	templateExternalParameter := &models.TemplateExternalParameter{}

	if err := db.Select(queryFields).First(templateExternalParameter, id).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func CreateTemplateExternalParameter(db *gorm.DB, data interface{}) (interface{}, error) {

	templateExternalParameter := data.(*models.TemplateExternalParameter)

	if err := db.Create(templateExternalParameter).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func UpdateTemplateExternalParameter(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	templateExternalParameter := data.(*models.TemplateExternalParameter)
	templateExternalParameter.ID, _ = strconv.Atoi(id)

	if err := db.Save(&templateExternalParameter).Error; err != nil {
		return nil, err
	}

	return templateExternalParameter, nil

}

func DeleteTemplateExternalParameter(db *gorm.DB, id string) error {

	templateExternalParameter := &models.TemplateExternalParameter{}

	if err := db.First(&templateExternalParameter, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&templateExternalParameter).Error; err != nil {
		return err
	}

	return nil

}

func GetTemplates(db *gorm.DB, queryFields string) ([]interface{}, error) {

	templates := []*models.Template{}

	if err := db.Select(queryFields).Find(&templates).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(templates))
	for i, data := range templates {
		result[i] = data
	}

	return result, nil

}

func GetTemplate(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	template := &models.Template{}

	if err := db.Select(queryFields).First(template, id).Error; err != nil {
		return nil, err
	}

	return template, nil

}

func CreateTemplate(db *gorm.DB, data interface{}) (interface{}, error) {
	template := data.(*models.Template)

	if err := db.Create(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func UpdateTemplate(db *gorm.DB, id string, data interface{}) (interface{}, error) {
	template := data.(*models.Template)
	template.ID, _ = strconv.Atoi(id)

	if err := db.Save(template).Error; err != nil {
		return nil, err
	}

	return template, nil
}

func DeleteTemplate(db *gorm.DB, id string) error {

	template := &models.Template{}

	if err := db.First(&template, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&template).Error; err != nil {
		return err
	}

	return nil

}

func ApplyTemplate(db *gorm.DB, id string, queryFields string) (interface{}, error) {
	type templateParameter struct {
		Nodes                      []*models.Node
		Ports                      []*models.Port
		NodePvs                    []*models.NodePv
		NodeTypes                  []*models.NodeType
		NodeGroups                 []*models.NodeGroup
		Segments                   []*models.Segment
		TemplateExternalParameters map[string]string
	}

	nodePvs := []*models.NodePv{}
	if err := db.Select(queryFields).Find(&nodePvs).Error; err != nil {
		return "", err
	}

	nodeTypes := []*models.NodeType{}
	if err := db.Select(queryFields).Find(&nodeTypes).Error; err != nil {
		return "", err
	}

	nodes := []*models.Node{}
	if err := db.Preload("Ports").Select(queryFields).Find(&nodes).Error; err != nil {
		return "", err
	}

	ports := []*models.Port{}
	if err := db.Preload("Node").Select(queryFields).Find(&ports).Error; err != nil {
		return "", err
	}

	nodeGroups := []*models.NodeGroup{}
	if err := db.Preload("Nodes").Select(queryFields).Find(&nodeGroups).Error; err != nil {
		return "", err
	}

	template := &models.Template{}
	template.ID, _ = strconv.Atoi(id)

	if err := db.Preload("TemplateExternalParameters").Select(queryFields).First(template, template.ID).Error; err != nil {
		return nil, err
	}

	nodeMap := make(map[int]*models.Node)
	portMap := make(map[int]*models.Port)
	consumedPortMap := make(map[int]*models.Port)

	for _, node := range nodes {
		nodeMap[node.ID] = node
	}
	for _, port := range ports {
		portMap[port.ID] = port
	}

	segments := CreateSegments(nodeMap, portMap, consumedPortMap)

	templateExternalParameterMap := make(map[string]string)
	for _, templateExternalParameter := range template.TemplateExternalParameters {
		templateExternalParameterMap[templateExternalParameter.Name] = templateExternalParameter.Value
	}

	parameter := &templateParameter{
		Nodes:                      nodes,
		Ports:                      ports,
		NodePvs:                    nodePvs,
		NodeTypes:                  nodeTypes,
		NodeGroups:                 nodeGroups,
		Segments:                   segments,
		TemplateExternalParameters: templateExternalParameterMap,
	}

	tpl, _ := tplpkg.New("template").Parse(template.TemplateContent)

	var doc bytes.Buffer
	tpl.Execute(&doc, parameter)
	result := doc.String()

	return result, nil
}
