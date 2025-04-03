package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type stFile struct {
	fileName    string
	fileModTime string
	fileMode    fs.FileMode
	fileIsDir   bool
	fileSize    string
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

func ls(args []string) {
	showHidden := false
	rootDir, _ := os.Getwd()
	if len(args) > 0 {
		for _, arg := range args {
			if arg == "-h" || arg == "-la" {
				showHidden = true
			}
			if arg[0] != '-' {
				_, err := os.Stat(arg)
				if !errors.Is(err, fs.ErrNotExist) {
					rootDir = arg
				}
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
