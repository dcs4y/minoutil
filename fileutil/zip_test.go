package fileutil

import (
	"fmt"
	"testing"
)

func TestZip(t *testing.T) {
	err := Zip("D:\\usr\\local\\temp", []string{}, "中文.zip")
	if err != nil {
		fmt.Println(err)
	}
}

func TestUnzip(t *testing.T) {
	fmt.Println(Unzip("D:\\usr\\local\\zip\\中文.zip", "D:\\usr\\local\\zip"))
}
