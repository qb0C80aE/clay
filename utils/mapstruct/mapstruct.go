package mapstruct

import (
	"encoding/json"
)

// MapToStruct maps the map instance into the struct instance using JSON marshal logic
func MapToStruct(m []interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return err
	}
	return nil
}
