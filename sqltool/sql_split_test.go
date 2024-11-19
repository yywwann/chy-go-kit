package sqltool

import (
	"reflect"
	"testing"
)

func TestSplitSQLUnsafe(t *testing.T) {
	type args struct {
		sql string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 string
	}{
		{
			name: "1",
			args: args{
				sql: "-- 注释\nalter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';\nalter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';",
			},
			want: []string{
				"-- 注释\nalter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';",
				"alter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';",
			},
			want1: "",
		},
		{
			name: "2",
			args: args{
				sql: "\n-- 注释\nalter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';\nalter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 ",
			},
			want: []string{
				"-- 注释\nalter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 0:关闭 1:开启';",
			},
			want1: "\nalter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT '表单修改权限 ",
		},
		{
			name: "3",
			args: args{
				sql: "''",
			},
			want:  nil,
			want1: "''",
		},
		{
			name: "3",
			args: args{
				sql: "alter table qqq",
			},
			want:  nil,
			want1: "alter table qqq",
		},
		{
			name: "4",
			args: args{
				sql: `alter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT  '奇怪的注释\\\'';`,
			},
			want: []string{
				`alter table emr_template_config add COLUMN ep_flag tinyint(1) DEFAULT 1 COMMENT  '奇怪的注释\\\'';`,
			},
			want1: ``,
		},
		{
			name: "",
			args: args{
				sql: `/** 请在晚上空闲时执行 **/
alter table dtc_doctor_order
modify extension mediumtext null comment '扩展信息' ,ALGORITHM = INPLACE,LOCK = NONE;`,
			},
			want:  []string{`/** 请在晚上空闲时执行 **/
alter table dtc_doctor_order
modify extension mediumtext null comment '扩展信息' ,ALGORITHM = INPLACE,LOCK = NONE;`},
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := SplitSQLUnsafe(tt.args.sql)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitSQLUnsafe() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SplitSQLUnsafe() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
