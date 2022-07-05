package utils_test

import (
	"testing"

	"github.com/yahuian/marker/utils"
)

func TestParseMarkdownImage(t *testing.T) {
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
			if got := utils.ParseMarkdownImage(tt.args.text); got != tt.want {
				t.Errorf("parseImages() = %v, want %v", got, tt.want)
			}
		})
	}
}
