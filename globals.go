package main

import (
	"io/fs"
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
