package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

func echo(args []string) {
	for i := range args {
		fmt.Printf("%s ", args[i])
	}
	fmt.Println()
}
func cp(args []string) {
	if len(args) < 2 {
		return
	}
	source := args[0]
	dest := args[1]
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("ERROR: %s does not exist\n", source)
		return
	}
	contents, err := os.ReadFile(source)
	if err != nil {
		fmt.Printf("ERROR: unable to copy file %s\n", source)
		return
	}
	os.WriteFile(dest, contents, os.FileMode(os.O_CREATE))
}
func ls() {
	dir, _ := os.Getwd()
	files, err := os.ReadDir(".")
	if err != nil {
		return
	}
	if len(files) <= 0 {
		return
	}

	fmt.Printf("\n  Directory: %s\n\n", dir)
	fmt.Printf("Mode\t\tModified\t\tSize\tName\n")
	fmt.Printf("----\t\t--------\t\t----\t----\n")
	for _, file := range files {
		file, _ := file.Info()
		size := strconv.FormatInt(file.Size(), 10)
		name := file.Name()
		if file.IsDir() {
			size = " "
			name += string(os.PathSeparator)
		}
		if file.Name()[0] == '.' {
			continue
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", file.Mode(), file.ModTime().Format("1/2/2006 3:04 PM"), size, name)
	}
	fmt.Println()
}
