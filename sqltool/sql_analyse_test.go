package sqltool

import (
	"reflect"
	"testing"
)

func TestSqlAnalyseWithFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *SqlAnalyseResult
		wantErr bool
	}{
		{
			name:    "",
			args:    args{filename: "/Users/xiniu/Downloads/tmp_2024-07-10/total.sql"},
			want:    &SqlAnalyseResult{},
			wantErr: false,
		},
		{
			name:    "",
			args:    args{filename: "/Users/xiniu/Downloads/ops_op_log.sql"},
			want:    &SqlAnalyseResult{},
			wantErr: false,
		},
		{
			name: "",
			args: args{filename: "/Users/xiniu/Downloads/matris-test.sql"},

			want:    &SqlAnalyseResult{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SqlAnalyseWithFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("SqlAnalyseWithFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SqlAnalyseWithFile() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
