package utils

import (
	"os"
	"path"
	"vuerd/types"
)

func WriteFile(file types.File) error {
	folder := path.Dir(file.Path)
	err := os.MkdirAll(folder, os.ModeAppend)
	if err != nil {
		return err
	}
	err = os.WriteFile(file.Path, []byte(file.Buffer), os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}

func WriteFiles(files []types.File) error {
	for _, file := range files {
		if err := WriteFile(file); err != nil {
			return err
		}
	}
	return nil
}
