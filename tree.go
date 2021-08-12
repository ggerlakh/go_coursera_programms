package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"strconv"
)

func getTab(path string, elem string, m map[string][]string) string {
	var tab string
	elem = filepath.Dir(elem)
	for elem != path {
		if elem == m[filepath.Dir(elem)][len(m[filepath.Dir(elem)]) - 1] {
			tab = "\t" + tab
		} else {
			tab = "│\t" + tab
		}
		elem = filepath.Dir(elem)
	}
	return tab
}

func printTree(out io.Writer, sl []string, path string) error {
	m := make(map[string][]string)
	last := "└───"
	graphic := "├───"
	for idx := range sl {
		m[filepath.Dir(sl[idx])] = append(m[filepath.Dir(sl[idx])], sl[idx])
	}
	for i := range sl {
		if sl[i] == m[filepath.Dir(sl[i])][len(m[filepath.Dir(sl[i])]) - 1] {
			fmt.Fprintln(out, getTab(path, sl[i], m) + last + filepath.Base(sl[i]))
		} else {
			fmt.Fprintln(out, getTab(path, sl[i], m) + graphic + filepath.Base(sl[i]))
		}
	}
	return nil
}

func dirTree(output io.Writer, path string, printFiles bool) error {
	var lines string
	var fileLines []string
	if printFiles {
		filepath.Walk(path, func(line string, info os.FileInfo, err error) error {
			if line != "." && line != path {
				if !(info.IsDir()) {
					if info.Size() != 0 {
						line = line + " (" + strconv.Itoa(int(info.Size())) + "b" + ")"
					} else {
						line = line + " (" + "empty" + ")"
					}
				}
				fileLines = append(fileLines, line)
			}
			return nil
		})
		return printTree(output, fileLines, path)
	} else {
		filepath.Walk(path, func(line string, info os.FileInfo, err error) error {
			dirLine := filepath.Dir(line)
			if dirLine != "." && !(strings.Contains(lines, dirLine)) && dirLine != path {
				lines = lines + " " + dirLine
			}
			return nil
		})
		slice := strings.Split(lines, " ")[1:]
		return printTree(output, slice, path)
	}
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
