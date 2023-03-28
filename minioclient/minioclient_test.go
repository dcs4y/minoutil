package minioclient

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMinioBucket_UploadFile(t *testing.T) {
	mb, err := GetClient().NewBucket("picture", "")
	if err != nil {
		log.Println("NewBucket", err)
		return
	}
	fPath := "F:\\OneDrive - abc\\图片\\美女壁纸\\4b93593fdd007d319cfec9bf08584018.jpg"
	objectName := fPath[strings.LastIndex(fPath, "\\")+1:]
	size, err := mb.UploadFile("image/"+objectName, fPath,
		map[string]string{
			"name":        objectName,
			"create_Time": fmt.Sprintf("%d", time.Now().Unix()),
		})
	if err != nil {
		log.Println("UploadFile", err)
		return
	}
	log.Println("对象大小:", size)
}

func TestMinioBucket_UploadObject(t *testing.T) {
	mb, err := GetClient().NewBucket("picture", "")
	if err != nil {
		log.Println("NewBucket", err)
		return
	}
	readFile("F:\\OneDrive - abc\\图片", "", func(mb *minioBucket) func(root, path, name string) {
		return func(root, path, name string) {
			fPath := filepath.Join(root, path, name)
			file, err := os.Open(fPath)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			fileStat, err := file.Stat()
			if err != nil {
				fmt.Println(err)
				return
			}
			objectName := path + name
			size, err := mb.UploadObject(objectName, file, fileStat.Size(),
				map[string]string{
					"name":        objectName,
					"create_Time": fmt.Sprintf("%d", time.Now().Unix()),
				})
			if err != nil {
				log.Println("UploadFile", err)
				return
			}
			log.Println("对象大小:", size)
		}
	}(mb))
}

func readFile(root, path string, fun func(root, path, name string)) {
	files, err := ioutil.ReadDir(filepath.Join(root, path))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			fmt.Println(root, path+f.Name()+"/")
			readFile(root, path+f.Name()+"/", fun)
		} else {
			fmt.Println(root, path, f.Name())
			fun(root, path, f.Name())
		}
	}
}

func TestMinioBucket_GetObjectMeta(t *testing.T) {
	mb, err := GetClient().NewBucket("picture", "")
	if err != nil {
		log.Println("NewBucket", err)
		return
	}
	objectName := "image/4b93593fdd007d319cfec9bf08584018.jpg"
	info := mb.GetObjectMeta(objectName)
	if info != nil {
		fmt.Println(*info)
	}
	infoList := mb.GetObjectList("", true)
	for _, info := range infoList {
		fmt.Println(*info)
	}
}

func TestMinioBucket_DownloadFile(t *testing.T) {
	mb, err := GetClient().NewBucket("picture", "")
	if err != nil {
		log.Println("NewBucket", err)
		return
	}
	objectName := "222.jpg"
	err = mb.DownloadFile(objectName, "C:\\Users\\jugg\\Desktop\\temp\\"+objectName)
	if err != nil {
		fmt.Println(err)
	}
}

func TestMinioBucket_DownloadObject(t *testing.T) {
	mb, err := GetClient().NewBucket("picture", "")
	if err != nil {
		log.Println("NewBucket", err)
		return
	}
	objectName := "4K美女图片/000706-161461482650a2.jpg"
	read, err := mb.DownloadObject(objectName)
	if err != nil {
		fmt.Println(err)
		return
	}
	localFile, err := os.Create("D:\\tmp\\image\\" + objectName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if _, err = io.Copy(localFile, read); err != nil {
		fmt.Println(err)
		return
	}
}

func TestMinioBucket_RemoveObject(t *testing.T) {
	mb, err := GetClient().NewBucket("picture", "")
	if err != nil {
		log.Println("NewBucket", err)
		return
	}
	objectName := "111.jpg"
	err = mb.RemoveObject(objectName)
	if err != nil {
		fmt.Println(err)
	}
}

func TestMinioBucket_GetObjectUrl(t *testing.T) {
	// 私有文件
	{
		mb, err := GetClient().NewBucket("picture", "")
		if err != nil {
			log.Println("NewBucket", err)
			return
		}
		objectName := "4K美女图片/000706-161461482650a2.jpg"
		url, err := mb.GetObjectUrl(objectName, "161461482650a2.jpg", 60, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(url.String())
	}
	// 公有文件
	{
		mb, err := GetClient().NewBucket("file", "")
		if err != nil {
			log.Println("NewBucket", err)
			return
		}
		objectName := "cae13b82615b37c430407055ca7b4dbf.jpg"
		url, err := mb.GetObjectUrl(objectName, "", 0, true)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(url.String())
	}
}
