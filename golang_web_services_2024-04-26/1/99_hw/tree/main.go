package main

import (
	"os"
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
		panic("")
	}
}

func dirTree(path string, printFiles bool) (out *os.File) {

	printFiles = false

	files, err := os.ReadDir(path)
	if err != nil {
		panic("Bad dir")
	}

	for _, entries := range files {
		if entries.IsDir() {
			printFiles = false
			if printFiles == false {
				dirTree(path, printFiles)
			}
		} else {
			printFiles = true

		}
	}
	return out
}
