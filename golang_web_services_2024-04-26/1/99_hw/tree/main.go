package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err)
	}
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	err := dirTreeRecursive(out, path, printFiles, "")
	if err != nil {
		return err
	}
	return nil
}

func dirTreeRecursive(out io.Writer, path string, printFiles bool, prefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	var filteredFiles []os.DirEntry

	for _, entries := range files {
		if printFiles {
			filteredFiles = append(filteredFiles, entries)
		} else {
			if entries.IsDir() {
				filteredFiles = append(filteredFiles, entries)
			}
		}
	}
	for i, filteredEntries := range filteredFiles {
		ending := ""
		fileInfo, err := filteredEntries.Info()
		if err != nil {
			return err
		}
		if fileInfo.Mode().IsRegular() {
			if fileInfo.Size() == 0 {
				ending = " (empty)"
			} else {
				ending = fmt.Sprintf(" (%db)", fileInfo.Size())
			}
		}
		if i == len(filteredFiles)-1 {
			_, err := fmt.Fprintf(out, "%s└───%v%s\n", prefix, filteredEntries.Name(), ending)
			if err != nil {
				return err
			}
		} else {
			_, err := fmt.Fprintf(out, "%s├───%v%s\n", prefix, filteredEntries.Name(), ending)
			if err != nil {
				return err
			}
		}
		if filteredEntries.IsDir() {
			nextPrefix := prefix
			if i == len(filteredFiles)-1 {
				nextPrefix += "\t"
			} else {
				nextPrefix += "│\t"
			}
			nextPath := filepath.Join(path, filteredEntries.Name())
			err := dirTreeRecursive(out, nextPath, printFiles, nextPrefix)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
