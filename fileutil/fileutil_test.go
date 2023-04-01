package fileutil

import (
	"fmt"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	path := "\\Log.txt"
	data := "Hello World!"
	err := WriteFile(path, []byte(data), os.ModeType)
	if err != nil {
		fmt.Println(err)
	}
	b, err := ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(b))
}

func TestGetFileList(t *testing.T) {
	fileList, err := GetFileList("/", true, false)
	if err != nil {
		fmt.Println(err)
	}
	for _, fileInfo := range fileList {
		fmt.Println(fileInfo.GetAbsolutePath())
	}
}
