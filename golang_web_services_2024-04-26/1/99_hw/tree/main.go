package main

import (
	"fmt"
	"io/fs"
	"os"
)

func main() {

	err, _ := dirTree()
	if err != nil {
		return
	}

	/*out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
	*/
}

type DirEntry = fs.DirEntry

func dirTree(name string) ([]DirEntry, error) {

	files, err := os.ReadDir(name)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return nil, err
}
