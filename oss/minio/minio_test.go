package minio

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/yywwann/chy-go-kit/oss"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

var (
	minioConfig = &oss.Config{
		UseSSL:          false,
		AccessKeyId:     "ikvVLFicIG",
		AccessKeySecret: "ikvVLFicIG",
		Bucket:          "fe-static-resources",
		EndPoint:        "minio-fygs.seenew.info:180",
	}
	smFilePath = "/Users/xiniu/install.sh"
	bgFilePath = "/Users/xiniu/Documents/install/thunder_4.2.1.65254.dmg"
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
	defer file.Close()

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

//获取指定目录下的所有文件,包含子目录下的文件

type File struct {
	Path string
	Key  string
	Size int64
}

func GetAllFiles(dirPth string) (files []File, err error) {
	var dirs []string
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	//PthSep := string(os.PathSeparator)
	for _, fi := range dir {
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, filepath.Join(dirPth, fi.Name()))
		} else {
			files = append(files, File{
				Path: dirPth,
				Key:  fi.Name(),
				Size: fi.Size(),
			})
		}
	}

	// 读取子目录下文件
	for _, table := range dirs {
		temp, _ := GetAllFiles(table)
		for _, temp1 := range temp {
			files = append(files, temp1)
		}
	}

	return files, nil
}

func TestMinio_Upload_Folder(t *testing.T) {
	basePath := "/Users/xxx/Documents/"
	baseDir := "深入剖析Kubernetes/代码/"
	files, err := GetAllFiles(basePath + baseDir)
	if err != nil {
		t.Log(err)
		return
	}
	for _, fileInfo := range files {
		file, err := os.OpenFile(filepath.Join(fileInfo.Path, fileInfo.Key), os.O_RDONLY, os.ModePerm)
		if err != nil {
			t.Log(err)
			continue
		}
		defer file.Close()

		url, err := OssMinio.Upload(strings.TrimPrefix(filepath.Join(fileInfo.Path, fileInfo.Key), basePath), file)
		if err != nil {
			t.Log(err)
			return
		}
		t.Log(url)
	}
}
