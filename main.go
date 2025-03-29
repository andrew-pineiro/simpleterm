package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

func get_newline() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
func try_execute(buf string) bool {

	program := buf
	if strings.Contains(buf, " ") {
		program = strings.Split(buf, " ")[0]
	}
	_, err := exec.LookPath(program)
	if err != nil {
		return false
	}
	cmd := exec.Command(program)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Args = append(cmd.Args, strings.Split(buf, " ")[1:]...)

	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	cmd.Run()
	return true

}
func main() {

	//curr_dir, _ := os.Getwd()
	curr_user, _ := user.Current()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s> ", curr_user.Username)

		buffer, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			os.Exit(0)
		}
		if err != nil {
			panic(err)
		}
		if len(strings.TrimSpace(buffer)) > 0 {
			cmd := strings.Replace(buffer, get_newline(), "", -1)

			//SYSTEM APPS
			if try_execute(cmd) {
				continue
			}

			//CUSTOM CMDS
			switch cmd {
			case "exit":
				os.Exit(0)
			default:
				fmt.Printf("%s\n", cmd)
			}
		}

	}
}
