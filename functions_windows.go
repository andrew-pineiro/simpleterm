package main

import (
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

func sortByNameAsc(entries []stFile) {
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].fileName) < strings.ToLower(entries[j].fileName)
	})
}

// checks if a file is hidden on Windows.
func isHidden(path string, isDir bool) bool {
	//const dotChar = 46
	// if filepath.Base(path)[0] == dotChar && isDir {
	// 	return true
	// }
	if strings.Contains(strings.ToLower(filepath.Base(path)), "ntuser") {
		return true
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Appending `\\?\` to the absolute path helps with
	// preventing 'Path Not Specified Error' when accessing long paths and filenames
	// https://docs.microsoft.com/en-us/windows/win32/fileio/maximum-file-path-limitation?tabs=cmd

	pointer, err := syscall.UTF16PtrFromString(`\\?\` + absPath)
	if err != nil {
		return false
	}

	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false
	}

	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}
