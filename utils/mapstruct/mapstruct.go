package mapstruct

import (
	"encoding/json"
	"fmt"
)

func MapToStruct(m []interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return err
	}
	fmt.Println(string(tmp))
	err = json.Unmarshal(tmp, val)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
