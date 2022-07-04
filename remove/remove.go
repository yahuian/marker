package remove

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yahuian/marker/config"
	"github.com/yahuian/marker/utils"
)

var (
	Yes  = false
	Root = "./"
)

func RunE(cmd *cobra.Command, args []string) error {
	images, err := getUselessImages(os.DirFS(Root))
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

func getUselessImages(fsys fs.FS) ([]string, error) {
	// get all markdown files
	list, err := utils.GetAllFiles(fsys, func(d fs.DirEntry) bool {
		return config.SkipFiles(d) || (!d.IsDir() && path.Ext(d.Name()) != ".md")
	})
	if err != nil {
		return nil, fmt.Errorf("get all markdown files err: %w", err)
	}

	/*
		maybe there are many markdown files in one dir
		group markdown files with prefix dir
			a.md
			b.md
			image/a.png
	*/
	filesMap := make(map[string][]string)
	for _, v := range list {
		k := path.Dir(v)
		filesMap[k] = append(filesMap[k], v)
	}

	var results []string

	for _, files := range filesMap {
		used := make(map[string]struct{})
		for _, v := range files {
			// get used images path in markdown file
			list, err := getUsedImages(fsys, v)
			if err != nil {
				return nil, fmt.Errorf("get used images err: %w", err)
			}
			if len(list) == 0 {
				continue
			}
			for _, v := range list {
				used[v] = struct{}{}
			}
		}

		if len(used) == 0 {
			continue
		}

		/*
			get all images in relative path
			typical relative file struct for markdown:
			image/
				a.png
				b.png
			README.md
		*/
		var dir string
		for k := range used {
			dir = path.Dir(k)
			break // we think there is only one image dir
		}
		subFS, err := fs.Sub(fsys, dir) // image/ sub fs
		if err != nil {
			return nil, fmt.Errorf("get sub fs err: %w", err)
		}
		images, err := utils.GetAllFiles(subFS, config.SkipFiles) // a.png b.png
		if err != nil {
			return nil, fmt.Errorf("get image files err: %w", err)
		}

		for _, image := range images {
			p := path.Join(dir, image) // image/a.png
			if _, ok := used[p]; ok {
				continue
			}
			results = append(results, p)
		}
	}

	return results, nil
}

func getUsedImages(fsys fs.FS, filePath string) ([]string, error) {
	f, err := fsys.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file err: %w", err)
	}
	defer f.Close()

	var images []string

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		image := parseImagePath(scan.Text())
		if image != "" {
			// filepath: blog/README.md
			// image: image/a.png
			// result is blog/image/a.png
			images = append(images, path.Join(path.Dir(filePath), image))
		}
	}

	return images, nil
}

// BUG one line maybe have many images
var imageRegex = regexp.MustCompile(`!\[.*\]\((.*)\)`)

func parseImagePath(text string) string {
	list := imageRegex.FindAllStringSubmatch(text, 1)
	for _, v := range list {
		if len(v) != 2 {
			continue
		}
		image := v[1]
		// skip online image
		if strings.HasPrefix(image, "http://") || strings.HasPrefix(image, "https://") {
			continue
		}
		return image
	}
	return ""
}
