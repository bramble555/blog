package pkg

import (
	"log"
	"os"
	"path/filepath"
)

// 判断目录是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		log.Println(err)
		return false
	}
	return true
}
func CreateFolder(filePath string) {
	ok := pathExists(filePath)
	// 目录不存在就创建目录
	if !ok {
		err := os.MkdirAll(filePath, 0777)
		if err != nil {
			log.Fatalf("err:%s\n", err.Error())
		}
	}
}

// filePath 是文件路径，fileName 是文件名字(需要带扩展名，并且不能带文件路径)
func CreateFile(filePath string, fileName string) *os.File {
	ok := pathExists(filePath)
	// 目录不存在就创建目录
	if !ok {
		err := os.MkdirAll(filePath, 0777)
		if err != nil {
			log.Fatalf("err:%s\n", err.Error())
		}
	}

	// 输出到哪
	FilePath := filepath.Join(filePath, fileName)
	file, err := os.OpenFile(FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("err:%s\n", err.Error())
	}
	log.Printf("file:%v\n", file)
	return file
}
