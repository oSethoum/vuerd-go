package utils

import (
	"os"
	"path"
	"vuerd/types"
)

func WriteFile(file types.File, dir *string) error {
	if dir != nil {
		file.Path = path.Join(*dir, file.Path)
	}
	folder := path.Dir(file.Path)
	err := os.MkdirAll(folder, 0666)
	if err != nil {
		return err
	}
	err = os.WriteFile(file.Path, []byte(file.Buffer), 0666)
	if err != nil {
		return err
	}
	return nil
}

func WriteFiles(files []types.File, dir *string) error {
	for _, file := range files {
		if err := WriteFile(file, dir); err != nil {
			return err
		}
	}
	return nil
}
