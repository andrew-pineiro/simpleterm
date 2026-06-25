package main

import (
	"fmt"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func getDrives() ([]string, error) {
	n, err := unix.Getfsstat(nil, unix.MNT_NOWAIT)
	if err != nil {
		return nil, err
	}

	stats := make([]unix.Statfs_t, n)
	n, err = unix.Getfsstat(stats, unix.MNT_NOWAIT)
	if err != nil {
		return nil, err
	}

	var drives []string
	for _, stat := range stats[:n] {
		dev := bytesToString(stat.Mntfromname[:])
		if len(dev) > 0 {
			drives = append(drives, dev)
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
	return false
}

func getDiskSpaceAvailable(drives []string) []stDisk {
	var disks []stDisk
	for _, v := range drives {
		var stat unix.Statfs_t
		unix.Statfs(v, &stat)

		bsize := uint64(stat.Bsize)
		total := stat.Blocks * bsize
		avail := stat.Bavail * bsize

		var used uint64
		if total > avail {
			used = total - avail
		}
		disks = append(disks, stDisk{
			driveName:  v,
			availSpace: avail,
			totalSpace: total,
			usedSpace:  used,
		})
	}
	return disks
}
func getMemoryInfo() (stMem, error) {
	var mem stMem
	fmt.Println("not implemented for MacOS")
	return mem, nil
}
