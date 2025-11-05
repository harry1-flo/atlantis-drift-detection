package utils

import (
	"os"
	"path/filepath"
)

func GetPwd() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return dir, nil
}

// testcase/*.txt
func MustGetCurrentFileUseTest(prfixPath string, filename string) []byte {
	dir, err := GetPwd()
	if err != nil {
		panic(err)
	}
	filePath := filepath.Join(dir, prfixPath, "testcase", filename)
	content, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	return content
}
