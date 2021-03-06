package utils

import (
	"os"
	"path"
)

// MkdirForFile make directory for file path
func MkdirForFile(filepath string, perm os.FileMode) error {
	dir := path.Dir(filepath)
	return os.MkdirAll(dir, perm)
}

// CreateFileIfNotExists create file if file not exists
func CreateFileIfNotExists(filePath string, flag int, perm os.FileMode) (*os.File, error) {
	if err := MkdirForFile(filePath, 0666); err != nil {
		return nil, err
	}

	if ok, err := PathExists(filePath); err != nil {
		return nil, err
	} else if !ok {
		return os.OpenFile(filePath, os.O_CREATE|flag, perm)
	}

	return os.OpenFile(filePath, flag, perm)
}

// PathExists to judge file's existence
func PathExists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}
