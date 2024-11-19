package sqltool

import (
	"os"
	"strings"

	"github.com/pkg/errors"
)

func EnsureSemicolon(s string) string {
	if !strings.HasSuffix(s, ";") {
		s += ";"
	}
	return s
}

func SplitSQLWithFile(filename string) ([]string, error) {
	sql, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	oksqls, lastSql := SplitSQLUnsafe(string(sql))
	if strings.TrimSpace(lastSql) != "" {
		return nil, errors.Errorf("sql不完整, lastSql: %s", lastSql)
	}

	return oksqls, nil
}

// SplitSQLUnsafe
// 人工模拟分割, 不安全
// 分割 sql 语句, sql语句可能不完整
func SplitSQLUnsafe(sql string) ([]string, string) {
	var (
		startIndex int
		sqls       []string
		lastSql    string
	)

	for i := 0; i < len(sql); i++ {
		// 判断 comment
		if sql[i] == '\'' { // ''
			nextIdx := 0
			tmpi := i
			for {
				if tmpi+1 >= len(sql) {
					goto end
				}

				nextIdx = strings.Index(sql[tmpi+1:], "'")
				if nextIdx == -1 {
					goto end
				}

				tmpi += nextIdx + 1
				if sql[tmpi-1] == '\\' {
					count := 0
					for ttmpi := tmpi - 1; ttmpi >= 0 && sql[ttmpi] == '\\'; ttmpi-- {
						count++
					}
					if count%2 == 0 {
						break
					} else {
						continue
					}
				} else {
					break
				}
			}

			i = tmpi
		} else if sql[i] == '"' { // ""
			nextIdx := 0
			tmpi := i
			for {
				if tmpi+1 >= len(sql) {
					goto end
				}

				nextIdx = strings.Index(sql[tmpi+1:], "\"")
				if nextIdx == -1 {
					goto end
				}

				tmpi += nextIdx + 1
				if sql[tmpi-1] == '\\' {
					count := 0
					for ttmpi := tmpi - 1; ttmpi >= 0 && sql[ttmpi] == '\\'; ttmpi-- {
						count++
					}
					if count%2 == 0 {
						break
					} else {
						continue
					}
				} else {
					break
				}
			}

			i = tmpi
		} else if sql[i] == '`' { // ``
			nextIdx := 0
			tmpi := i
			for {
				if tmpi+1 >= len(sql) {
					goto end
				}

				nextIdx = strings.Index(sql[tmpi+1:], "`")
				if nextIdx == -1 {
					goto end
				}

				tmpi += nextIdx + 1
				if sql[tmpi-1] == '\\' {
					count := 0
					for ttmpi := tmpi - 1; ttmpi >= 0 && sql[ttmpi] == '\\'; ttmpi-- {
						count++
					}
					if count%2 == 0 {
						break
					} else {
						continue
					}
				} else {
					break
				}
			}

			i = tmpi
		} else if sql[i] == '/' && i+1 < len(sql) && sql[i+1] == '*' { // /* */
			nextIdx := strings.Index(sql[i+1:], "*/")
			if nextIdx == -1 {
				break
			}
			i += strings.Index(sql[i+1:], "*/") + 1
		} else if sql[i] == '-' && i+2 < len(sql) && sql[i+1] == '-' && sql[i+2] == ' ' { // --
			nextIdx := strings.Index(sql[i+1:], "\n")
			if nextIdx == -1 {
				break
			}
			i += strings.Index(sql[i+1:], "\n") + 1
		} else if sql[i] == '#' { // #
			nextIdx := strings.Index(sql[i+1:], "\n")
			if nextIdx == -1 {
				break
			}
			i += strings.Index(sql[i+1:], "\n") + 1
		} else if sql[i] == ';' {
			sqls = append(sqls, strings.TrimSpace(sql[startIndex:i+1]))
			startIndex = i + 1
		}
	}
end:
	if startIndex != len(sql) {
		lastSql = sql[startIndex:]
	}

	return sqls, lastSql
}
