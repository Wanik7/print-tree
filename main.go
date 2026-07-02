package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func printObj(name string, fileType bool, nest int) {

	const wide = 2

	values := map[bool]string{
		true:  "dir",
		false: "file",
	}

	spaces := strings.Repeat(" ", nest*wide)

	fmt.Printf("%s%s: %s\n", spaces, values[fileType], name)
}

func printTree(dest string, nest int) {
	files, err := os.ReadDir(dest)
	if err != nil {
		log.Fatal(err)
	}

	sort.Slice(files, func(i, j int) bool {
		first := files[i]
		second := files[j]

		if first.IsDir() != second.IsDir() {
			return !first.IsDir()
		}

		return first.Name() < first.Name()
	})

	for _, file := range files {
		isDir := file.IsDir()

		printObj(file.Name(), isDir, nest)

		if isDir {
			nestDir := filepath.Join(dest, file.Name())
			printTree(nestDir, nest+1)
		}
	}

}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file>")
		os.Exit(1)
	}

	path := os.Args[1]

	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Error: Path %s is invalid!\n\n", path)
		} else {
			fmt.Printf("Error: Cannot find such path... (%v)\n\n", err)
		}
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		fmt.Printf("Error: %s - file, not a directory.\n", path)
		os.Exit(1)
	}

	fmt.Printf("--- Tree of %s ---\n\n", path)
	printTree(path, 0)

}
