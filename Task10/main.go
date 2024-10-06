package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	startPath := "./"
	err := directoryTree(startPath, 0)
	if err != nil {
		fmt.Println(err)
	}
}

// directoryTree - display tree of all directories in the path
func directoryTree(startPath string, level int) error {
	dir, err := os.ReadDir(startPath)
	if err != nil {
		return err
	}
	for _, file := range dir {
		thisPath := strings.Repeat(".    ", level)
		fmt.Printf("%s|- %s\n", thisPath, file.Name())

		if file.IsDir() {
			newPath := filepath.Join(startPath, file.Name())
			err := directoryTree(newPath, level+1)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
