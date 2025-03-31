package main

import (
	"os"

	"github.com/chzyer/readline"
)

var completer = readline.NewPrefixCompleter(
	readline.PcItem("exit"),
	readline.PcItem("ls"),
	readline.PcItem("dir"),
	readline.PcItem("cp",
		readline.PcItemDynamic(list_files(c.working_dir))),
	readline.PcItem("cd",
		readline.PcItemDynamic(list_files(c.working_dir))),
	readline.PcItem("echo"),
)

func list_files(path string) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := os.ReadDir(path)
		for _, f := range files {
			names = append(names, f.Name())
		}
		return names
	}
}
