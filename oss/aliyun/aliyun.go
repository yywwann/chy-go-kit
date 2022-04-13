package aliyun

import (
	"fmt"
	"io"
	"sort"

	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pkg/errors"

	"github.com/yywwann/chy-go-kit/oss"
)

const (
	Cloud = "aliyun"
)

type Aliyun struct {
	cfg    *oss.Config
	client *alioss.Client
}

func New(cfg *oss.Config) (*Aliyun, error) {
	ali := &Aliyun{
		cfg:    cfg,
		client: nil,
	}

	// 创建并维护客户端
	err := ali.createClient()
	if err != nil {
		return nil, err
	}

	return ali, nil
}

func (ali *Aliyun) createClient() (err error) {
	ali.client, err = alioss.New(
		ali.cfg.EndPoint,
		ali.cfg.AccessKeyId,
		ali.cfg.AccessKeySecret,
	)
	if err != nil {
		return errors.WithMessage(err, "ali.oss.New()")
	}

	return nil
}

// Upload   上传文件
// key string 文件名, 包括路径 "dtalk/2021-07-01/1.jpg"
// body io.Reader 文件内容.
// url string 文件资源链接.
// error it's nil if no error, otherwise it's an error object.
func (ali *Aliyun) Upload(key string, body io.Reader) (url string, err error) {
	bucket, err := ali.client.Bucket(ali.cfg.Bucket)
	if err != nil {
		return "", errors.WithMessage(err, "client.Bucket()")
	}

	// 指定存储类型为标准存储，缺省也为标准存储。
	storageType := alioss.ObjectStorageClass(alioss.StorageStandard)

	// 指定存储类型为归档存储。
	// storageType := alioss.ObjectStorageClass(alioss.StorageArchive)

	// 指定访问权限为公共读，缺省为继承bucket的权限。
	objectAcl := alioss.ObjectACL(alioss.ACLPublicRead)

	// 上传字符串。
	err = bucket.PutObject(key, body, storageType, objectAcl)
	if err != nil {
		return "", errors.WithMessage(err, "bucket.PutObject()")
	}

	return ali.parseUrl(key), nil
}

// InitiateMultipartUpload 初始化分段上传任务
// 使用分段上传方式传输数据前，必须先通知OBS初始化一个分段上传任务。
// 该操作会返回一个OBS服务端创建的全局唯一标识（Upload ID），用于标识本次分段上传任务。
// 您可以根据这个唯一标识来发起相关的操作，如取消分段上传任务、列举分段上传任务、列举已上传的段等。
//
// key string
//
// uploadId string
// err 		error
func (ali *Aliyun) InitiateMultipartUpload(key string) (uploadId string, err error) {
	bucket, err := ali.client.Bucket(ali.cfg.Bucket)
	if err != nil {
		return "", errors.WithMessage(err, "client.Bucket()")
	}

	imur, err := bucket.InitiateMultipartUpload(key)
	if err != nil {
		return "", errors.WithMessage(err, "bucket.InitiateMultipartUpload")
	}

	return imur.UploadID, nil
}

// UploadPart 上传段
// 初始化一个分段上传任务之后，可以根据指定的对象名和Upload ID来分段上传数据。
// 每一个上传的段都有一个标识它的号码——分段号（Part Number，范围是1~10000）。
// 对于同一个Upload ID，该分段号不但唯一标识这一段数据，也标识了这段数据在整个对象内的相对位置。
// 如果您用同一个分段号上传了新的数据，那么OBS上已有的这个段号的数据将被覆盖。
// 除了最后一段以外，其他段的大小范围是100KB~5GB；最后段大小范围是0~5GB。
// 每个段不需要按顺序上传，甚至可以在不同进程、不同机器上上传，OBS会按照分段号排序组成最终对象。
//
// key 			string
// uploadId 	string
// body 		io.Reader
// partNumber 	int32
// offset 		int64
// partSize 	int64
//
// ETag	string
// err	error
func (ali *Aliyun) UploadPart(key, uploadId string, body io.Reader, partNumber int32) (ETag string, err error) {
	bucket, err := ali.client.Bucket(ali.cfg.Bucket)
	if err != nil {
		return "", errors.WithMessage(err, "client.Bucket")
	}

	imur := alioss.InitiateMultipartUploadResult{
		Key:      key,
		UploadID: uploadId,
		Bucket:   ali.cfg.Bucket,
	}

	size, err := oss.GetReaderLen(body)
	if err != nil {
		return
	}

	part, err := bucket.UploadPart(imur, body, size, int(partNumber))
	if err != nil {
		return "", errors.WithMessage(err, "bucket.UploadPart")
	}
	return part.ETag, nil
}

// CompleteMultipartUpload 合并段
// 所有分段上传完成后，需要调用合并段接口来在OBS服务端生成最终对象。
// 在执行该操作时，需要提供所有有效的分段列表（包括分段号和分段ETag值）；
// OBS收到提交的分段列表后，会逐一验证每个段的有效性。当所有段验证通过后，OBS将把这些分段组合成最终的对象。
//
// key 		string
// uploadId string
// parts 	[]model.Part
//
// url string
// err error
func (ali *Aliyun) CompleteMultipartUpload(key, uploadId string, parts []oss.Part) (url string, err error) {
	bucket, err := ali.client.Bucket(ali.cfg.Bucket)
	if err != nil {
		return "", errors.WithMessage(err, "client.Bucket")
	}

	imur := alioss.InitiateMultipartUploadResult{
		Key:      key,
		UploadID: uploadId,
		Bucket:   ali.cfg.Bucket,
	}

	ossParts := make([]alioss.UploadPart, len(parts), len(parts))
	sort.Sort(oss.Parts(parts))
	for i := range parts {
		ossParts[i] = alioss.UploadPart{
			PartNumber: int(parts[i].PartNumber),
			ETag:       parts[i].ETag,
		}
	}

	_, err = bucket.CompleteMultipartUpload(imur, ossParts)
	if err != nil {
		return "", errors.WithMessage(err, "bucket.CompleteMultipartUpload")
	}
	return ali.parseUrl(key), nil
}

// AbortMultipartUpload 取消分段上传任务
// 分段上传任务可以被取消，当一个分段上传任务被取消后，就不能再使用其Upload ID做任何操作，已经上传段也会被OBS删除。
// 采用分段上传方式上传对象过程中或上传对象失败后会在桶内产生段，这些段会占用您的存储空间，您可以通过取消该分段上传任务来清理掉不需要的段，节约存储空间。
//
// key 		string
// uploadId string
//
// err error
func (ali *Aliyun) AbortMultipartUpload(key, uploadId string) error {
	bucket, err := ali.client.Bucket(ali.cfg.Bucket)
	if err != nil {
		return errors.WithMessage(err, "client.Bucket")
	}

	imur := alioss.InitiateMultipartUploadResult{
		Key:      key,
		UploadID: uploadId,
		Bucket:   ali.cfg.Bucket,
	}

	err = bucket.AbortMultipartUpload(imur)
	if err != nil {
		return errors.WithMessage(err, "bucket.AbortMultipartUpload")
	}
	return nil
}

// parseUrl 获取对象在阿里云OSS上的完整访问URL
func (ali *Aliyun) parseUrl(key string) string {
	return fmt.Sprintf("https://%s.%s/%s", ali.cfg.Bucket, ali.cfg.EndPoint, key)
}
