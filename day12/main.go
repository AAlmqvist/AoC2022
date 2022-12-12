package main

import (
	"fmt"
	"math"
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

type pQueue struct {
	Next *pQueue
	x    int
	y    int
	val  int
}

func Add(p, new *pQueue) *pQueue {
	if p == nil {
		return new
	}
	if new.val < p.val {
		new.Next = p
		return new
	}
	curr := p
	for curr.Next != nil {
		if new.val < curr.val {
			curr.Next, new.Next = new, curr.Next
			return p
		}
		curr = curr.Next
	}
	curr.Next = new
	return p
}

func makeGrid(input []string) ([][]rune, int, int) {
	var startx, starty int
	grid := [][]rune{}
	for x, line := range input {
		row := []rune{}
		for y, char := range line {
			if char == 'S' {
				startx, starty = x, y
			}
			row = append(row, char)
		}
		grid = append(grid, row)
	}
	return grid, startx, starty
}

func getHeight(in rune) int {
	switch in {
	case 'S':
		return 0
	case 'E':
		return 'z' - 'a'
	default:
		return int(in - 'a')
	}
}

func neigh(x, y, nx, ny, val int, grid [][]rune, visited [][]bool, reverse bool) *pQueue {
	if nx < 0 || nx >= len(grid) {
		return nil
	}
	if ny < 0 || ny >= len(grid[x]) {
		return nil
	}
	if visited[nx][ny] {
		return nil
	}
	if reverse {
		if getHeight(grid[x][y])-getHeight(grid[nx][ny]) > 1 {
			return nil
		}

	} else {
		if getHeight(grid[nx][ny])-getHeight(grid[x][y]) > 1 {
			return nil
		}
	}
	return &pQueue{x: nx, y: ny, val: val}
}

func findShortestPath(sx, sy int, grid [][]rune, reverse bool) int {
	visited := [][]bool{}
	for _, row := range grid {
		visited = append(visited, make([]bool, len(row)))
	}
	pq := &pQueue{x: sx, y: sy, val: 0}
	for pq != nil {
		n := pq
		pq = n.Next
		if reverse {
			if getHeight(grid[n.x][n.y]) == 0 {
				return n.val
			}
		} else {
			if grid[n.x][n.y] == 'E' {
				return n.val
			}
		}
		if !visited[n.x][n.y] {
			visited[n.x][n.y] = true
			new := neigh(n.x, n.y, n.x-1, n.y, n.val+1, grid, visited, reverse)
			if new != nil {
				pq = Add(pq, new)
			}
			new = neigh(n.x, n.y, n.x+1, n.y, n.val+1, grid, visited, reverse)
			if new != nil {
				pq = Add(pq, new)
			}
			new = neigh(n.x, n.y, n.x, n.y-1, n.val+1, grid, visited, reverse)
			if new != nil {
				pq = Add(pq, new)
			}
			new = neigh(n.x, n.y, n.x, n.y+1, n.val+1, grid, visited, reverse)
			if new != nil {
				pq = Add(pq, new)
			}
		}
	}
	return math.MaxInt
}

func main() {
	input := readInput("input.txt")
	grid, sx, sy := makeGrid(input)
	part1 := findShortestPath(sx, sy, grid, false)
	fmt.Println(part1)
	var ex, ey int
	for x, row := range grid {
		for y, char := range row {
			if char == 'E' {
				ex, ey = x, y
			}
		}
	}
	part2 := findShortestPath(ex, ey, grid, true)
	fmt.Println(part2)
}
