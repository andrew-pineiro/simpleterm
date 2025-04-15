package main

import (
	"path/filepath"
	"sort"
	"strings"
)

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
