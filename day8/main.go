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

func Any(bs []bool) bool {
	for _, b := range bs {
		if b {
			return true
		}
	}
	return false
}

func visibleUp(tree, i, j int, grid [][]int) (int, bool) {
	see := 0
	for k := i - 1; k > -1; k-- {
		see += 1
		if tree <= grid[k][j] {
			return see, false
		}
	}
	return see, true
}

func visibleRight(tree, i, j int, grid [][]int) (int, bool) {
	see := 0
	for k := j + 1; k < len(grid[0]); k++ {
		see += 1
		if tree <= grid[i][k] {
			return see, false
		}
	}
	return see, true
}

func visibleDown(tree, i, j int, grid [][]int) (int, bool) {
	see := 0
	for k := i + 1; k < len(grid); k++ {
		see += 1
		if tree <= grid[k][j] {
			return see, false
		}
	}
	return see, true
}

func visibleLeft(tree, i, j int, grid [][]int) (int, bool) {
	see := 0
	for k := j - 1; k > -1; k-- {
		see += 1
		if tree <= grid[i][k] {
			return see, false
		}
	}
	return see, true
}

func readInput(filename string) []string {
	input, _ := os.ReadFile(filename)
	return Filter(strings.Split(string(input), "\n"), func(s string) bool { return len(s) > 0 })
}

func makeGrid(in []string) [][]int {
	grid := [][]int{}
	for _, line := range in {
		row := []int{}
		for _, char := range line {
			tree, _ := strconv.Atoi(string(char))
			row = append(row, tree)
		}
		grid = append(grid, row)
	}
	return grid
}

func keep(tree, i, j int, grid [][]int) bool {
	_, b1 := visibleUp(tree, i, j, grid)
	_, b2 := visibleRight(tree, i, j, grid)
	_, b3 := visibleDown(tree, i, j, grid)
	_, b4 := visibleLeft(tree, i, j, grid)
	test := []bool{b1, b2, b3, b4}
	return Any(test)
}

func keepVisible(grid [][]int) [][]int {
	kept := [][]int{}
	for _, row := range grid {
		kept = append(kept, make([]int, len(row)))
	}
	for i, row := range grid {
		for j, tree := range row {
			if keep(tree, i, j, grid) {
				kept[i][j] = 1
			}
		}
	}
	return kept
}

func getBestScenicScore(grid [][]int) int {
	best := 0
	for i, row := range grid {
		for j, tree := range row {
			s1, _ := visibleUp(tree, i, j, grid)
			s2, _ := visibleLeft(tree, i, j, grid)
			s3, _ := visibleDown(tree, i, j, grid)
			s4, _ := visibleRight(tree, i, j, grid)
			scenicScore := s1 * s2 * s3 * s4
			if best < scenicScore {
				best = scenicScore
			}
		}
	}
	return best
}

func main() {
	input := readInput("input.txt")
	grid := makeGrid(input)
	kept := keepVisible(grid)
	part1 := 0
	for _, row := range kept {
		for _, tree := range row {
			part1 += tree
		}
	}
	part2 := getBestScenicScore(grid)
	fmt.Println(part1)
	fmt.Println(part2)
}
