package file

import (
	"os"
)

//格式化目皮后缀，统一加/
func FmtDirSuffix(pathname string) string {
	length := len(pathname)
	suffix := pathname[length-1 : length]
	if suffix != "/" {
		pathname = pathname + "/"
	}
	return pathname
}

//文件是否存在
func File_exists(file string) bool {
	_, err := os.Stat(file)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
