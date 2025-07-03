package common

import (
	"os"
	"strings"
	"time"
)

func GetTimestamp() int64 {
	return time.Now().UnixMilli()
}

// 将文件中的~/xx 替换为绝对路径
func NormalizePath(path string) string {
	if strings.HasPrefix(path, "~") {
		path, _ = Expand(path)
	}
	return path
}

// 判断文件是否存在
func FileExists(filename string) bool {
	filename = NormalizePath(filename)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// 判断文件夹是否存在
func FolderExists(foldername string) bool {
	foldername = NormalizePath(foldername)
	if _, err := os.Stat(foldername); os.IsNotExist(err) {
		return false
	}
	return true
}
