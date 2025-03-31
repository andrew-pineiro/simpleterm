package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func echo() {
	if len(c.args) == 0 {
		return
	}
	for i := range c.args {
		fmt.Printf("%s ", c.args[i])
	}
	fmt.Println()
}
func cp() {
	if len(c.args) < 2 {
		return
	}
	source := c.args[0]
	dest := c.args[1]
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
func cd() {
	if len(c.args) != 1 {
		fmt.Printf("ERROR: invalid arguments.")
		return
	}
	new_path, err := filepath.Abs(c.args[0])
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("ERROR: %s does not exist\n", new_path)
		return
	}
	// if new_path_str[:1] == "."+string(os.PathSeparator) ||
	// 	!strings.Contains(new_path_str, string(os.PathSeparator)) {
	// 	new_path_str = path.Join(c.working_dir, new_path_str)
	// }
	_, err = os.Stat(new_path)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("ERROR: %s does not exist\n", new_path)
		return
	}
	// if !new_path.IsDir() {
	// 	fmt.Printf("ERROR: %s is not a directory\n", new_path.Name())
	// }
	c.working_dir = new_path
	os.Chdir(new_path)
}
func sort_by_name_asc(entries []os.DirEntry) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
}
func ls() {
	var dirs []os.DirEntry
	var files []os.DirEntry
	f, err := os.ReadDir(c.working_dir)
	for _, file := range f {
		if file.IsDir() {
			dirs = append(dirs, file)
		} else {
			files = append(files, file)
		}
	}

	sort_by_name_asc(dirs)
	sort_by_name_asc(files)
	if err != nil {
		return
	}
	if len(files) <= 0 {
		return
	}

	fmt.Printf("\n  Directory: %s\n\n", c.working_dir)
	fmt.Printf("Mode\t\tModified\t\tSize\tName\n")
	fmt.Printf("----\t\t--------\t\t----\t----\n")
	for _, file := range dirs {
		file, _ := file.Info()
		size := " "
		name := file.Name()
		if file.Name()[0] == '.' {
			continue
		}
		fmt.Printf("%s\t%s\t%s\t\033[37;44;1m%s\033[0m\n", file.Mode(), file.ModTime().Format("1/2/2006 3:04 PM"), size, name)
	}
	for _, file := range files {
		file, _ := file.Info()
		size := strconv.FormatInt(file.Size(), 10)
		name := file.Name()

		if file.Name()[0] == '.' {
			continue
		}
		fmt.Printf("%s\t%s\t%s\t%s\n", file.Mode(), file.ModTime().Format("1/2/2006 3:04 PM"), size, name)
	}
	fmt.Println()
}
