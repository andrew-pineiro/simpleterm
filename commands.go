package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var commands = map[string]func([]string){
	"rm":      rm,
	"cd":      cd,
	"ls":      ls,
	"dir":     ls, // alias
	"echo":    echo,
	"cp":      cp,
	"cat":     cat,
	"ping":    ping,
	"checkip": checkip,
	"pwd": func(args []string) {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println(wd)
	},
}

type stFile struct {
	fileName    string
	fileModTime string
	fileMode    fs.FileMode
	fileIsDir   bool
	fileSize    string
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
	if fn, ok := commands[cmd]; ok {
		fn(args)
		return true
	}
	return false
}

func echo(args []string) {
	if len(args) == 0 {
		return
	}
	for i := range args {
		fmt.Printf("%s ", args[i])
	}
	fmt.Println()
}
func cp(args []string) {
	if len(args) < 2 {
		//TODO: print help for cp
		return
	}
	source := args[0]
	dest := args[1]
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		fmt.Printf("ERROR: %s does not exist\n", source)
		return
	}
	contents, err := os.ReadFile(source)
	if err != nil {
		fmt.Printf("ERROR: unable to copy file %s\n", source)
		return
	}
	os.WriteFile(dest, contents, os.FileMode(os.O_CREATE))
}
func cd(args []string) {
	if len(args) < 1 {
		//TODO: print help for cd
		fmt.Printf("ERROR: invalid arguments.")
		return
	}
	newPath, err := filepath.Abs(args[0])
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("ERROR: %s does not exist\n", newPath)
		return
	}
	_, err = os.Stat(newPath)
	if errors.Is(err, fs.ErrNotExist) {
		fmt.Printf("ERROR: %s does not exist\n", newPath)
		return
	}
	if errors.Is(err, fs.ErrPermission) {
		fmt.Printf("ERROR: permission denied to %s", newPath)
		return
	}
	os.Chdir(newPath)
}
func checkip(args []string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		//TODO: handle error for checkip
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			//TODO: handle error for check ip addrs
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.DefaultMask() == nil {
				fmt.Printf("Interface: %s IPv6: %s\n", i.Name, ip.String())
			} else {
				fmt.Printf("Interface: %s IPv4: %s\n", i.Name, ip.String())
			}
		}
	}
}
func ping(args []string) {
	if len(args) < 1 {
		//TODO: print help for ping
		return
	}
	//TODO: add more args
	ip := args[0]
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		fmt.Printf("ERROR: unable to establish socket - %s\n", err.Error())
		return
	}
	defer conn.Close()

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  1,
			Data: []byte("000000000000000"),
		},
	}

	bytes, err := msg.Marshal(nil)
	if err != nil {
		fmt.Printf("ERROR: could not marshal content - %s\n", err.Error())
		return
	}
	for i := 0; i < 4; i++ {
		time.Sleep(time.Second * 1)
		_, err := conn.WriteTo(bytes, &net.IPAddr{IP: net.ParseIP(ip)})
		if err != nil {
			fmt.Printf("ERROR: could not open connection - %s\n", err.Error())
			return
		}
		fmt.Printf("Successfully connected to %s\n", ip)
	}

}
func ls(args []string) {
	showHidden := false
	rootDir, _ := os.Getwd()
	if len(args) > 0 {
		for _, arg := range args {
			if arg == "-h" || arg == "-la" {
				showHidden = true
				continue
			}
			_, err := os.Stat(arg)
			if !errors.Is(err, fs.ErrNotExist) {
				rootDir = arg
			}
		}
	}
	absDir, _ := filepath.Abs(rootDir)
	var dirs []stFile
	var files []stFile
	f, err := os.ReadDir(absDir)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	for _, file := range f {
		fileInfo, _ := file.Info()
		var newFile = stFile{
			fileName:    file.Name(),
			fileModTime: fileInfo.ModTime().Format("1/2/2006 3:04 PM"),
			fileMode:    fileInfo.Mode(),
			fileIsDir:   file.IsDir(),
			fileSize:    strconv.FormatInt(fileInfo.Size(), 10),
		}
		if newFile.fileIsDir || strings.HasPrefix(newFile.fileMode.String(), "L") {
			dirs = append(dirs, newFile)
		} else {
			files = append(files, newFile)
		}
	}
	sortByNameAsc(dirs)
	sortByNameAsc(files)

	if len(files)+len(dirs) <= 0 {
		return
	}

	fmt.Printf("\n  Directory: %s\n\n", absDir)
	maxModeWidth := len("Mode")
	maxDateWidth := len("Modified")
	maxSizeWidth := len("Size")
	maxNameWidth := len("Name")

	var allEntries []stFile
	allEntries = append(allEntries, dirs...)
	allEntries = append(allEntries, files...)

	for _, entry := range allEntries {
		if !showHidden && isHidden(entry.fileName, entry.fileIsDir) {
			continue
		}
		mode := entry.fileMode.String()
		modified := entry.fileModTime
		size := entry.fileSize
		name := entry.fileName
		if len(mode) > maxModeWidth {
			maxModeWidth = len(mode)
		}
		if len(modified) > maxDateWidth {
			maxDateWidth = len(modified)
		}
		if len(size) > maxSizeWidth {
			maxSizeWidth = len(size)
		}
		if len(name) > maxNameWidth {
			maxNameWidth = len(name)
		}
	}

	fmt.Printf("%-*s  %-*s  %*s  %-*s\n",
		maxModeWidth, "Mode",
		maxDateWidth, "Modified",
		maxSizeWidth, "Size",
		maxNameWidth, "Name")

	fmt.Printf("%s  %s  %s  %s\n",
		strings.Repeat("-", maxModeWidth),
		strings.Repeat("-", maxDateWidth),
		strings.Repeat("-", maxSizeWidth),
		strings.Repeat("-", maxNameWidth))

	for _, entry := range dirs {
		size := " "
		name := fmt.Sprintf("\033[37;44;1m%s\033[0m", entry.fileName) // Blue background for directories

		if !showHidden {
			path, _ := filepath.Abs(entry.fileName)
			if isHidden(path, true) {
				continue
			}
		}

		fmt.Printf("%-*s  %-*s  %*s  %-*s\n",
			maxModeWidth, entry.fileMode,
			maxDateWidth, entry.fileModTime,
			maxSizeWidth, size,
			maxNameWidth, name)
	}

	for _, entry := range files {
		if !showHidden {
			path, _ := filepath.Abs(entry.fileName)
			if isHidden(path, false) {
				continue
			}
		}

		fmt.Printf("%-*s  %-*s  %*s  %-*s\n",
			maxModeWidth, entry.fileMode,
			maxDateWidth, entry.fileModTime,
			maxSizeWidth, entry.fileSize,
			maxNameWidth, entry.fileName)
	}
	fmt.Println()
}

func rm(args []string) {
	if len(args) < 1 {
		fmt.Printf("ERROR: invalid arguments.")
		return
	}
	rawPath := args[0]
	if strings.Contains(rawPath, "*") {
		//TODO(#1): implement wildcard matching
	}

	path, _ := filepath.Abs(rawPath)
	_, err := os.Stat(path)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
		return
	}

	err = os.Remove(path)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
}

func cat(args []string) {
	if len(args) < 1 {
		return
	}
	file := args[0]
	_, err := os.Stat(file)
	if err != nil {
		return
	}

	content, _ := os.ReadFile(file)
	fmt.Printf("%s", content)
}
