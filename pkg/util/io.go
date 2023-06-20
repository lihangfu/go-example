package util

import (
	"io"
	"os"
	"path/filepath"
)

// Exists 判断文件或目录是否存在
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// CreatNestedFile 给定path创建文件，如果目录不存在就递归创建
func CreatNestedFile(path string) (*os.File, error) {
	basePath := filepath.Dir(path)
	if !Exists(basePath) {
		err := os.MkdirAll(basePath, 0700)
		if err != nil {
			return nil, err
		}
	}
	return os.Create(path)
}

// IsEmpty 返回给定目录是否为空目录
func IsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)，读取目录下的一个文件名
	if err == io.EOF {
		return true, nil
	}
	return false, err // 非空或者错误
}
