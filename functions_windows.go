package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"

	"golang.org/x/sys/windows"
)

func isWSL() bool {
	return false
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

func getDrives() ([]string, error) {
	n, err := windows.GetLogicalDriveStrings(0, nil)
	if err != nil {
		return nil, fmt.Errorf("failed getting buffer size: %w", err)
	}

	buf := make([]uint16, n)
	_, err = windows.GetLogicalDriveStrings(uint32(len(buf)), &buf[0])
	if err != nil {
		return nil, fmt.Errorf("failed retrieving drive strings: %w", err)
	}

	var drives []string
	from := 0
	for i, v := range buf {
		if v == 0 {
			if i > from {
				drives = append(drives, string(windows.UTF16ToString(buf[from:i])))
			}
			from = i + 1
		}
	}

	return drives, nil
}
func getDiskSpaceAvailable(drives []string) []stDisk {
	var disks []stDisk
	for _, v := range drives {
		var freeBytesAvailable uint64
		var totalNumberOfBytes uint64
		var totalNumberOfFreeBytes uint64
		driveName := strings.ReplaceAll(v, ":\\", "")
		windows.GetDiskFreeSpaceEx(windows.StringToUTF16Ptr(driveName+":"),
			&freeBytesAvailable, &totalNumberOfBytes, &totalNumberOfFreeBytes)
		disks = append(disks, stDisk{
			driveName:  driveName,
			totalSpace: totalNumberOfBytes,
			availSpace: totalNumberOfFreeBytes,
			usedSpace:  totalNumberOfBytes - totalNumberOfFreeBytes,
		})
	}

	return disks
}
func getMemoryInfo() (stMem, error) {
	var mem stMem
	fmt.Println("not implemented for Windows")
	return mem, nil
}
