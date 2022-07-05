package tree_test

import (
	"io/fs"
	"reflect"
	"regexp"
	"testing"
	"testing/fstest"

	"github.com/yahuian/marker/pkg/tree"
)

func TestTree(t *testing.T) {
	fsys := fstest.MapFS{
		".vscode/settings.json": &fstest.MapFile{},
		".vscode/xyz.png":       &fstest.MapFile{},
		"README.md":             &fstest.MapFile{Data: []byte("![](./a.png)\n![](images/a.png)")},
		"a.png":                 &fstest.MapFile{},
		"b.png":                 &fstest.MapFile{},
		"c.png":                 &fstest.MapFile{},
		"note.md":               &fstest.MapFile{Data: []byte(`![](images/dir/x.jpg)`)},
		"images/a.png":          &fstest.MapFile{},
		"images/b.png":          &fstest.MapFile{},
		"images/readme.md":      &fstest.MapFile{Data: []byte(`![](../b.png)`)},
		"images/dir/x.jpg":      &fstest.MapFile{},
		"images/dir/y.jpg":      &fstest.MapFile{},
		"images/dir/test.md":    &fstest.MapFile{Data: []byte(`![](../../b.png)`)},
	}

	root, err := tree.NewTree(fsys, func(d fs.DirEntry) bool {
		return regexp.MustCompile(`^\.`).MatchString(d.Name())
	})
	if err != nil {
		t.Fatal(err)
	}
	images, err := root.GetUselessImages(fsys, []string{".png", ".jpg"})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(images, []string{"c.png", "images/b.png", "images/dir/y.jpg"}) {
		t.Fatal("get useless images err")
	}

	// for debug
	// data, err := json.Marshal(root)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Error(string(data))
}
