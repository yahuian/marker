package utils

import (
	"io/fs"
	"path"
	"reflect"
	"regexp"
	"testing"
	"testing/fstest"
)

func TestGetAllFiles(t *testing.T) {
	type args struct {
		fsys fs.FS
		skip func(d fs.DirEntry) bool
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				fsys: fstest.MapFS{
					"1.txt":     &fstest.MapFile{},
					"dir/a.txt": &fstest.MapFile{},
					"dir/b.txt": &fstest.MapFile{},
				},
				skip: nil,
			},
			want:    []string{"1.txt", "dir/a.txt", "dir/b.txt"},
			wantErr: false,
		},
		{
			name: "skip . file",
			args: args{
				fsys: fstest.MapFS{
					".idea":     &fstest.MapFile{},
					".git":      &fstest.MapFile{},
					"dir/a.txt": &fstest.MapFile{},
					"dir/b.txt": &fstest.MapFile{},
				},
				skip: func(d fs.DirEntry) bool {
					return regexp.MustCompile(`^\.`).MatchString(d.Name())
				},
			},
			want:    []string{"dir/a.txt", "dir/b.txt"},
			wantErr: false,
		},
		{
			name: "only markdown file",
			args: args{
				fsys: fstest.MapFS{
					".idea":     &fstest.MapFile{},
					".git":      &fstest.MapFile{},
					"dir/a.md":  &fstest.MapFile{},
					"dir/b.txt": &fstest.MapFile{},
				},
				skip: func(d fs.DirEntry) bool {
					return !d.IsDir() && path.Ext(d.Name()) != ".md"
				},
			},
			want:    []string{"dir/a.md"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllFiles(tt.args.fsys, tt.args.skip)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}
