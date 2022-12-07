package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type InFolder interface {
	Name() string
}

type folder struct {
	name     string
	size     int
	parent   *folder
	children []InFolder
}

func (f *folder) Name() string {
	return f.name
}

func (f *folder) Get(sf string) *folder {
	for _, child := range f.children {
		switch c := child.(type) {
		case *folder:
			if c.Name() == sf {
				return c
			}
		}
	}
	return nil
}

func (f *folder) Add(in InFolder) {
	f.children = append(f.children, in)
}

type file struct {
	name string
	size int
}

func (f *file) Name() string {
	return f.name
}

func readInput(filename string) []string {
	input, _ := os.ReadFile(filename)
	toRet := make([]string, len(input))
	cutOff := 0
	for i, e := range strings.Split(string(input), "\n") {
		if len(e) == 0 {
			continue
		}
		toRet[i] = e
		cutOff++
	}
	return toRet[:cutOff]
}

func execCmd(cwd *folder, root *folder, cmds []string) *folder {
	cmd := strings.Split(cmds[0], " ")
	switch cmd[1] {
	case "cd":
		switch cmd[2] {
		case "/":
			return root
		case "..":
			return cwd.parent
		default:
			newWd := cwd.Get(cmd[2])
			if newWd == nil {
				fmt.Println(cmd)
				fmt.Println("newWd is nil, things will break")
			}
			return newWd
		}

	case "ls":
		for _, line := range cmds[1:] {
			stuff := strings.Split(line, " ")
			size, err := strconv.Atoi(stuff[0])
			if err != nil {
				// subfolder
				fold := &folder{name: stuff[1], parent: cwd}
				cwd.Add(fold)
				continue
			}
			newFile := &file{name: stuff[1], size: size}
			cwd.Add(newFile)
		}
		return cwd
	}
	return nil
}

func buildDir(cmds []string) *folder {
	root := folder{name: "/"}
	wd := &root
	ind := 0
	for ind < len(cmds) {
		if cmds[ind][0] == '$' {
			next := ind + 1
			for next < len(cmds) && cmds[next][0] != '$' {
				next++
			}
			res := execCmd(wd, &root, cmds[ind:next])
			if res != nil {
				wd = res
			}
			ind = next
			continue
		}
		ind++
	}
	return &root
}

func CountSize(f *folder, total *int) int {
	folder_size := 0
	for _, child := range f.children {
		switch c := child.(type) {
		case *folder:
			folder_size += CountSize(c, total)
		case *file:
			folder_size += c.size
		}
	}
	f.size = folder_size
	if folder_size < 100000 {
		*total += folder_size
	}
	return folder_size
}

func FindOptimalSize(minimal int, f *folder, optimal *int) {
	if f.size < *optimal {
		*optimal = f.size
	}
	for _, child := range f.children {
		switch c := child.(type) {
		case *folder:
			if c.size > minimal {
				FindOptimalSize(minimal, c, optimal)
			}
		}
	}
}

// Visualizes the File tree in terminal
// (for debugging purposes)
func PrintDirs(wd *folder, layer int) {
	printLayer(fmt.Sprintf("- %s (dir)", wd.Name()), layer)
	for _, child := range wd.children {
		switch c := child.(type) {
		case *folder:
			PrintDirs(c, layer+1)
		case *file:
			printLayer(fmt.Sprintf("- %s (file, size=%d)", c.name, c.size), layer+1)
		}
	}
	return
}

// helper method to visualisation
func printLayer(toPrint string, layer int) {
	fmt.Println(strings.Repeat("  ", layer), toPrint)
}

func main() {
	cmds := readInput("input.txt")
	root := buildDir(cmds)
	// PrintDirs(root, 0)
	part1 := 0
	usedMem := CountSize(root, &part1)
	fmt.Println(part1)
	totalMem := 70000000
	minimalDelete := 30000000 - (totalMem - usedMem)
	part2 := totalMem
	FindOptimalSize(minimalDelete, root, &part2)
	fmt.Println(part2)
}
