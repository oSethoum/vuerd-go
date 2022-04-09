package utils

import (
	"encoding/json"
	"os"
	"vuerd/types"
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
	buffer, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return WriteFile(types.File{Path: path, Buffer: string(buffer)}, "")
}
