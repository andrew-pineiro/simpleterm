package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func get_newline() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}
func write(msg string) {
	fmt.Printf(">> %s", msg)
}
func try_execute(name string) bool {
	_, err := exec.LookPath(name)
	if err != nil {
		return false
	}
	cmd := exec.Command(name)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	if err := cmd.Run(); err != nil {
		fmt.Printf("Unable to run program %s - %s\n", name, err)
	}
	return true

}
func main() {

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")

		buffer, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		cmd := strings.Replace(buffer, get_newline(), "", -1)

		//SYSTEM APPS
		if try_execute(cmd) {
			continue
		}

		//CUSTOM CMDS
		switch cmd {
		case "exit":
			os.Exit(0)
		case "hello":
			write("Hello World!\n")
		default:
			write(fmt.Sprintf("Command not implemented %s", buffer))
		}
	}
}
