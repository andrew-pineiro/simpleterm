package main

import (
	"os"

	"github.com/chzyer/readline"
)

func createCompleter() *readline.PrefixCompleter {
	wd, _ := os.Getwd()
	return readline.NewPrefixCompleter(
		readline.PcItem("exit"),
		readline.PcItem("cat",
			readline.PcItemDynamic(listFiles(wd, false, true))),
		readline.PcItem("ls",
			readline.PcItemDynamic(listFiles(wd, true, true))),
		readline.PcItem("dir",
			readline.PcItemDynamic(listFiles(wd, true, true))),
		readline.PcItem("cp",
			readline.PcItemDynamic(listFiles(wd, false, true))),
		readline.PcItem("cd",
			readline.PcItemDynamic(listFiles(wd, true, true))),
		readline.PcItem("echo"),
		readline.PcItem("rm",
			readline.PcItemDynamic(listFiles(wd, false, true))),
	)
}

func listFiles(path string, onlyDir bool, showHidden bool) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := os.ReadDir(path)
		for _, f := range files {
			if onlyDir && !f.IsDir() {
				continue
			}
			if !showHidden && isHidden(f.Name(), onlyDir) {
				continue
			}
			names = append(names, f.Name())
		}
		return names
	}
}
