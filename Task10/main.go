package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var startPath = flag.String("startPath", "./Task10", "path to read")

func main() {
	flag.Parse()
	result, err := directoryTree(*startPath, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}

// directoryTree - display tree of all directories in the path
func directoryTree(path string, level int) ([]string, error) {
	dir, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, file := range dir {
		thisPath := strings.Repeat(".    ", level)
		result = append(result, thisPath+file.Name())
		if file.IsDir() {
			newPath := filepath.Join(path, file.Name())
			subDirRes, err := directoryTree(newPath, level+1)
			if err != nil {
				return nil, err
			}

			result = append(result, subDirRes...)
		}
	}

	return result, nil
}
