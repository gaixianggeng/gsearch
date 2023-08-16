package utils

import "os"

// IsFileExist  判断所给路径文件/文件夹是否存在
func IsFileExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}
