package utils

import (
	"io/fs"
)

// GetAllFiles return file path in fsys, like test/a/a.txt
func GetAllFiles(fsys fs.FS, skip func(d fs.DirEntry) bool) ([]string, error) {
	var paths []string

	err := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// . is current dir
		if d.Name() == "." {
			return nil
		}

		// skip files
		if skip != nil && skip(d) {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			paths = append(paths, p)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}
