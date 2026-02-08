package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/chzyer/readline"
)

const HIST_FILENAME = "sterm.hst"

func filterInput(r rune) (rune, bool) {
	switch r {
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func newRlInstance() *readline.Instance {
	//currDir, _ := os.Getwd()
	homeDir, _ := os.UserHomeDir()
	histDir := path.Join(homeDir, ".sterm")
	_ = os.Mkdir(histDir, 0700)
	cfg := &readline.Config{
		HistoryFile:     path.Join(histDir, HIST_FILENAME),
		InterruptPrompt: "^C",
		AutoComplete:    createCompleter(),
		EOFPrompt:       "",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	}
	reader, err := readline.NewEx(cfg)
	if err != nil {
		panic(err)
	}
	return reader
}

func parseCmdArgs(line string) (string, []string) {
	var args []string

	l := strings.Split(line, " ")
	cmd := l[0]
	if len(l) > 1 {
		a := l[1:]
		for i := 0; i < len(a); i++ {
			buffer := a[i]
			if strings.TrimSpace(a[i]) != "" && a[i][0] == '"' {
				for j := i + 1; j < len(a); j++ {
					buffer += fmt.Sprintf(" %s", a[j])

					i = j + 1

					//TODO(#2): check for quote in middle of string
					if strings.ContainsRune(a[j], '"') {
						break
					}

					continue
				}
			}
			args = append(args, buffer)
		}
	}
	return cmd, args
}
func main() {
	//currUser, _ := user.Current()
	homeDir, _ := os.UserHomeDir()
	os.Chdir(homeDir)

	reader := newRlInstance()
	reader.CaptureExitSignal()

	defer reader.Close()

	for {
		wd, _ := os.Getwd()
		reader.SetPrompt(wd + "> ")

		line, err := reader.Readline()
		if err == io.EOF {
			goto exit
		}

		//CHECK FOR INTERRUPT OR BLANK
		if err == readline.ErrInterrupt || len(line) == 0 {
			continue
		}

		cmd, args := parseCmdArgs(line)

		//EXIT
		if cmd == "exit" {
			goto exit
		}

		if tryCmd(cmd, args) {
			reader.Config.AutoComplete = createCompleter()
			continue
		}

		if tryExecute(cmd, args) {
			continue
		}

		fmt.Printf("Command not found %s\n", cmd)
	}
exit:
	return
}
