package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Filter[T any](in []T, fn func(T) bool) []T {
	res := make([]T, len(in))
	cnt := 0
	for i, e := range in {
		if fn(e) {
			res[i] = e
			cnt++
		}
	}
	return res[:cnt]
}

func readInput(filename string) []string {
	input, _ := os.ReadFile(filename)
	return Filter(strings.Split(string(input), "\n"), func(s string) bool { return len(s) > 0 })
}

func hash(x, y, z int) string {
	h := fmt.Sprintf("%d,%d,%d", x, y, z)
	// fmt.Println(h)
	return h
}

func makeCubes(in []string, ns [][]int) (map[string][]int, map[string][]int, []int, int) {
	cubes := make(map[string][]int)
	pockets := make(map[string][]int)
	match := 0
	maxVals := make([]int, 3)
	for _, e := range in {
		s := strings.Split(e, ",")
		c := []int{}
		for i, coo := range s {
			v, _ := strconv.Atoi(coo)
			if v > maxVals[i] {
				maxVals[i] = v
			}
			c = append(c, v)
		}
		for _, n := range ns {
			h := hash(c[0]+n[0], c[1]+n[1], c[2]+n[2])
			_, m := cubes[h]
			if m {
				match++
			}
		}
		cubes[e] = c
	}
	for _, c := range cubes {
		for _, n := range ns {
			x, y, z := c[0]+n[0], c[1]+n[1], c[2]+n[2]
			h := hash(x, y, z)
			_, hit := cubes[h]
			if !hit {
				pockets[h] = []int{x, y, z}
			}
		}
	}
	return cubes, pockets, maxVals, match
}

func makeGrid(gridSizes []int, filled map[string][]int) [][][]int {
	grid := [][][]int{}
	for i := 0; i < gridSizes[0]+2; i++ {
		r := [][]int{}
		for j := 0; j < gridSizes[1]+2; j++ {
			c := make([]int, gridSizes[2]+2)
			r = append(r, c)
		}
		grid = append(grid, r)
	}
	for _, pos := range filled {
		x, y, z := pos[0], pos[1], pos[2]
		grid[x][y][z] = 1
	}
	return grid
}

func isTrapped(p []int, ns [][]int, cubes map[string][]int, visited map[string]bool, maxVals []int) bool {
	for i, px := range p {
		if px == 0 || px > maxVals[i] {
			return false
		}
	}
	for _, n := range ns {
		x, y, z := p[0]+n[0], p[1]+n[1], p[2]+n[2]
		h := hash(x, y, z)
		_, found := visited[h]
		_, hit := cubes[h]
		if !found && !hit {
			visited[h] = true
			trapped := isTrapped([]int{x, y, z}, ns, cubes, visited, maxVals)
			if !trapped {
				return false
			}
		}
	}
	return true
}

func removePockets(sa int, ns [][]int, cubes, pockets map[string][]int, maxVals []int) (int, []string) {
	removed := []string{}
	for _, p := range pockets {
		visited := make(map[string]bool)
		h := hash(p[0], p[1], p[2])
		visited[h] = true
		rem := isTrapped(p, ns, cubes, visited, maxVals)
		na := 0
		if rem {
			removed = append(removed, h)
			for _, n := range ns {
				_, hit := cubes[hash(p[0]+n[0], p[1]+n[1], p[2]+n[2])]
				if hit {
					na++
				}
			}
		}
		sa -= na
	}
	return sa, removed
}

func findSurfaceArea(grid [][][]int, ns [][]int) (map[string]bool, int) {
	sa := 0
	for _, rc := range grid {
		lay := [][]int{}
		for _, r := range rc {
			row := make([]int, len(r))
			for i, v := range r {
				row[i] = v
			}
			lay = append(lay, row)
		}
	}
	visited := make(map[string]bool)
	start := []int{0, 0, 0}
	visited[hash(start[0], start[1], start[2])] = true
	toVisit := [][]int{}
	toVisit = append(toVisit, start)
	for len(toVisit) > 0 {
		c := toVisit[0]

		toVisit = toVisit[1:]
		for _, n := range ns {
			x, y, z := c[0]+n[0], c[1]+n[1], c[2]+n[2]
			// fmt.Printf("-- %d, %d, %d --\n", x, y, z)
			if x < 0 || x >= len(grid) {
				continue
			}
			if y < 0 || y >= len(grid[0]) {
				continue
			}
			if z < 0 || z >= len(grid[0][0]) {
				continue
			}
			switch grid[x][y][z] {
			case 0:
				h := hash(x, y, z)
				_, vis := visited[h]
				if !vis {
					visited[h] = true
					toVisit = append(toVisit, []int{x, y, z})
				}
			case 1:
				sa++
			}
		}
	}
	return visited, sa
}

func main() {
	input := readInput("input.txt")
	ns := [][]int{
		[]int{1, 0, 0},
		[]int{-1, 0, 0},
		[]int{0, 1, 0},
		[]int{0, -1, 0},
		[]int{0, 0, 1},
		[]int{0, 0, -1},
	}
	cubes, pockets, maxVals, matches := makeCubes(input, ns)
	part1 := 6*len(cubes) - 2*matches
	fmt.Println(part1)
	part2, removed := removePockets(part1, ns, cubes, pockets, maxVals)
	grid := makeGrid(maxVals, cubes)
	visited, _ := findSurfaceArea(grid, ns)
	for _, h := range removed {
		_, m := visited[h]
		if !m {
			fmt.Println(h)
		} else {
			fmt.Println("Overlap")
		}
	}
	fmt.Println(part2)
}
