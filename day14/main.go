package main

import (
	"fmt"
	"math"
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

func Map[T, U any](in []T, fn func(T) U) []U {
	res := make([]U, len(in))
	for i, e := range in {
		res[i] = fn(e)
	}
	return res
}

func readInput(filename string) []string {
	input, _ := os.ReadFile(filename)
	return Filter(strings.Split(string(input), "\n"), func(s string) bool { return len(s) > 0 })
}

type coords struct {
	x int
	y int
}

func (c *coords) Add(dx, dy int) {
	c.x += dx
	c.y += dy
}

func getPaths(input []string) ([][]coords, int, int, int) {
	startMax := 0
	xT := &startMax
	yT := &startMax
	tmp := math.MaxInt
	x0 := &tmp
	paths := Map(Map(input, func(in string) []string {
		return strings.Split(in, " -> ")
	}), func(in []string) []coords {
		return Map(in, func(pair string) coords {
			stuff := strings.Split(pair, ",")
			x, _ := strconv.Atoi(stuff[0])
			if x < *x0 {
				x0 = &x
			}
			if x > *xT {
				xT = &x
			}
			y, _ := strconv.Atoi(stuff[1])
			if y > *yT {
				yT = &y
			}
			return coords{x: x, y: y}
		})
	})
	return paths, *x0, *xT, *yT
}

func createCaverns(x0, xT, yT int, part2 bool) [][]int {
	caverns := [][]int{}
	if part2 {
		yT = yT + 1
	}
	for i := 0; i < yT; i++ {
		if part2 {
			caverns = append(caverns, make([]int, 2*yT+1))
			continue
		}
		caverns = append(caverns, make([]int, xT-x0+1))
	}
	if part2 {
		caverns = append(caverns, Map(make([]int, 2*yT+1), func(i int) int { return 1 }))
	}
	return caverns
}

func abs(in int) int {
	if in < 0 {
		return -in
	}
	return in
}

func sign(in int) int {
	if in == 0 {
		return 0
	}
	return in / abs(in)
}

func markCaverns(cavs [][]int, paths [][]coords, x0 int) [][]int {
	for _, path := range paths {
		start := &coords{path[0].x, path[0].y}
		cavs[start.y][start.x-x0] = 1
		i := 0
		for i < len(path) {
			dx := sign(path[i].x - start.x)
			dy := sign(path[i].y - start.y)
			for start.x != path[i].x || start.y != path[i].y {
				start.Add(dx, dy)
				cavs[start.y][start.x-x0] = 1
			}
			i++
		}
	}
	return cavs
}

func fillSand(cavs [][]int, x0 int) ([][]int, int) {
	notFull := true
	nbrSand := 0
	for notFull {
		start := &coords{500, 0}
		atRest := false
		for !atRest {
			if start.y >= len(cavs)-1 {
				return cavs, nbrSand
			}
			if cavs[start.y+1][start.x-x0] != 0 {
				if start.x-1-x0 == -1 {
					return cavs, nbrSand
				}
				if cavs[start.y+1][start.x-1-x0] == 0 {
					start.Add(-1, 0)
				} else {
					if start.x+1-x0 == len(cavs[0]) {
						return cavs, nbrSand
					}
					if cavs[start.y+1][start.x+1-x0] != 0 {
						cavs[start.y][start.x-x0] = 2
						atRest = true
						if start.x == 500 && start.y == 0 {
							nbrSand++
							return cavs, nbrSand
						}
					} else {
						start.Add(1, 0)
					}
				}
			}
			start.Add(0, 1)
		}
		nbrSand++
	}
	return cavs, nbrSand
}

func main() {
	input := readInput("input.txt")
	paths, x0, xT, yT := getPaths(input)
	caverns := createCaverns(x0, xT+1, yT+1, false)
	caverns = markCaverns(caverns, paths, x0)
	caverns, part1 := fillSand(caverns, x0)
	fmt.Println(part1)
	cavs2 := createCaverns(0, xT, yT+1, true)
	cavs2 = markCaverns(cavs2, paths, 500-len(cavs2[0])/2)
	cavs2, part2 := fillSand(cavs2, 500-len(cavs2[0])/2)
	fmt.Println(part2)
}
