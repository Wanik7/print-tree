# CLI Tree Printer

A simple Go-based command-line utility for visualizing directory structures as a tree. The program features recursive directory traversal, prioritized file-first rendering,
and a GitHub Dark inspired color palette with syntax highlighting for specific configuration and source files.

## Features

- Smart Sorting: Files are always listed first, followed by directories. Alphabetical order is strictly preserved within each group.

- Custom Palette: Uses 24-bit TrueColor (RGB) ANSI escape codes to deliver the soft, eye-friendly tones of the GitHub Dark theme.

- Input Validation: Robust handling of command-line arguments, including checks for path existence and verifying that the target is a directory.

## Requirements

- Go 1.16 or higher

## Usage

- To run the utility, pass the path to the target directory as a command-line argument:

``` bash
go run main.go <path_to_directory>
```
- Example Command

``` bash
go run main.go ./my-project
```
- Example Output

``` bash
file: main.go
file: go.mod
dir:  config
  file: config.yaml
  file: app.json
dir:  internal
  file: router.go
  dir:  db
    file: connection.go
```
