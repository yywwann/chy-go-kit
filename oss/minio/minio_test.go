package minio

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/yywwann/chy-go-kit/oss"
	"math"
	"os"
	"sort"
	"testing"
)

var (
	minioConfig = &oss.Config{
		UseSSL:          false,
		AccessKeyId:     "AccessKeyId",
		AccessKeySecret: "AccessKeySecret",
		Bucket:          "Bucket",
		EndPoint:        "EndPoint",
	}
	smFilePath = "smFilePath"
	bgFilePath = "bgFilePath"
	OssMinio   = &Minio{}
)

func TestMain(m *testing.M) {
	var err error
	OssMinio, err = New(minioConfig)
	if err != nil {
		fmt.Println("main err", err)
		return
	}
	os.Exit(m.Run())
}

func TestMinio_Upload_sm(t *testing.T) {
	file, err := os.OpenFile(smFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	url, err := OssMinio.Upload("test_smFile.sh", file)
	require.NoError(t, err)
	t.Log(url)
}

func TestMinio_Upload_bg(t *testing.T) {
	file, err := os.OpenFile(bgFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	url, err := OssMinio.Upload("test_movie_bg.sh", file)
	require.NoError(t, err)
	t.Log(url)
}

func TestMinio_Multipart_bg(t *testing.T) {
	file, err := os.OpenFile(smFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}
	fileStat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	var chunkSize int64 = 1024 * 1024 * 5
	var num int64
	var key string = "xxx.sh"
	var uploadId string
	var parts []oss.Part

	uploadId, err = OssMinio.InitiateMultipartUpload(key)
	fmt.Println("Multipart", "uploadId = ", uploadId)
	require.NoError(t, err)

	num = int64(math.Ceil(float64(fileStat.Size()) / float64(chunkSize)))

	// 执行并发上传段
	partChan := make(chan oss.Part, num)
	var i int64 = 0
	for ; i < num; i++ {
		go func(i int64) {
			b := make([]byte, chunkSize)
			_, _ = file.Seek(i*(chunkSize), 0)
			if len(b) > int(fileStat.Size()-i*chunkSize) {
				b = make([]byte, fileStat.Size()-i*chunkSize)
			}

			file.Read(b)
			r := bytes.NewReader(b)
			etag, err := OssMinio.UploadPart(key, uploadId, r, int32(i+1))
			require.NoError(t, err)
			partChan <- oss.Part{
				ETag:       etag,
				PartNumber: int32(i + 1),
			}

			t.Log("partNumber", i+1, "len", len(b))

		}(i)
	}

	parts = make([]oss.Part, 0, num)
	// 等待上传完成
	for {
		part, ok := <-partChan
		if !ok {
			break
		}
		parts = append(parts, part)

		if len(parts) == int(num) {
			close(partChan)
		}
	}

	sort.Sort(oss.Parts(parts))

	url, err := OssMinio.CompleteMultipartUpload(key, uploadId, parts)
	require.NoError(t, err)
	t.Log(url)
}
