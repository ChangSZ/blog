package file

import (
	"os"
	"path/filepath"
)

func FileOrDirExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func MkdirAll(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// 如果文件夹不存在，递归创建
		err = os.MkdirAll(path, os.ModePerm)
	}
	return err
}

func OpenOrCreate(file string) (*os.File, error) {
	if err := MakeDirByFile(file); err != nil {
		return nil, err
	}
	return os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
}

func MakeDirByFile(file string) error {
	if !FileOrDirExists(file) {
		dir, _ := filepath.Split(file)
		return MkdirAll(dir)
	}
	return nil
}
