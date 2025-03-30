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

const TEMP_FILE_NAME = "sterm.tmp"

type context struct {
	working_dir string
	cmd         string
	args        []string
}

var c = &context{}

func filterInput(r rune) (rune, bool) {
	switch r {
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}
func try_execute() bool {

	program := c.cmd

	_, err := exec.LookPath(program)
	if err != nil {
		return false
	}
	cmd := exec.Command(program)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = c.args

	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	cmd.Run()
	return true

}

func try_cmd() bool {
	switch c.cmd {
	case "cd":
		cd()
	case "ls", "dir":
		ls()
	case "echo":
		echo()
	case "cp":
		cp()
	default:
		return false
	}
	return true
}
func main() {
	//curr_user, _ := user.Current()
	c.working_dir, _ = os.Getwd()
	temp_dir, _ := os.MkdirTemp("", ".sterm")
	cfg := &readline.Config{
		HistoryFile:     path.Join(temp_dir, "sterm.tmp"),
		InterruptPrompt: "^C",
		AutoComplete:    completer,
		EOFPrompt:       "exit",

		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	}
	reader, err := readline.NewEx(cfg)
	if err != nil {
		panic(err)
	}
	defer reader.Close()
	defer os.Remove(temp_dir)
	reader.CaptureExitSignal()

	for {
		reader.SetPrompt(c.working_dir + "> ")
		line, err := reader.Readline()
		if err == io.EOF {
			break
		}
		line = strings.TrimSpace(line)

		//CHECK FOR BLANK
		if len(line) == 0 {
			continue
		}
		c.cmd = line
		if strings.Contains(line, " ") {
			c.cmd = strings.Split(line, " ")[0]
		}
		c.args = strings.Split(line, " ")[1:]
		//EXIT
		if c.cmd == "exit" {
			goto exit
		}

		if try_cmd() {
			continue
		}
		if try_execute() {
			continue
		}
		fmt.Printf("Command not found %s\n", c.cmd)
	}
exit:
}
