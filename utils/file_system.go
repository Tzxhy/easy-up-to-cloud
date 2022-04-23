package utils

import (
	"log"
	"os"
)

func IsPathExists(pathName string) bool {
	info, err := os.Stat(pathName)
	log.Print(info)
	if err == nil {
		log.Print("存在111")
		return true
	}
	if os.IsNotExist(err) {
		log.Print("不存在")
		return false
	}
	return false
}
func MakeSurePathExists(dirName string) {
	if IsPathExists(dirName) {
		log.Print("存在")
		return
	}
	os.Mkdir(dirName, 0666)
}
