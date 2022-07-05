package remove

import (
	"fmt"
	"io/fs"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/yahuian/marker/config"
	"github.com/yahuian/marker/pkg/tree"
)

var (
	Yes  = false
	Root = "./"
)

func RunE(cmd *cobra.Command, args []string) error {
	fsys := os.DirFS(Root)
	t, err := tree.NewTree(fsys, func(d fs.DirEntry) bool {
		return config.SkipFiles(d)
	})
	if err != nil {
		return fmt.Errorf("new tree err: %w", err)
	}
	images, err := t.GetUselessImages(fsys, config.Val.ImageTypes)
	if err != nil {
		return fmt.Errorf("get useless images err: %w", err)
	}

	if len(images) == 0 {
		fmt.Println("Well done, your images are all used.")
		return nil
	}

	if !Yes {
		fmt.Println("These images are useless, you can remove them with --yes flag.")
	}
	for _, v := range images {
		if !Yes {
			fmt.Println(v)
			continue
		}

		p := path.Join(Root, v)
		if err := os.Remove(p); err != nil {
			return fmt.Errorf("remove %s err: %w", p, err)
		}

		fmt.Printf("[removed] %s\n", p)
	}

	return nil
}
