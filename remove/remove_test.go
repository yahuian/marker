package remove

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func Test_parseImages(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "one image",
			args: args{
				text: `![a](image/a.txt)`,
			},
			want: "image/a.txt",
		},
		{
			name: "no image",
			args: args{
				text: `this is title`,
			},
			want: "",
		},
		{
			name: "online",
			args: args{
				text: `![online-image](https://www.abc.com/a.png)`,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseImagePath(tt.args.text); got != tt.want {
				t.Errorf("parseImages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getUselessImages(t *testing.T) {
	type args struct {
		fsys fs.FS
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
					"README.md": &fstest.MapFile{Data: []byte(`
I am README file.
![image-a](image/a.png)
![online-image](https://www.abc.com/a.png)
					`)},
					"image/a.png": &fstest.MapFile{},
					"image/b.jpg": &fstest.MapFile{},
				},
			},
			want: []string{
				"image/b.jpg",
			},
			wantErr: false,
		},
		{
			name: "many markdown files",
			args: args{
				fsys: fstest.MapFS{
					"blog/README.md":   &fstest.MapFile{Data: []byte(`![image-a](image/a.png)`)},
					"blog/note.md":     &fstest.MapFile{Data: []byte(`![image-b](image/b.jpg)`)},
					"blog/image/a.png": &fstest.MapFile{},
					"blog/image/b.jpg": &fstest.MapFile{},
					"blog/image/c.jpg": &fstest.MapFile{},
				},
			},
			want: []string{
				"blog/image/c.jpg",
			},
			wantErr: false,
		},
		{
			name: "nested",
			args: args{
				fsys: fstest.MapFS{
					"a.md":               &fstest.MapFile{Data: []byte(`![image-a](image/a.png)`)},
					"image/b.png":        &fstest.MapFile{},
					"blog/README.md":     &fstest.MapFile{Data: []byte(`![image](image/a.png)`)},
					"blog/note.md":       &fstest.MapFile{Data: []byte(`![image](image/b.jpg)`)},
					"blog/image/a.png":   &fstest.MapFile{},
					"blog/image/xyz.jpg": &fstest.MapFile{},
					"blog/image/b.jpg":   &fstest.MapFile{},
				},
			},
			want: []string{
				"image/b.png",
				"blog/image/xyz.jpg",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getUselessImages(tt.args.fsys)
			if (err != nil) != tt.wantErr {
				t.Errorf("getUselessImages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getUselessImages() = %v, want %v", got, tt.want)
			}
		})
	}
}
