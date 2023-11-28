package main

import (
	"fmt"
	"os"
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

func createMap(in []string) [][]int {
	grid := [][]int{}
	for _, row := range in {
		r := make([]int, len(row))
		for i, v := range row {
			switch v {
			case '#':
				r[i] = 16
			case '>':
				r[i] = 1
			case 'v':
				r[i] = 2
			case '<':
				r[i] = 4
			case '^':
				r[i] = 8
			}
		}
		grid = append(grid, r)
	}
	return grid
}

func updateGrid(grid [][]int) [][]int {
	newGrid := [][]int{}
	for _, row := range grid {
		r := make([]int, len(row))
		for i, v := range row {
			if v == 16 {
				r[i] = 16
			}
		}
		newGrid = append(newGrid, r)
	}
	for i := 1; i < len(grid)-1; i++ {
		for j := 1; j < len(grid[0])-1; j++ {
			if grid[i][j]&1 > 0 {
				if j == len(grid[0])-2 {
					newGrid[i][1] += 1
				} else {
					newGrid[i][j+1] += 1
				}
			}
			if grid[i][j]&2 > 0 {
				if i == len(grid)-2 {
					newGrid[1][j] += 2
				} else {
					newGrid[i+1][j] += 2
				}
			}
			if grid[i][j]&4 > 0 {
				if j == 1 {
					newGrid[i][len(grid[0])-2] += 4
				} else {
					newGrid[i][j-1] += 4
				}
			}
			if grid[i][j]&8 > 0 {
				if i == 1 {
					newGrid[len(grid)-2][j] += 8
				} else {
					newGrid[i-1][j] += 8
				}
			}
		}
	}
	return newGrid
}

type queue struct {
	first *node
	last  *node
	size  int
}

func (q *queue) Add(e *node) {
	if q.last == nil {
		q.first = e
		q.last = e
		q.size = 1
		return
	}
	q.last.next = e
	q.last = e
	q.size += 1
}

func (q *queue) Pop() *node {
	m := q.first
	if m.next == nil {
		q.first = nil
		q.last = nil
		q.size = 0
		return m
	}
	q.first = m.next
	q.size -= 1
	return m
}

type node struct {
	r    int
	c    int
	min  int
	next *node
}

func hash(c *node) string {
	return fmt.Sprintf("%d,%d,%d", c.r, c.c, c.min)
}

func findShortestPath(time int, start, end []int, grid [][]int) (int, [][]int) {
	q := &queue{}
	sn := &node{r: start[0], c: start[1], min: time}
	visited := make(map[string]bool)
	q.Add(sn)
	min := -1
	for q.size > 0 {
		c := q.Pop()
		if c.r == end[0] && c.c == end[1] {
			fmt.Println("reached goal at", c.r, c.c)
			return c.min, grid
		}
		h := hash(c)
		_, ok := visited[h]
		if ok {
			continue
		}
		visited[h] = true
		if c.min > min {
			grid = updateGrid(grid)
			min = c.min
		}
		if c.r > 0 && grid[c.r][c.c] == 0 {
			cn := &node{r: c.r - 1, c: c.c, min: c.min + 1}
			q.Add(cn)
		}
		if c.r < len(grid)-1 && grid[c.r+1][c.c] == 0 {
			cn := &node{r: c.r + 1, c: c.c, min: c.min + 1}
			q.Add(cn)
		}
		if grid[c.r][c.c-1] == 0 {
			cn := &node{r: c.r, c: c.c - 1, min: c.min + 1}
			q.Add(cn)
		}
		if grid[c.r][c.c+1] == 0 {
			cn := &node{r: c.r, c: c.c + 1, min: c.min + 1}
			q.Add(cn)
		}
		if grid[c.r][c.c] != 0 {
			continue
		}
		c.min += 1
		q.Add(c)
	}
	return 0, grid
}

func printGrid(grid [][]int) {
	for _, r := range grid {
		for j, v := range r {
			if v < 10 || j == len(r)-1 {
				fmt.Print(" ")
			}
			fmt.Print(v)
		}
		fmt.Print("\n")
	}
}

func main() {
	input := readInput("test.txt")
	grid := createMap(input)
	start := []int{0}
	for i, v := range grid[0] {
		if v == 0 {
			start = append(start, i)
			break
		}
	}
	end := []int{len(grid) - 1}
	for i, v := range grid[len(grid)-1] {
		if v == 0 {
			end = append(end, i)
			break
		}
	}
	sp, grid := findShortestPath(0, start, end, grid)
	fmt.Println(sp)
	sp, grid = findShortestPath(sp, end, start, grid)
	fmt.Println(sp)
	sp, grid = findShortestPath(sp, start, end, grid)
	fmt.Println(sp)
}
