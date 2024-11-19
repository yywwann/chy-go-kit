package sqltool

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"unicode"

	"github.com/pkg/errors"
	"github.com/xwb1989/sqlparser"
	"github.com/yywwann/chy-go-kit/convert"
)

const (
	UnKnow = iota
	DML
	DDL
	Select
	DbDDL // create、drop databse
	Use
)

var SyntaxTypeDesc = map[int]string{
	UnKnow: "未知",
	DML:    "DML",
	DDL:    "DDL",
	Select: "SELECT",
	DbDDL:  "DbDDL",
	Use:    "USE",
}

var (
	reDbDDL  = regexp.MustCompile("(?i)^create( )*database|drop( )*database")
	reDDL    = regexp.MustCompile("(?i)^alter|^create|^drop|^rename|^truncate")
	reDML    = regexp.MustCompile("(?i)^call|^delete|^do|^handler|^insert|^load\\s+data|^load\\s+xml|^replace|^update")
	reSelect = regexp.MustCompile("(?i)^select|^show|^explain")
	reUse    = regexp.MustCompile("(?i)^use")
	reAlter  = regexp.MustCompile("(?i)^alter\\s*table")

	reMap = map[int]*regexp.Regexp{
		DbDDL:  reDbDDL,
		DDL:    reDDL,
		DML:    reDML,
		Select: reSelect,
		Use:    reUse,
	}
)

// GetSyntaxType 获取sql类型，1为dml 2为ddl 3为select 4为DbDDL 0为未知
func GetSyntaxType(sql string, parser bool) int {
	sqlType := UnKnow
	if parser {
		sqlType = parserType(sql)
	}

	if sqlType == UnKnow {
		sqlType = regexpType(sql)
	}

	return sqlType
}

func parserType(sql string) int {
	sql = warpSqlChinese2(sql)
	sql = StripLeadingComments(sql)
	st, err := sqlparser.Parse(sql)
	if err != nil {
		return UnKnow
	}
	switch st.(type) {
	case *sqlparser.DDL:
		return DDL
	case *sqlparser.Select:
		return Select
	case *sqlparser.DBDDL:
		return DbDDL
	case *sqlparser.Update, *sqlparser.Insert, *sqlparser.Delete:
		return DML
	case *sqlparser.Use:
		return Use
	default:
		return UnKnow
	}
}

func regexpType(sql string) int {
	sql = StripLeadingComments(sql)
	for k, v := range reMap {
		if len(v.FindStringIndex(sql)) == 2 {
			return k
		}
	}

	return UnKnow
}

func StripLeadingComments(sql string) string {
	return sqlparser.StripLeadingComments(sql)
}

var (
	reGetUpdatePrefix  = regexp.MustCompile("(?i)^update[\\s\\w.`]{3,}set")
	reUpdate2Select    = regexp.MustCompile("(?i)^update\\s+([^\\s]+).*(where.*)")
	selectSubstitution = "select count(*) from $1 $2"
)

// GenUpdate2Select
// update <table> set <> where ..... -> select count(*) from <table> where ...
// sqlContent 必须全小写，没有前导注释
func GenUpdate2Select(sqlContent string) (string, error) {
	return genUpdate2SelectV2(sqlContent)
}

// genUpdate2Select
// update <table> set <> where ..... -> select count(*) from <table> where ...
// sqlContent 必须全小写，没有前导注释
func genUpdate2Select(sqlContent string) (string, error) {
	if reGetUpdatePrefix.MatchString(sqlContent) != true {
		return "", errors.New("sqlContent 不是 update 语句")
	}

	sqlContent = strings.ReplaceAll(sqlContent, "\n", " ")
	sqlContent = reUpdate2Select.ReplaceAllString(sqlContent, selectSubstitution)

	return sqlContent, nil
}

// genUpdate2SelectV2
// update <table> set <> where ..... -> select count(*) from <table> where ...
// sqlContent 必须全小写，没有前导注释
func genUpdate2SelectV2(sql string) (string, error) {
	// 解析 SQL
	sql = StripLeadingComments(sql)
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return "", err
	}

	// 如果解析结果不是 UPDATE 语句，返回错误
	updateStmt, ok := stmt.(*sqlparser.Update)
	if !ok {
		return "", errors.Errorf("%s is not an update statement", sql)
	}

	// 构造 SELECT 语句
	table := updateStmt.TableExprs
	where := updateStmt.Where

	selectStmt := &sqlparser.Select{
		Distinct: "",
		Where:    where,
		SelectExprs: sqlparser.SelectExprs{
			&sqlparser.AliasedExpr{
				Expr: &sqlparser.FuncExpr{
					Qualifier: sqlparser.TableIdent{},
					Name:      sqlparser.NewColIdent("count"),
					Distinct:  false,
					Exprs: sqlparser.SelectExprs{
						&sqlparser.StarExpr{},
					},
				},
			},
		},
		From: table,
	}

	// 生成 SQL
	buf := sqlparser.NewTrackedBuffer(nil)
	selectStmt.Format(buf)

	return buf.String(), nil
}

// CheckUpdateContainField
func CheckUpdateContainField(sql string, key string, value string) (bool, error) {
	// 解析 SQL
	sql = StripLeadingComments(sql)
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return false, err
	}

	// 如果解析结果不是 UPDATE 语句，返回错误
	updateStmt, ok := stmt.(*sqlparser.Update)
	if !ok {
		return true, nil
	}

	exprs := updateStmt.Exprs
	for _, expr := range exprs {
		if expr.Name.Name.String() == key {
			if value != "" {
				buf := sqlparser.NewTrackedBuffer(nil)
				expr.Expr.Format(buf)
				// fmt.Println("value = ", buf.String()) // debug
				if buf.String() == value {
					return true, nil
				} else {
					return false, errors.Errorf("字段 '%s' 的值不符合要求, 期待: %s, 实际上: %s", key, value, buf.String())
				}
			}
			return true, nil
		}
	}
	// TODO 判断字段是否为空
	return false, errors.Errorf("备注字段 '%s' 未更新", key)
}

// Deprecated: warpSqlChinese 给中文增加”
func warpSqlChinese(sql string) string {
	res := make([]rune, 0, len(sql))
	hasPrefix := 0
	for _, r := range sql {

		if unicode.Is(unicode.Scripts["Han"], r) {
			if hasPrefix == 0 {
				res = append(res, '\'')
				//if i-1 > 0 && sql[i-1] != '\'' {
				//	res = append(res, '\'')
				//}
				hasPrefix ^= 1
			}
			//res = append(res, rune(r))
		} else if r == '\'' {
			hasPrefix ^= 1
		} else {
			if hasPrefix == 1 && (r == ',' || r == '\t' || r == ' ' || r == '\n') {
				//if i-1 > 0 && sql[i-1] != '\'' {
				//	res = append(res, '\'')
				//}
				res = append(res, '\'')
				hasPrefix ^= 1
			}

		}
		res = append(res, rune(r))
		//fmt.Println(hasPrefix, string(r), string(res))
	}

	resSql := string(res)
	//return strings.ReplaceAll(resSql, "''", "'")
	return resSql
}

func warpSqlChinese2(sql string) string {
	sql = strings.ReplaceAll(sql, ",", ", ")
	strs := strings.Fields(strings.TrimSpace(sql))
	res := ""
	for _, str := range strs {
		needQuote := false

		for _, r := range str {
			if r == '\'' {
				needQuote = false
			}
			if unicode.Is(unicode.Scripts["Han"], r) {
				needQuote = true
			}
		}

		if needQuote {
			if strings.HasSuffix(str, ",") {
				str = strings.TrimSuffix(str, ",")
				res += " '" + str + "'" + ","
			} else {
				res += " '" + str + "'"
			}
		} else {
			res += " " + str
		}
		res += "\n"
	}
	return res
}

// GetQueryFields 获得查询语句要查的字段
func GetQueryFields(sql string) ([]string, error) {
	// 解析SQL语句，返回抽象语法树
	sql = warpSqlChinese2(sql)
	sql = StripLeadingComments(sql)
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		//return nil, errors.New(fmt.Sprintf("%v, 可尝试将查询字段由< name 姓名 >改成< name as '姓名'>", err))
		return []string{}, nil
	}

	// 将语法树强制转换为SELECT语句
	selectStmt, ok := stmt.(*sqlparser.Select)
	if !ok {
		//return nil, errors.New("Not a select statement")
		return []string{}, nil
	}

	fields := make([]string, 0)
	// 遍历SELECT语句中的列，并打印它们的名称
	for _, expr := range selectStmt.SelectExprs {
		switch expr := expr.(type) {
		case *sqlparser.StarExpr:
			fields = append(fields, "*")
		case *sqlparser.AliasedExpr:
			if !expr.As.IsEmpty() {
				field := sqlparser.String(expr.As)
				field = strings.Trim(field, "`")
				fields = append(fields, field)
			} else {
				field := sqlparser.String(expr.Expr)
				if strings.Contains(field, ".") {
					field = field[strings.LastIndex(field, ".")+1:]
				}
				fields = append(fields, field)
			}
		}
	}

	return fields, nil
}

// EscapeSql 转义 sql 语句
// \ -> \\
// ' -> \'
func EscapeSql(sql string) string {
	sql = strings.ReplaceAll(sql, "\\", "\\\\")
	sql = strings.ReplaceAll(sql, "'", "\\'")
	return sql
}

// GetTablesFromSelectStatement 解析 SELECT 查询语句，获取相关的表名
func GetTablesFromSelectStatement(sql string) ([]string, error) {
	// 解析查询语句
	sql = warpSqlChinese2(sql)
	sql = StripLeadingComments(sql)
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return nil, errors.Errorf("sqlparser.Parse 失败, err: %v", err)
	}

	// 判断是否是 SELECT 查询语句
	sel, ok := stmt.(*sqlparser.Select)
	if !ok {
		return nil, errors.Errorf("%s is not an select statement", sql)
	}

	// 获取 SELECT 查询语句中的表名
	tables := make([]string, 0)
	for _, expr1 := range sel.From {
		ts, err := parseTableExpr(expr1)
		if err != nil {
			return nil, err
		}
		tables = append(tables, ts...)
		//switch t := expr1.(type) {
		//case *sqlparser.AliasedTableExpr:
		//	fmt.Println(1)
		//	expr := sqlparser.String(t.Expr)
		//	if strings.Contains(expr, "(") {
		//		fmt.Println(expr)
		//		expr = strings.TrimPrefix(expr, "(")
		//		expr = strings.TrimSuffix(expr, ")")
		//		subTables, err := GetTablesFromSelectStatement(expr)
		//		if err != nil {
		//			return nil, err
		//		}
		//		tables = append(tables, subTables...)
		//	} else {
		//		tables = append(tables, expr)
		//	}
		//case *sqlparser.JoinTableExpr:
		//	fmt.Println(2)
		//	leftExpr := sqlparser.String(t.LeftExpr)
		//	rightExpr := sqlparser.String(t.RightExpr)
		//	tables = append(tables, leftExpr)
		//	tables = append(tables, rightExpr)
		//case *sqlparser.ParenTableExpr:
		//	fmt.Println(3)
		//	subTables, err := GetTablesFromSelectStatement(sqlparser.String(t))
		//	if err != nil {
		//		return nil, err
		//	}
		//	tables = append(tables, subTables...)
		//default:
		//	return nil, fmt.Errorf("unsupported table expression: %T", t)
		//}
	}

	return tables, nil
}

func parseTableExpr(expr1 sqlparser.TableExpr) ([]string, error) {
	tables := make([]string, 0)

	switch t := expr1.(type) {
	case *sqlparser.AliasedTableExpr:
		//fmt.Println(1)
		expr := sqlparser.String(t.Expr)
		switch t.Expr.(type) {
		case *sqlparser.Subquery:
			//fmt.Println(1.1)
			//fmt.Println(expr)
			expr = strings.TrimPrefix(expr, "(")
			expr = strings.TrimSuffix(expr, ")")
			subTables, err := GetTablesFromSelectStatement(expr)
			if err != nil {
				return nil, err
			}
			tables = append(tables, subTables...)
		default:
			//fmt.Println(1.2)
			tables = append(tables, expr)
		}
		//expr := sqlparser.String(t.Expr)
		//if strings.Contains(expr, "(") {
		//	fmt.Println(expr)
		//	expr = strings.TrimPrefix(expr, "(")
		//	expr = strings.TrimSuffix(expr, ")")
		//	subTables, err := GetTablesFromSelectStatement(expr)
		//	if err != nil {
		//		return nil, err
		//	}
		//	tables = append(tables, subTables...)
		//} else {
		//	tables = append(tables, expr)
		//}
	case *sqlparser.JoinTableExpr:
		//fmt.Println(2)

		//leftExpr := sqlparser.String(t.LeftExpr)
		//rightExpr := sqlparser.String(t.RightExpr)
		//fmt.Println(leftExpr, rightExpr)
		leftTables, err := parseTableExpr(t.LeftExpr)
		if err != nil {
			return nil, err
		}
		rightTables, err := parseTableExpr(t.RightExpr)
		if err != nil {
			return nil, err
		}
		tables = append(tables, leftTables...)
		tables = append(tables, rightTables...)
	case *sqlparser.ParenTableExpr:
		//fmt.Println(3)
		subTables, err := GetTablesFromSelectStatement(sqlparser.String(t))
		if err != nil {
			return nil, err
		}
		tables = append(tables, subTables...)
	default:
		return nil, fmt.Errorf("unsupported table expression: %T", t)
	}

	return tables, nil
}

type SelectTableInfo struct {
	DbCode    string
	TableCode string
}

// GetSelectTableInfosFromSelectStatement 获得查询语句涉及到的库表
func GetSelectTableInfosFromSelectStatement(query string) ([]SelectTableInfo, error) {
	tableInfos := make([]SelectTableInfo, 0)
	tables, err := GetTablesFromSelectStatement(query)
	if err != nil {
		return tableInfos, nil
		//return nil, err
	}

	//tableInfos := make([]SelectTableInfo, 0, len(tables))
	for _, table := range tables {
		dbCode, tableCode, err := ParseTableName(table)
		if err != nil {
			continue
		}
		tableInfos = append(tableInfos, SelectTableInfo{
			DbCode:    dbCode,
			TableCode: tableCode,
		})
	}

	return tableInfos, nil
}

const dns1035LabelFmt string = "[a-z]([_a-z0-9]*[a-z0-9])?"

var dns1035LabelRegexp = regexp.MustCompile("^(?i)" + dns1035LabelFmt + "$")

func checkTableName(value string) bool {
	return dns1035LabelRegexp.MatchString(value)
}

func ParseTableName(table string) (string, string, error) {
	tables := strings.Split(table, " as")
	if len(tables) == 0 {
		return "", "", errors.New("table name is empty")
	}

	table = tables[0]
	if strings.Contains(table, ".") {
		ts := strings.Split(table, ".")
		if len(ts) == 2 {
			if !checkTableName(ts[0]) || !checkTableName(ts[1]) {
				return "", "", errors.New("table name is invalid")
			}
			return ts[0], ts[1], nil
		} else {
			return "", "", errors.New("table name is invalid")
		}
	} else {
		if !checkTableName(table) {
			return "", "", errors.New("table name is invalid")
		}
		return "", table, nil
	}
}

// CheckInsertContainField 检查 sql 是否为 insert 语句, 并判断是否 insert 了指定字段
func CheckInsertContainField(sql string, key string) (bool, error) {
	sql = StripLeadingComments(sql)
	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return false, err
	}

	// 如果解析结果不是 INSERT 语句，返回
	insertStmt, ok := stmt.(*sqlparser.Insert)
	if !ok {
		return true, nil
	}

	columns := insertStmt.Columns
	rows := insertStmt.Rows
	keyColumn := sqlparser.NewColIdent(key)
	idx := columns.FindColumn(keyColumn)
	if idx == -1 {
		return false, errors.Errorf("insert 语句未包含 %s 字段", key)
	}

	flag := true

	// 可能是 sqlparser.Select
	values, ok := rows.(sqlparser.Values)
	if !ok {
		return true, nil
	}

	for _, row := range values {
		for i, val := range row {
			if i != idx {
				continue
			}

			if len(sqlparser.String(val)) == 2 {
				flag = false
				break
			}
		}
		if !flag {
			break
		}
	}

	if flag {
		return true, nil
	}
	return false, errors.Errorf("insert 语句 %s 字段值为空", key)
}

var (
	reLimitN           = regexp.MustCompile("(?i)limit\\s+(\\d+)\\s*$")
	reLimitOffset      = regexp.MustCompile("(?i)limit\\s+(\\d+)\\s+offset\\s+(\\d+)\\s*$")
	reOffsetCommaLimit = regexp.MustCompile("(?i)limit\\s+(\\d+)\\s*,\\s*(\\d+)\\s*$")

	reOnlySelect     = regexp.MustCompile("(?i)^select")
	reSelectNow      = regexp.MustCompile("(?i)^select\\snow()")
	reSelectVersion  = regexp.MustCompile("(?i)^select\\sversion()")
	reSpecialSelects = []*regexp.Regexp{reSelectNow, reSelectVersion}
)

func SqlAddLimit(sqlContent string, limit int) string {
	sqlContent = StripLeadingComments(sqlContent)
	sql := strings.TrimSpace(strings.TrimRight(strings.TrimSpace(sqlContent), ";"))

	// 特定语句不添加 limit
	for _, re := range reSpecialSelects {
		if re.MatchString(sql) {
			return sql + ";"
		}
	}

	if reOnlySelect.MatchString(sql) {
		// LIMIT N
		limitN := reLimitN.FindStringSubmatch(sql)
		// LIMIT M OFFSET N
		limitOffsetN := reLimitOffset.FindStringSubmatch(sql)
		// LIMIT M,N
		offsetCommaLimitN := reOffsetCommaLimit.FindStringSubmatch(sql)

		if len(limitN) > 1 {
			sqlLimit := convert.ToInt(limitN[1])
			l := math.Min(float64(limit), float64(sqlLimit))
			sql = reLimitN.ReplaceAllString(sql, fmt.Sprintf("limit %d", int64(l)))
		} else if len(limitOffsetN) > 2 {
			sqlLimit := convert.ToInt(limitOffsetN[1])
			sqlOffset := convert.ToInt(limitOffsetN[2])
			l := math.Min(float64(sqlLimit), float64(limit))
			sql = reLimitOffset.ReplaceAllString(sql, fmt.Sprintf("limit %d offset %d", int64(l), sqlOffset))
		} else if len(offsetCommaLimitN) > 2 {
			sqlOffset := convert.ToInt(offsetCommaLimitN[1])
			sqlLimit := convert.ToInt(offsetCommaLimitN[2])
			l := math.Min(float64(sqlLimit), float64(limit))
			sql = reOffsetCommaLimit.ReplaceAllString(sql, fmt.Sprintf("limit %d,%d", sqlOffset, convert.ToInt64(l)))
		} else {
			sql = fmt.Sprintf("%s limit %d", sql, limit)
		}
	}

	sql = sql + ";"
	return sql
}

// CheckSqlAlterAndGetTableCode 检查语句是否是 alter 语句, 并返回修改的表名
func CheckSqlAlterAndGetTableCode(sql string) string {
	sql = StripLeadingComments(sql)

	if !reAlter.MatchString(sql) {
		return ""
	}

	stmt, err := sqlparser.Parse(sql)
	if err != nil {
		return ""
	}

	// 如果解析结果不是 ddlStmt 语句，返回错误
	ddlStmt, ok := stmt.(*sqlparser.DDL)
	if !ok {
		return ""
	}

	if ddlStmt.Action != sqlparser.AlterStr {
		return ""
	}

	return ddlStmt.Table.Name.String()
}
