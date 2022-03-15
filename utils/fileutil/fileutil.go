package fileutil

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileExist 判断文件是否存在
func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

// CreateDir 创建目录
func CreateDir(path string) {
	if !FileExist(path) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// Zip 压缩zipPath目录下的文件，存储到zipPath目录或outputPath位置。
func Zip(zipPath string, includeFiles []string, outputPath string) error {
	outputFileName := outputPath
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(zipPath, outputPath)
	}
	if len(includeFiles) == 0 {
		des, err := os.ReadDir(zipPath)
		if err != nil {
			return err
		}
		for _, de := range des {
			if de.Name() == outputFileName {
				return errors.New("文件[" + outputFileName + "]已存在！")
			}
			includeFiles = append(includeFiles, de.Name())
		}
	}
	newZipFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer newZipFile.Close()
	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()
	// Add files to zip
	for _, file := range includeFiles {
		if err = addFileToZip(zipWriter, zipPath, file); err != nil {
			return err
		}
	}
	return nil
}

func addFileToZip(zipWriter *zip.Writer, zipPath, filename string) error {
	fileToZip, err := os.Open(filepath.Join(zipPath, filename))
	if err != nil {
		return err
	}
	defer fileToZip.Close()
	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}
	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename
	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// Unzip 解压
func Unzip(src string, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()
	for _, f := range r.File {
		// Store filename/path for returning and using later on
		fPath := filepath.Join(dest, f.Name)
		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fPath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fPath)
		}
		filenames = append(filenames, fPath)
		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}
		// Make File
		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return filenames, err
		}
		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}
		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		_, err = io.Copy(outFile, rc)
		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()
		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
