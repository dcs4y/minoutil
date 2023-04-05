package minioclient

import (
	"context"
	"github.com/dcs4y/minoutil/logutil"
	"github.com/dcs4y/minoutil/netutil"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mino/common"
	"net/url"
	"strings"
	"sync"
	"time"
)

var clients = make(map[string]*minioClient)

var log = logutil.GetLog("minio")

type MinioConfig struct {
	Endpoint  string `yaml:"endpoint"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	UseSSL    bool   `yaml:"use_ssl"`
}

type minioClient struct {
	*minio.Client
	bucketMap map[string]*minioBucket // 桶名的缓存。key=location:bucketName。
	sync.RWMutex
}

type minioBucket struct {
	Client     *minio.Client
	BucketName string
	Location   string
}

func NewClient(config MinioConfig) (*minioClient, error) {
	// 初使化minio client对象。
	mc, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	client := &minioClient{
		Client:    mc,
		bucketMap: make(map[string]*minioBucket),
	}
	clients[""] = client
	return client, nil
}

func GetClientByName(name string) *minioClient {
	return clients[name]
}

func GetClient() *minioClient {
	return clients[""]
}

func (mc *minioClient) NewBucket(bucketName, location string) (*minioBucket, error) {
	mc.Lock()
	defer mc.Unlock()
	ctx := context.Background()
	if mb := mc.bucketMap[bucketName+":"+location]; mb != nil {
		return mb, nil
	}
	err := mc.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, err := mc.BucketExists(ctx, bucketName)
		if exists {
			log.Println("桶已存在", bucketName+":"+location)
		} else {
			return nil, err
		}
	}
	return &minioBucket{
		Client:     mc.Client,
		BucketName: bucketName,
		Location:   location,
	}, nil
}

// UploadObject 上传对象
func (mb *minioBucket) UploadObject(objectName string, reader io.Reader, objectSize int64, metadata map[string]string) (info minio.UploadInfo, err error) {
	ctx := context.Background()
	return mb.Client.PutObject(ctx, mb.BucketName, objectName, reader, objectSize, minio.PutObjectOptions{ContentType: netutil.GetContentType(objectName), UserMetadata: metadata})
}

// UploadFile 上传文件。根据objectName覆盖文件。
func (mb *minioBucket) UploadFile(objectName, filePath string, metadata map[string]string) (info minio.UploadInfo, err error) {
	ctx := context.Background()
	return mb.Client.FPutObject(ctx, mb.BucketName, objectName, filePath, minio.PutObjectOptions{ContentType: netutil.GetContentType(objectName), UserMetadata: metadata})
}

// DownloadObject 下载对象
func (mb *minioBucket) DownloadObject(objectName string) (io.ReadCloser, error) {
	ctx := context.Background()
	return mb.Client.GetObject(ctx, mb.BucketName, objectName, minio.GetObjectOptions{})
}

// DownloadFile 下载文件
func (mb *minioBucket) DownloadFile(objectName, filePath string) error {
	ctx := context.Background()
	return mb.Client.FGetObject(ctx, mb.BucketName, objectName, filePath, minio.GetObjectOptions{})
}

type objectInfo struct {
	BucketName   string
	ETag         string    // 对象的MD5校验码
	ObjectName   string    // 对象名称
	LastModified time.Time // 对象的最后修改时间
	Size         int64     // 对象的大小
	ContentType  string    // 对象的Content type
	metadata     map[string]string
}

// GetObjectMeta 查询对象元数据。objectInfo为固有属性，metadata为自定义属性(key为仅首字母大写)。
func (mb *minioBucket) GetObjectMeta(objectName string) *objectInfo {
	ctx := context.Background()
	object, err := mb.Client.StatObject(ctx, mb.BucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		log.Println(err)
		return nil
	} else {
		return buildObjectInfo(mb.BucketName, object)
	}
}

// GetObjectList 查询桶所有对象列表。recursive是否递归查询子目录。
func (mb *minioBucket) GetObjectList(namePrefix string, recursive bool) []*objectInfo {
	var objectList []*objectInfo
	ctx := context.Background()
	objectCh := mb.Client.ListObjects(ctx, mb.BucketName, minio.ListObjectsOptions{Prefix: namePrefix, Recursive: recursive})
	for object := range objectCh {
		if object.Err != nil {
			log.Println(object.Err)
		} else {
			objectList = append(objectList, buildObjectInfo(mb.BucketName, object))
		}
	}
	return objectList
}

func buildObjectInfo(bucketName string, object minio.ObjectInfo) *objectInfo {
	info := &objectInfo{
		BucketName: bucketName,
		metadata:   make(map[string]string),
	}
	info.ETag = object.ETag
	info.ObjectName = object.Key
	info.LastModified = object.LastModified
	info.Size = object.Size
	info.ContentType = object.ContentType
	for k, v := range object.Metadata {
		if strings.HasPrefix(k, "X-Amz-Meta-") {
			k = strings.TrimPrefix(k, "X-Amz-Meta-")
			info.metadata[k] = v[0]
		}
	}
	return info
}

// GetObjectUrl 生成分享链接。生成公开桶里文件连接时，second=0。
func (mb *minioBucket) GetObjectUrl(objectName, showName string, second int64, inline bool) (*url.URL, error) {
	if second == 0 {
		config := common.WS.Minio
		scheme := "http"
		if config.UseSSL {
			scheme = "https"
		}
		return url.Parse(scheme + "://" + config.Endpoint + "/" + mb.BucketName + "/" + objectName)
	}
	// 额外的响应头，支持response-expires， response-content-type， response-cache-control， response-content-disposition。
	reqParams := make(url.Values)
	reqParams.Set("response-content-type", netutil.GetContentType(objectName))
	reqParams.Set("response-cache-control", "no-store")
	if inline {
		reqParams.Set("response-content-disposition", "inline;filename=\""+showName+"\"")
	} else {
		reqParams.Set("response-content-disposition", "attachment;filename=\""+showName+"\"")
	}
	ctx := context.Background()
	return mb.Client.PresignedGetObject(ctx, mb.BucketName, objectName, time.Second*time.Duration(second), reqParams)
}

// RemoveObject 删除对象
func (mb *minioBucket) RemoveObject(objectName string) error {
	ctx := context.Background()
	return mb.Client.RemoveObject(ctx, mb.BucketName, objectName, minio.RemoveObjectOptions{})
}
