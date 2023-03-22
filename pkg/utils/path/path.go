package path

import "os"

func IsPathExist(path string) (res bool) {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func MkdirPath(path string) (err error) {
	err = os.MkdirAll(path, os.ModePerm)
	return err
}

func RemovePath(path string) {
	os.RemoveAll(path)
}
