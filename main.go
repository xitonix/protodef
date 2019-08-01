package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jhump/protoreflect/desc/protoparse"
)

func main() {
	root := flag.String("root", "./src", "Root of proto files")
	flag.Parse()

	importPaths, err := getImportPaths(*root)
	if err != nil {
		fatal(err, 1)
	}

	print(importPaths, "Import Paths")

	files, err := protoFiles(*root)
	if err != nil {
		fatal(err, 2)
	}

	print(files, "Proto Files")

	resolved, err := protoparse.ResolveFilenames(importPaths, files...)
	if err != nil {
		fatal(err, 3)
	}

	print(resolved, "Proto Files")

	parser := protoparse.Parser{
		ImportPaths: importPaths,
	}

	fileDescriptors, err := parser.ParseFiles(resolved...)
	if err != nil {
		fatal(err, 4)
	}

	fmt.Println("\nParsed Messages:")
	for _, fd := range fileDescriptors {
		fmt.Printf(" %s\n", fd.GetName())
		messages := fd.GetMessageTypes()
		for _, m := range messages {
			fmt.Printf("    %s\n", m.GetName())
		}
	}
}

func print(in []string, msg string) {
	fmt.Println(msg + ":")
	for i, p := range in {
		fmt.Printf("    %d: %s\n", i, p)
	}
}

func fatal(err error, code int) {
	fmt.Printf("E%d: %s\n", code, err)
	os.Exit(code)
}

func getImportPaths(root string) ([]string, error) {
	var dirs []string
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			dirs = append(dirs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func protoFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".proto") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
