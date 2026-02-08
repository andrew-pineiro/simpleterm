package main

import (
	"os"

	"github.com/chzyer/readline"
)

func createCompleter() *readline.PrefixCompleter {
	wd, _ := os.Getwd()
	completer := readline.NewPrefixCompleter(
		readline.PcItem("exit"),
		readline.PcItem("file",
			readline.PcItemDynamic(listFiles(wd, false, true, true))),
		readline.PcItem("cat",
			readline.PcItemDynamic(listFiles(wd, false, true, false))),
		readline.PcItem("ls",
			readline.PcItemDynamic(listFiles(wd, true, true, false))),
		readline.PcItem("dir",
			readline.PcItemDynamic(listFiles(wd, true, true, false))),
		readline.PcItem("cp",
			readline.PcItemDynamic(listFiles(wd, false, true, false))),
		readline.PcItem("cd",
			readline.PcItemDynamic(listFiles(wd, true, true, false))),
		readline.PcItem("echo"),
		readline.PcItem("ping"),
		readline.PcItem("checkip"),
		readline.PcItem("open"),
		readline.PcItem("track"),
		readline.PcItem("rm",
			readline.PcItemDynamic(listFiles(wd, false, true, false))),
	)
	return completer
}

func listFiles(path string, onlyDir bool, showHidden bool, onlyFile bool) func(string) []string {
	return func(line string) []string {
		names := make([]string, 0)
		files, _ := os.ReadDir(path)
		for _, f := range files {
			if onlyDir && !f.IsDir() {
				continue
			}
			if onlyFile && f.IsDir() {
				continue
			}
			if !showHidden && isHidden(f.Name(), onlyDir) {
				continue
			}
			name := f.Name()
			//TODO: get autocomplete to work with spaces on folder / file names
			// if strings.ContainsAny(name, " ") {
			// 	name = "\"" + name + "\""
			// }
			names = append(names, name)
		}
		return names
	}
}
