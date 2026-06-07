package main

import (
	"bufio"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

		if len(fields) >= 4 && strings.Contains(fields[0], "/dev/") {
			drives = append(drives, fields[0])
		}
	}
	return drives, nil
}

func sortByNameAsc(entries []stFile) {
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].fileName) < strings.ToLower(entries[j].fileName)
	})
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

func getDiskSpaceAvailable(drives []string) map[string]uint64 {
	disks := make(map[string]uint64)
	for _, v := range drives {
		var stat unix.Statfs_t
		unix.Statfs(v, &stat)
		disks[v] = (stat.Bfree * uint64(stat.Bsize)) / 1024 / 1024 / 1024
	}
	return disks
}
