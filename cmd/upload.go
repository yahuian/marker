/*
Copyright © 2022 yahuian <yahuian@126.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/kolo/xmlrpc"
	"github.com/spf13/cobra"
	"github.com/yahuian/marker/config"
	"github.com/yahuian/marker/pkg/metaweblog"
	"github.com/yahuian/marker/pkg/tree"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload local images to blog platform and generate a new markdown file.",
	RunE:  runUpload,
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func runUpload(cmd *cobra.Command, args []string) error {
	if len(config.Val.BlogPlatforms) == 0 {
		return fmt.Errorf("blog platforms config not found")
	}

	fsys := os.DirFS(root)
	t, err := tree.NewTree(fsys, func(d fs.DirEntry) bool {
		return config.SkipFiles(d)
	})
	if err != nil {
		return fmt.Errorf("new tree err: %w", err)
	}

	for _, v := range config.Val.BlogPlatforms {
		client, err := metaweblog.NewClient(v.API)
		if err != nil {
			return fmt.Errorf("new client err: %w", err)
		}
		defer client.Close()

		upload := func(name, b64 string) (string, error) {
			file := metaweblog.FileData{Bits: xmlrpc.Base64(b64), Name: name, Type: ""}
			url, err := client.NewMediaObject(v.BlogID, v.Username, v.AppKey, file)
			if err != nil {
				return "", err
			}
			return url, nil
		}

		if err := t.UploadImage(root, fsys, upload, v.Kind); err != nil {
			return fmt.Errorf("upload image err: %w", err)
		}
	}

	fmt.Println("finished upload.")
	return nil
}
