package utils

import (
	"encoding/json"
	"os"
)

func ReadJSON(v interface{}, path string) error {

	buffer, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	json.Unmarshal(buffer, v)
	return nil
}

func WriteJSON(v interface{}, path string) error {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	json.Unmarshal(buffer, v)
	return nil
}
