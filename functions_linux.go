package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/sys/unix"
)

func getDrives() ([]string, error) {
	var drives []string

	file, err := os.Open("/proc/mounts")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) >= 4 && strings.Contains(fields[0], "/dev/") && !slices.Contains(drives, fields[0]) {
			drives = append(drives, fields[0])
		}
	}
	return drives, nil
}

func isHidden(path string, _ bool) bool {
	const dotChar = 46
	base := filepath.Base(path)

	return base[0] == dotChar
}

func isWSL() bool {
	releaseData, err := exec.Command("uname", "-r").Output()
	if err != nil {
		return false
	}
	return strings.Contains(strings.ToLower(string(releaseData)), "microsoft")
}

func getDiskSpaceAvailable(drives []string) []stDisk {
	var disks []stDisk
	for _, v := range drives {
		var stat unix.Statfs_t
		unix.Statfs(v, &stat)
		disks = append(disks, stDisk{
			driveName:  v,
			availSpace: stat.Bfree * uint64(stat.Bsize),
			totalSpace: stat.Bavail * uint64(stat.Bsize),
			usedSpace:  (stat.Bavail - stat.Bfree) * uint64(stat.Bsize),
		})
	}
	return disks
}
