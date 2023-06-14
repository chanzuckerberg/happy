package utils

import (
	"encoding/json"
)

func StructToMap(payload interface{}) (map[string]interface{}, error) {
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(payload)
	err := json.Unmarshal(inrec, &inInterface)
	if err != nil {
		return nil, err
	}
	return inInterface, nil
}
