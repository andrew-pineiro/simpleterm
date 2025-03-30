package main

import (
	"bufio"
	"errors"
	"fmt"
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
func try_cmd(buf string) bool {
	cmd := buf
	if strings.Contains(buf, " ") {
		cmd = strings.Split(buf, " ")[0]
	}
	args := strings.Split(buf, " ")[1:]
	if len(cmd) > 0 {
		switch cmd {
		case "exit":
			os.Exit(0)
		case "ls", "dir":
			ls()
		case "echo":
			echo(args)
		case "cp":
			cp(args)
		default:
			return false
		}
		return true
	}
	return false
}
func main() {

	//curr_dir, _ := os.Getwd()
	curr_user, _ := user.Current()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%s> ", curr_user.Username)

		buffer, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		cmd := strings.Replace(buffer, get_newline(), "", -1)
		if try_execute(cmd) {
			continue
		}
		if try_cmd(cmd) {
			continue
		}
		fmt.Printf("Command not found %s\n", cmd)
	}
}
