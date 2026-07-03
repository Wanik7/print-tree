package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ANSI presets for GitHub Dark Theme (TrueColor RGB)
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
	colors  [5]string
	showAll bool
	noColor bool
	flat    bool
	wide    bool
}

func (tp *treePrinter) printObj(name string, isDir bool, nest int) {
	// determine count of spaces considering *flat* and *wide* flags
	indentSize := 2
	if tp.flat {
		indentSize = 0
	} else if tp.wide {
		indentSize = 4
	}

	spaces := strings.Repeat(" ", nest*indentSize)

	// if color are off, output plain text
	if tp.noColor {
		typeStr := "file:"
		if isDir {
			typeStr = "dir:"
		}
		fmt.Printf("%s%s %s\n", spaces, typeStr, name)
		return
	}

	// output logic
	values := map[bool]string{
		true:  fmt.Sprintf("%sdir: ", Bold),
		false: "file:",
	}

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
		return first.Name() < second.Name()
	})

	for _, file := range files {
		name := file.Name()

		// *a* flag: if obj hidden and a flag disabled => skip
		if strings.HasPrefix(name, ".") && !tp.showAll {
			continue
		}

		isDir := file.IsDir()
		tp.printObj(name, isDir, nest)

		if isDir {
			nestDir := filepath.Join(dest, name)
			tp.walkTree(nestDir, nest+1)
		}
	}
}

func main() {
	flag.Usage = func() {
		fmt.Println("CLI Tree Printer — Utility that shows file-tree structure")
		fmt.Println("\nUsage:")
		fmt.Println("  go run main.go [flags] <path_to_dir>")
		fmt.Println("\nFlags:")
		fmt.Println("  -a, --all      Show hidden files (e.g. .git)")
		fmt.Println("  -n, --no-color Turn off color output")
		fmt.Println("  -h, --help     Show this help message")
		fmt.Println("  -f, --flat     Turn off spaces in output")
		fmt.Println("  -w, --wide     Use 4 spaces instead of 2")
	}

	var showAll bool
	var noColor bool
	var flat bool
	var wide bool

	flag.BoolVar(&showAll, "a", false, "Show hidden files")
	flag.BoolVar(&showAll, "all", false, "Show hidden files")

	flag.BoolVar(&noColor, "n", false, "Turn off color output")
	flag.BoolVar(&noColor, "no-color", false, "Turn off color output")

	flag.BoolVar(&flat, "f", false, "Turn off spaces in output")
	flag.BoolVar(&flat, "flat", false, "Turn off spaces in output")

	flag.BoolVar(&wide, "w", false, "Use 4 spaces instead of 2")
	flag.BoolVar(&wide, "wide", false, "Use 4 spaces instead of 2")

	flag.Parse()

	args := flag.Args()

	// check for flag-like arguments after positional args
	// (Go's flag package stops parsing at the first non-flag argument)
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			fmt.Printf("%sError: Flag %q must be placed before the path argument!%s\n", GHRed, arg, Reset)
			fmt.Printf("Usage: go run main.go [flags] <path_to_dir>\n")
			os.Exit(1)
		}
	}

	// check *f* and *w* conflict
	if flat && wide {
		fmt.Printf("%sError: Flags --flat (-f) and --wide (-w) cannot be used together!%s\n", GHRed, Reset)
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Printf("%sError: There is no dir path!%s\n\n", GHRed, Reset)
		flag.Usage()
		os.Exit(1)
	}

	path := args[0]

	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%sError: Path %s is invalid!%s\n\n", GHRed, path, Reset)
		} else {
			fmt.Printf("%sError: Cannot find such path... (%v)%s\n\n", GHRed, err, Reset)
		}
		os.Exit(1)
	}

	if !fileInfo.IsDir() {
		fmt.Printf("%sError: %s - file, not a directory.%s\n", GHRed, path, Reset)
		os.Exit(1)
	}

	printer := treePrinter{
		colors:  [5]string{GHBlue, GHGreen, GHRed, GHPurple, GHOrange},
		noColor: noColor,
		showAll: showAll,
		flat:    flat,
		wide:    wide,
	}

	if !printer.noColor {
		fmt.Print(GHDarkText)
	}

	fmt.Printf("--- Tree of %s ---\n\n", path)
	printer.walkTree(path, 0)
}
