package oss

import (
	"bytes"
	"io"
)

type OSS interface {
	Upload(key string, body io.Reader) (url string, err error)
	InitiateMultipartUpload(key string) (uploadId string, err error)
	UploadPart(key, uploadId string, body io.Reader, partNumber int32) (ETag string, err error)
	CompleteMultipartUpload(key, uploadId string, parts []Part) (url string, err error)
	AbortMultipartUpload(key, uploadId string) error
}

type Config struct {
	UseSSL          bool   // 是否使用安全配置（用于minio和local云服务商模式）
	Cloud           string // 云服务商（当前支持aliyun、huawei、tencent、minio）
	AccessKeyId     string
	AccessKeySecret string
	Bucket          string
	EndPoint        string
}

type Part struct {
	// 段数据的MD5值
	ETag string `json:"ETag" form:"ETag" binding:"required"`
	// 分段序号, 范围是1~10000
	PartNumber int32 `json:"partNumber" form:"partNumber" binding:"required"`
}

// Parts part数组
type Parts []Part

func (p Parts) Len() int {
	return len(p)
}
func (p Parts) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p Parts) Less(i, j int) bool {
	return p[i].PartNumber < p[j].PartNumber
}

func GetSize(stream io.Reader) int64 {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return int64(buf.Len())
}
