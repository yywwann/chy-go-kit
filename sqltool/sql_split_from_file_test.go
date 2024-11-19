package sqltool

import (
	"testing"
)

func TestSplitSqlFromFileAndSaveFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				filename: "/Users/xiniu/Downloads/matris-test.sql",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := splitSqlFromFileAndSaveFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitSqlFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
