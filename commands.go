package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"math/big"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

type command struct {
	fn          func([]string)
	description string
	usage       string
}

func (c command) printHelp(name string) {
	fmt.Printf("Usage: %s %s\n", name, c.usage)
	fmt.Printf("  %s\n", c.description)
}

var commands = map[string]command{
	"rm": {
		fn:          rm,
		description: "Remove a file or files matching a pattern",
		usage:       "<file|pattern>",
	},
	"cd": {
		fn:          cd,
		description: "Change the current working directory",
		usage:       "<path>",
	},
	"ls": {
		fn:          ls,
		description: "List directory contents",
		usage:       "[-h|-la] [path] [*.ext]",
	},
	"dir": {
		fn:          ls,
		description: "List directory contents (alias for ls)",
		usage:       "[-h|-la] [path] [*.ext]",
	},
	"file": {
		fn:          file,
		description: "Show the MIME type of a file",
		usage:       "<file>",
	},
	"echo": {
		fn:          echo,
		description: "Print arguments to stdout",
		usage:       "<text...>",
	},
	"cp": {
		fn:          cp,
		description: "Copy a file to a destination",
		usage:       "<source> <dest>",
	},
	"cat": {
		fn:          cat,
		description: "Print the contents of a file",
		usage:       "<file>",
	},
	"ping": {
		fn:          ping,
		description: "Send ICMP echo requests to an IP address (4 pings)",
		usage:       "<ip>",
	},
	"checkip": {
		fn:          checkip,
		description: "Display local network interfaces and public IP address",
		usage:       "",
	},
	"track": {
		fn:          track,
		description: "Open a package tracking page for a UPS, FedEx, or USPS tracking number",
		usage:       "<tracking-number>",
	},
	"open": {
		fn:          open,
		description: "Open a URL in the default browser",
		usage:       "<url>",
	},
	"disk": {
		fn:          disk,
		description: "Show disk space usage for all drives",
		usage:       "",
	},
	"pwd": {
		fn: func(args []string) {
			wd, err := os.Getwd()
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Println(wd)
		},
		description: "Print the current working directory",
		usage:       "",
	},
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
	if strings.ToLower(cmd) == "help" {
		for _, cmd := range slices.Collect(maps.Keys(commands)) {
			c := commands[cmd]
			fmt.Printf("%s\n", cmd)

			c.printHelp(cmd)
			fmt.Println()
		}
		return true
	}

	if c, ok := commands[cmd]; ok {
		if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
			c.printHelp(cmd)
			return true
		}
		c.fn(args)
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
	var output string
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("ERROR: Unable to establish network interfaces: %s\n", err)
		return
	}

	for _, i := range interfaces {
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Printf("ERROR: Cannot retrieve address from interface %s: %s\n", i.Name, err)
			return
		}
		output = "Interface: " + i.Name
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.DefaultMask() == nil {
				output += " - IPv6: " + ip.String()
			} else {
				output += " - IPv4: " + ip.String()
			}
		}
		fmt.Println(output)
	}

	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		fmt.Printf("ERROR: Unable to get public IP: %s\n", err)
		return
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("ERROR: Unable to read response: %s\n", err)
		return
	}
	fmt.Printf("Public IP: %s\n", string(ip))
}
func ping(args []string) {
	if len(args) < 1 {
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
func file(args []string) {
	if len(args) > 0 {
		absPath, _ := filepath.Abs(args[0])
		file, err := os.Stat(absPath)
		if err != nil {
			fmt.Printf("ERROR: could not read file - %s\n", err.Error())
			return
		}
		if file.IsDir() {
			return
		}
		mtype, err := mimetype.DetectFile(absPath)
		if err != nil {
			fmt.Printf("ERROR: could not parse file type - %s\n", err.Error())
			return
		}
		fmt.Printf("%s - %s\n", file.Name(), mtype.String())
	}
}
func ls(args []string) {
	showHidden := false
	wildCardExten := ""
	rootDir, _ := os.Getwd()
	if len(args) > 0 {
		for _, arg := range args {
			if arg == "-h" || arg == "-la" {
				showHidden = true
				continue
			}
			if strings.HasPrefix(arg, "*") && len(arg) > 2 {
				wildCardExten = strings.Replace(arg, "*.", "", -1)
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
		if wildCardExten != "" && !strings.HasSuffix(newFile.fileName, wildCardExten) {
			continue
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
		fmt.Printf("ERROR: invalid arguments.\n")
		return
	}
	var files []string
	rawPath := args[0]
	if strings.Contains(rawPath, "*") {
		wd, _ := os.Getwd()
		path, _ := filepath.Abs(wd)
		if len(rawPath) > strings.Index(rawPath, "*")+1 && string(rawPath[strings.Index(rawPath, "*")+1]) == "." {
			ext := string(filepath.Ext(rawPath))
			f, err := os.ReadDir(path)
			if err != nil {
				fmt.Printf("ERROR: unable to read directory - %s\n", err)
				return
			}
			for _, file := range f {
				if strings.HasSuffix(file.Name(), ext) {
					absFilePath, _ := filepath.Abs(file.Name())
					files = append(files, absFilePath)
				}
			}
		} else if strings.Replace(rawPath, "*", "", -1) != "" {
			if strings.HasSuffix(rawPath, "*") {
				wordMatch := strings.Replace(rawPath, "*", "", -1)
				f, err := os.ReadDir(path)
				if err != nil {
					fmt.Printf("ERROR: unable to read directory - %s\n", err)
					return
				}
				for _, file := range f {
					if strings.Contains(file.Name(), wordMatch) {
						absFilePath, _ := filepath.Abs(file.Name())
						files = append(files, absFilePath)
					}
				}
			}
		} else {
			f, err := os.ReadDir(path)
			if err != nil {
				fmt.Printf("ERROR: unable to read directory - %s\n", err)
				return
			}
			for _, file := range f {
				absFilePath, _ := filepath.Abs(file.Name())
				files = append(files, absFilePath)
			}
		}
	} else {
		path, _ := filepath.Abs(rawPath)
		_, err := os.Stat(path)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			return
		}
		files = append(files, path)
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
		}
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

func track(args []string) {
	if len(args) < 1 {
		return
	}
	trackId := strings.ToUpper(args[0])

	if len(trackId) <= 0 {
		return
	}

	var fedexLengths = map[int]struct{}{
		10: {},
		12: {},
		15: {},
		20: {},
		22: {},
		34: {},
	}

	if strings.HasPrefix(trackId, "1Z") || len(trackId) == 9 || strings.HasPrefix(trackId, "K") || strings.HasPrefix(trackId, "H") {
		open([]string{"https://www.ups.com/track?tracknum=" + trackId + "&AgreeToTermsAndConditions=yes&loc=en_US"})
	} else if _, ok := new(big.Int).SetString(trackId, 10); ok {
		if (strings.HasPrefix(trackId, "94") && len(trackId) == 22) || strings.HasSuffix(trackId, "US") {
			open([]string{"http://trkcnfrm1.smi.usps.com/PTSInternetWeb/InterLabelInquiry.do?origTrackNum=" + trackId})
		} else if _, ok := fedexLengths[len(trackId)]; ok {
			open([]string{"https://fedex.com/fedextrack/?tracknumbers=" + trackId})
		}

	} else {
		fmt.Println("ERROR: Unable to verify carrier for tracking id: " + trackId)
	}
}
func open(args []string) {
	if len(args) < 1 {
		return
	}
	var cmd string
	var execArgs []string
	url := args[0]

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd.exe"
		execArgs = []string{"/c", "rundll32", "url.dll,FileProtocolHandler", strings.ReplaceAll(url, "&", "^&")}
	case "darwin":
		cmd = "open"
		execArgs = []string{url}
	default:
		if isWSL() {
			cmd = "cmd.exe"
			execArgs = []string{"start", url}
		} else {
			cmd = "xdg-open"
			execArgs = []string{url}
		}
	}

	e := exec.Command(cmd, execArgs...)
	err := e.Start()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
	err = e.Wait()
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	return
}

func disk(args []string) {
	drives, err := getDrives()
	if err != nil {
		fmt.Printf("ERROR: Unable to retreive drives - %s", err)
		return
	}
	slices.Sort(drives)
	maxDiskLen := len("Disk")
	maxSizeLen := len("Size GB")
	maxUsedLen := len("Used GB")
	maxAvailLen := len("Available")

	disks := getDiskSpaceAvailable(drives)
	for _, drive := range disks {
		if len(drive.driveName) > maxDiskLen {
			maxDiskLen = len(drive.driveName)
		}
	}
	fmt.Printf("%-*s  %-*s  %-*s  %-*s\n",
		maxDiskLen, "Disk",
		maxSizeLen, "Size",
		maxUsedLen, "Used",
		maxAvailLen, "Available")

	fmt.Printf("%s  %s  %s  %s\n",
		strings.Repeat("-", maxDiskLen),
		strings.Repeat("-", maxSizeLen),
		strings.Repeat("-", maxUsedLen),
		strings.Repeat("-", maxAvailLen))
	for _, drive := range disks {
		//TODO: implement a check for non-storage drives like disc drives.
		if drive.availSpace <= 0 {
			continue
		}
		tot := fmt.Sprint(drive.totalSpace/1024/1024/1024) + " GB"
		used := fmt.Sprint(drive.usedSpace/1024/1024/1024) + " GB"
		avail := fmt.Sprint(drive.availSpace/1024/1024/1024) + " GB"
		fmt.Printf("%-*s   %-*s  %-*s  %-*s\n",
			maxDiskLen, drive.driveName,
			maxSizeLen, tot,
			maxUsedLen, used,
			maxAvailLen, avail)
	}
	return
}
