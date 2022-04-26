package utils

import (
	"os"
)

func IsPathExists(pathName string) bool {
	_, err := os.Stat(pathName)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
func MakeSurePathExists(dirName string) {
	if IsPathExists(dirName) {
		return
	}
	os.Mkdir(dirName, 0777)
}
