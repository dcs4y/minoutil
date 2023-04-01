package fileutil

import (
	"os"
	"path/filepath"
	"time"
)

// FileExist 判断文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// CreateDir 创建目录
func CreateDir(path string) error {
	if !FileExist(path) {
		return os.MkdirAll(path, os.ModeType)
	}
	return nil
}

// WriteFile 保存文件
func WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

// ReadFile 读取文件内容
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// GetFileList 获取一个目录下的所有文件列表。includeDir.是否返回文件夹路径。loop.是否递归查询子目录。
func GetFileList(dirPath string, includeDir bool, loop bool) ([]*FileInfo, error) {
	dir, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	var fileList []*FileInfo
	for _, f := range dir {
		if !f.IsDir() || f.IsDir() && includeDir {
			ff, _ := f.Info()
			fileInfo := &FileInfo{
				ParentPath: dirPath,
				Name:       ff.Name(),
				IsDir:      f.IsDir(),
			}
			fileList = append(fileList, fileInfo)
		}
		if f.IsDir() {
			if loop {
				subFileList, _ := GetFileList(filepath.Join(dirPath, f.Name()), includeDir, loop)
				fileList = append(fileList, subFileList...)
			}
		}
	}
	return fileList, nil
}

type FileInfo struct {
	ParentPath string
	Name       string
	IsDir      bool
	Size       int64
	ModifyTime time.Time
}

func (file FileInfo) GetAbsolutePath() string {
	return filepath.Join(file.ParentPath, file.Name)
}
