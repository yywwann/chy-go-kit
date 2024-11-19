package sqltool

import (
	"github.com/pkg/errors"
	"github.com/xwb1989/sqlparser"
)

type SqlAnalyseResult struct {
	AlterCount       int `json:"alterCount,omitempty"`
	AlterAddCount    int `json:"alterAddCount,omitempty"`
	AlterModifyCount int `json:"alterModifyCount,omitempty"`
	AlterDropCount   int `json:"alterDropCount,omitempty"`
	InsertCount      int `json:"insertCount,omitempty"`
	UpdateCount      int `json:"updateCount,omitempty"`
	DeleteCount      int `json:"deleteCount,omitempty"`
	CreateTableCount int `json:"createTableCount,omitempty"`
	DropTableCount   int `json:"dropTableCount,omitempty"`
	SelectCount      int `json:"selectCount,omitempty"`
	UnknownCount     int `json:"unknownCount,omitempty"`
	SetCount         int `json:"setCount,omitempty"`
	UseCount         int `json:"useCount,omitempty"`
	DBDDLCount       int `json:"dBDDLCount,omitempty"`
	BeginCount       int `json:"beginCount,omitempty"`
	CommitCount      int `json:"commitCount,omitempty"`
	RollbackCount    int `json:"rollbackCount,omitempty"`
	ShowCount        int `json:"showCount,omitempty"`
}

func (a *SqlAnalyseResult) MergeResult(partResult *SqlAnalyseResult) {
	a.AlterCount += partResult.AlterCount
	a.AlterAddCount += partResult.AlterAddCount
	a.AlterModifyCount += partResult.AlterModifyCount
	a.AlterDropCount += partResult.AlterDropCount
	a.InsertCount += partResult.InsertCount
	a.UpdateCount += partResult.UpdateCount
	a.DeleteCount += partResult.DeleteCount
	a.CreateTableCount += partResult.CreateTableCount
	a.DropTableCount += partResult.DropTableCount
	a.SelectCount += partResult.SelectCount
	a.UnknownCount += partResult.UnknownCount
	a.SetCount += partResult.SetCount
	a.UseCount += partResult.UseCount
	a.DBDDLCount += partResult.DBDDLCount
	a.BeginCount += partResult.BeginCount
	a.CommitCount += partResult.CommitCount
	a.RollbackCount += partResult.RollbackCount
	a.ShowCount += partResult.ShowCount
}

func SqlAnalyseWithFile(filename string) (*SqlAnalyseResult, error) {
	sqls, err := SplitSQLWithFile(filename)
	if err != nil {
		return nil, err
	}

	return SqlAnalyse(sqls)
}

func SqlAnalyse(sqls []string) (*SqlAnalyseResult, error) {
	resp := &SqlAnalyseResult{}

	for _, sql := range sqls {
		sql = StripLeadingComments(sql)
		stmt, err := sqlparser.Parse(sql)
		if err != nil {
			return nil, errors.Errorf("解析 sql 失败, err: %v, sql: %s", err, sql)
		}

		switch stmt := stmt.(type) {
		case *sqlparser.DDL:
			switch stmt.Action {
			case sqlparser.CreateStr:
				resp.CreateTableCount += 1
			case sqlparser.DropStr:
				resp.DropTableCount += 1
			case sqlparser.AlterStr:
				resp.AlterCount += 1
			}
		case *sqlparser.Insert:
			resp.InsertCount += 1
		case *sqlparser.Update:
			resp.UpdateCount += 1
		case *sqlparser.Delete:
			resp.DeleteCount += 1
		case *sqlparser.Select:
			resp.SelectCount += 1
		case *sqlparser.Show:
			resp.ShowCount += 1
		case *sqlparser.Set:
			resp.SetCount += 1
		case *sqlparser.Use:
			resp.UseCount += 1
		case *sqlparser.DBDDL:
			resp.DBDDLCount += 1
		case *sqlparser.Begin:
			resp.BeginCount += 1
		case *sqlparser.Commit:
			resp.CommitCount += 1
		case *sqlparser.Rollback:
			resp.RollbackCount += 1
		default:
			resp.UnknownCount += 1
		}

	}

	return resp, nil
}
