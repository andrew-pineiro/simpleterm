package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
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
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
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

func getMemoryInfo() (stMem, error) {
	var mem stMem

	////////

	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return mem, err
	}
	defer file.Close()
	metrics := make(map[string]int64)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 2 {
			continue
		}

		// Clean the key name (strip the trailing colon)
		key := strings.TrimSuffix(fields[0], ":")

		// Parse the numeric string into an integer
		valKb, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			continue
		}

		// Most /proc/meminfo values are in kB. Convert to bytes.
		metrics[key] = valKb * 1024
	}
	mem.totalMemory = metrics["MemTotal"]
	mem.availMemory = metrics["MemAvailable"]
	mem.freeMemory = metrics["MemFree"]
	if err := scanner.Err(); err != nil {
		return mem, err
	}

	/////////

	fmt.Println("not implemented for Linux")
	return mem, nil

}
