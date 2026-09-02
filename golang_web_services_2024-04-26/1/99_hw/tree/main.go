package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {

	dirTree(os.Stdout, "/Users/mattew/Desktop/testFolder", false)
	/*
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

	*/
}

func dirTree(out *os.File, path string, printFiles bool) {

	files, err := os.ReadDir(path)
	if err != nil {
		panic("Bad dir")
	}
	for i, entries := range files {
		if entries.IsDir() {
			printFiles = false
			fmt.Printf("└───"+"%v \n\t", entries)
			nextPath := filepath.Join(path, entries.Name())
			dirTree(out, nextPath, printFiles)
		} else {
			printFiles = true
			if i == len(files)-1 {
				fmt.Printf("└───"+"%v", entries)
			} else {
				fmt.Printf("├───%v \n", entries)
			}
		}
	}
}

/*
if i == len(files)-1 {
fmt.Printf("└───"+"%v", entries)
}*/
