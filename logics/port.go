package logics

import (
	"github.com/jinzhu/gorm"
	"github.com/qb0C80aE/clay/models"
	"strconv"
)

type PortLogic struct {
}

func NewPortLogic() *PortLogic {
	return &PortLogic{}
}

func updateLink(db *gorm.DB, inputPort *models.Port) error {

	port := &models.Port{}

	if err := db.Model(port).Where("destination_port_id = ?", inputPort.ID).Update("destination_port_id", nil).Error; err != nil {
		return err
	}

	if inputPort.DestinationPortID.Valid {
		if err := db.Model(port).Where("id <> ? and destination_port_id = ?", inputPort.ID, inputPort.DestinationPortID.Int64).Update("destination_port_id", nil).Error; err != nil {
			return err
		}
		if err := db.Model(port).Where("id = ?", inputPort.DestinationPortID.Int64).Update("destination_port_id", inputPort.ID).Error; err != nil {
			return err
		}
	}

	return nil

}

func (_ *PortLogic) GetSingle(db *gorm.DB, id string, queryFields string) (interface{}, error) {

	port := &models.Port{}

	if err := db.Select(queryFields).First(port, id).Error; err != nil {
		return nil, err
	}

	return port, nil

}

func (_ *PortLogic) GetMulti(db *gorm.DB, queryFields string) ([]interface{}, error) {

	ports := []*models.Port{}

	if err := db.Select(queryFields).Find(&ports).Error; err != nil {
		return nil, err
	}

	result := make([]interface{}, len(ports))
	for i, data := range ports {
		result[i] = data
	}

	return result, nil

}

func (this *PortLogic) Create(db *gorm.DB, data interface{}) (interface{}, error) {

	port := data.(*models.Port)

	if err := db.Create(&port).Error; err != nil {
		return nil, err
	}

	if err := updateLink(db, port); err != nil {
		return nil, err
	}

	return port, nil
}

func (this *PortLogic) Update(db *gorm.DB, id string, data interface{}) (interface{}, error) {

	port := data.(*models.Port)
	port.ID, _ = strconv.Atoi(id)

	destinationPort := &models.Port{}

	if err := db.Model(destinationPort).Where("destination_port_id = ?", port.ID).Update("destination_port_id", nil).Error; err != nil {
		return nil, err
	}

	if err := db.Save(port).Error; err != nil {
		return nil, err
	}

	if err := updateLink(db, port); err != nil {
		return nil, err
	}

	return port, nil
}

func (_ *PortLogic) Delete(db *gorm.DB, id string) error {

	port := &models.Port{}

	if err := db.First(&port, id).Error; err != nil {
		return err
	}

	if err := db.Delete(&port).Error; err != nil {
		return err
	}

	return nil

}

func (_ *PortLogic) Patch(_ *gorm.DB, _ string, _ string) (interface{}, error) {
	return nil, nil
}

func (_ *PortLogic) Options(db *gorm.DB) error {
	return nil
}
