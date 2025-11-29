package utils

import (
	"os"
)

func Exists(path string) bool {
	var _, err = os.Stat(path)
	return err == nil
}

func Filter[T any](source []T, filterFunc func(T) bool) (ret []T) {
	var returnArray = []T{}
	for _, s := range source {
		if filterFunc(s) {
			returnArray = append(ret, s)
		}
	}
	return returnArray
}
