package main

import (
	"io/fs"
	"sort"
	"strings"
)

type stFile struct {
	fileName    string
	fileModTime string
	fileMode    fs.FileMode
	fileIsDir   bool
	fileSize    string
}

type stDisk struct {
	driveName  string
	totalSpace uint64
	usedSpace  uint64
	availSpace uint64
}

func sortByNameAsc(entries []stFile) {
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].fileName) < strings.ToLower(entries[j].fileName)
	})
}

func bytesToString(b []byte) string {
	for i, c := range b {
		if c == 0 {
			return string(b[:i])
		}
	}
	return string(b)
}
