package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ANSI Пресеты GitHub Dark (TrueColor RGB)
const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	GHDarkText = "\033[38;2;201;209;217m"
	GHBlue     = "\033[38;2;88;166;255m"
	GHGreen    = "\033[38;2;63;185;80m"
	GHRed      = "\033[38;2;248;81;73m"
	GHOrange   = "\033[38;2;240;138;93m"
	GHPurple   = "\033[38;2;188;121;243m"
)

type treePrinter struct {
	colors [5]string
}

func (tp *treePrinter) printObj(name string, isDir bool, nest int) {

	values := map[bool]string{
		true:  fmt.Sprintf("%sdir: ", Bold),
		false: "file:",
	}

	const wide = 2

	spaces := strings.Repeat(" ", nest*wide)
	objColor := tp.colors[nest%len(tp.colors)]

	fmt.Printf("%s%s%s%s %s%s%s\n", spaces, GHDarkText, values[isDir], Reset, objColor, name, Reset)
}

func (tp *treePrinter) walkTree(dest string, nest int) {
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

		tp.printObj(file.Name(), isDir, nest)

		if isDir {
			nestDir := filepath.Join(dest, file.Name())
			tp.walkTree(nestDir, nest+1)
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

	printer := treePrinter{
		colors: [5]string{GHBlue, GHGreen, GHRed, GHPurple, GHOrange},
	}

	fmt.Printf("--- Tree of %s ---\n\n", path)
	printer.walkTree(path, 0)

}
