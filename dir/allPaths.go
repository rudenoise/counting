package dir

import (
	"os"
	"path/filepath"
	"regexp"
)

func AllPaths(path, exclude, include string) ([]string, error) {
	paths := make([]string, 0)

	err := filepath.Walk(path, func(cPath string, f os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		excludeMatch, err := regexp.MatchString(exclude, cPath)
		if err != nil {
			panic(err)
		}
		includeMatch, err := regexp.MatchString(include, cPath)
		if err != nil {
			panic(err)
		}

		if f.IsDir() == false && excludeMatch == false && includeMatch {
			paths = append(paths, cPath)
		}
		return err
	})
	if err != nil {
		return paths, err
	}
	return paths, nil
}
