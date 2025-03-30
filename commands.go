package main

import (
	"fmt"
	"os"
)

func echo(args []string) {
	for i := range args {
		fmt.Printf("%s ", args[i])
	}
	fmt.Println()
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
		if file.Name()[0] == '.' {
			continue
		}
		fmt.Printf("%s\t%s\t%d\t%s\n", file.Mode(), file.ModTime().Format("1/2/2006 3:04 PM"), file.Size(), file.Name())
	}
	fmt.Println()
}
