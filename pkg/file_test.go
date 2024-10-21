package pkg

import (
	"testing"
)

func Test_pathExists(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathExists(tt.args.path); got != tt.want {
				t.Errorf("pathExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateFile(t *testing.T) {
	type args struct {
		filePath string
		fileName string
	}
	tests := []struct {
		name string
		args args
		// want *os.File
	}{
		// TODO: Add test cases.
		{"first", args{"uploads", "666.log"}},
		{"first", args{"uploads", ""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CreateFile(tt.args.filePath, tt.args.fileName)
			t.Logf("got:%v\n", got)
		})
	}
}
