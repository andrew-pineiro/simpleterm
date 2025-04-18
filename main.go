package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
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
func tryExecute(program string, args []string) bool {
	_, err := exec.LookPath(program)
	if err != nil {
		return false
	}

	//This is required because of the way exec handles args
	//Requires program to be first argument
	var _args []string
	_args = append(_args, program)
	for i := 0; i < len(args); i++ {
		_args = append(_args, args[i])

	}

	cmd := exec.Command(program)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Args = _args

	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	cmd.Run()
	return true

}

func tryCmd(cmd string, args []string) bool {
	switch cmd {
	case "rm":
		rm(args)
	case "cd":
		cd(args)
	case "ls", "dir":
		ls(args)
	case "echo":
		echo(args)
	case "cp":
		cp(args)
	case "cat":
		cat(args)
	case "pwd":
		wd, _ := os.Getwd()
		fmt.Printf("%s\n", wd)
	default:
		return false
	}
	return true
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

		cmd := strings.TrimSpace(line)
		var args []string
		if strings.Contains(line, " ") {
			cmd = strings.Split(line, " ")[0]
			args = strings.Split(line, " ")[1:]
		}

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
	//os.RemoveAll(tempDir)
}
