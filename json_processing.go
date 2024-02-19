package main

import (
	"encoding/json"
	"log"
)

func getUUIDFromJson(jsonData []byte) (string, []byte) {
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Fatal(err)
	}

	orderUID := data["order_uid"].(string)
	delete(data, "order_uid")

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return orderUID, jsonData
}

func addUUIDToJson(orderUID string, jsonData []byte) []byte {
	var data map[string]interface{}

	if err := json.Unmarshal(jsonData, &data); err != nil {
		log.Fatal(err)
	}

	data["order_uid"] = orderUID

	resultJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	return resultJSON
}
