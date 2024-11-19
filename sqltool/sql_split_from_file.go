package sqltool

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

const DefaultMaxSize = 1000 * DefaultChunkSize
const DefaultChunkSize = 10 * 1024 * 1024 // 10MB

type SplitSqlInfo struct {
	Err  error
	Sql  string
	Part int
}

func splitSqlFromFileAndSaveFile(filename string) error {
	sqlfile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer sqlfile.Close()

	resp := SplitSqlFromFile(sqlfile)
	for splitSqlInfo := range resp {
		if splitSqlInfo.Err != nil {
			return splitSqlInfo.Err
		}

		if err := writeToFile(fmt.Sprintf("./tmp/output_part_%d.sql", splitSqlInfo.Part), splitSqlInfo.Sql); err != nil {
			return errors.Errorf("Failed writeToFile, err: %v", err)
		}
	}

	return nil
}

type SplitSqlFromFileOpts struct {
	maxSize   int
	chunkSize int
}

func WithMaxSize(maxSize int) func(opt *SplitSqlFromFileOpts) {
	return func(opt *SplitSqlFromFileOpts) {
		if maxSize > 0 {
			opt.maxSize = maxSize
		}
	}
}

func WithChunkSize(chunkSize int) func(opt *SplitSqlFromFileOpts) {
	return func(opt *SplitSqlFromFileOpts) {
		if chunkSize > 0 {
			opt.chunkSize = chunkSize
		}
	}
}

// SplitSqlFromFile 支持读sql文件,然后将sql文件拆成 10mb 的sql小文件
func SplitSqlFromFile(sqlfile io.Reader, opts ...func(opt *SplitSqlFromFileOpts)) chan *SplitSqlInfo {
	splitSqlFromFileOpts := &SplitSqlFromFileOpts{
		maxSize:   DefaultMaxSize,
		chunkSize: DefaultChunkSize,
	}
	for _, opt := range opts {
		opt(splitSqlFromFileOpts)
	}

	var (
		maxSize   = splitSqlFromFileOpts.maxSize
		chunkSize = splitSqlFromFileOpts.chunkSize
		err       error
		resp      = make(chan *SplitSqlInfo, 1)
	)

	go func() {
		defer close(resp)

		reader := bufio.NewReader(sqlfile)
		var builder strings.Builder

		part := 1
		chunk := make([]byte, chunkSize)
		for {
			if true {
				n := 0
				n, err = reader.Read(chunk)
				if err != nil && err != io.EOF {
					resp <- &SplitSqlInfo{
						Err: errors.Wrapf(err, "Error reader.Read(chunk)"),
					}
					break
				}
				if n != 0 {
					builder.Write(chunk[:n])
				}
			} else {
				chunk, err = reader.ReadBytes('\n')
				if err != nil && err != io.EOF {
					resp <- &SplitSqlInfo{
						Err: errors.Wrapf(err, "Error reader.ReadBytes('\n')"),
					}
					break
				}

				builder.Write(chunk)
			}

			if builder.Len() >= chunkSize || err == io.EOF {
				oksqls, lastSql := SplitSQLUnsafe(builder.String())
				resp <- &SplitSqlInfo{
					Err:  nil,
					Sql:  strings.Join(oksqls, "\n"),
					Part: part,
				}

				// Clear the builder and reset totalRead
				builder.Reset()
				part++

				if lastSql != "" {
					builder.WriteString(lastSql)
				}

				if err == io.EOF {
					break
				}

				if part > maxSize/chunkSize {
					resp <- &SplitSqlInfo{
						Err:  errors.Errorf("文件过大, 超过 %dG", maxSize/1024/1024/1024),
						Sql:  "",
						Part: part,
					}
					break
				}
			}

		}

		if len(strings.TrimSpace(builder.String())) > 0 {
			resp <- &SplitSqlInfo{
				Err:  nil,
				Sql:  builder.String(),
				Part: part,
			}
		}
	}()

	return resp
}

func writeToFile(filename, content string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = outFile.WriteString(content)
	return err
}
