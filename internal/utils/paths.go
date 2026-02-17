package utils

import (
	"io"
	"os"
)

func Exists(path string) bool {
	var _, err = os.Stat(path)
	return err == nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func Filter[T any](source []T, filterFunc func(T) bool) (ret []T) {
	var returnArray = []T{}
	for _, s := range source {
		if filterFunc(s) {
			returnArray = append(returnArray, s)
		}
	}
	return returnArray
}
