package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const DEBUG = false

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

func placeElfs(in []string) map[string]elf {
	pos := make(map[string]elf)
	for i, row := range in {
		for j, v := range row {
			if v == '#' {
				e := elf{x: j, y: -i}
				pos[e.hash(0, 0)] = e
			}
		}
	}
	return pos
}

type elf struct {
	x int
	y int
}

func (e *elf) hash(dx, dy int) string {
	return fmt.Sprintf("%d,%d", e.x+dx, e.y+dy)
}

func dehash(h string) elf {
	s := strings.Split(h, ",")
	x, _ := strconv.Atoi(s[0])
	y, _ := strconv.Atoi(s[1])
	return elf{x: x, y: y}
}

func proposeMove(e elf, pos map[string]elf, dirs []int) string {
	ds := []int{-1, 0, 1}
	noNeighbors := true
	for _, dx := range ds {
		for _, dy := range ds {
			if dx == 0 && dy == 0 {
				continue
			}
			_, ok := pos[e.hash(dx, dy)]
			if ok {
				noNeighbors = false
			}
		}
	}
	// Don't move if we stand alone
	if noNeighbors {
		return e.hash(0, 0)
	}
	for _, d := range dirs {
		canMove := true
		switch d {
		// north
		case 0:
			for _, dx := range ds {
				_, occupied := pos[e.hash(dx, 1)]
				if occupied {
					canMove = false
					break
				}
			}
			if canMove {
				return e.hash(0, 1)
			}
		// south
		case 1:
			for _, dx := range ds {
				_, occupied := pos[e.hash(dx, -1)]
				if occupied {
					canMove = false
					break
				}
			}
			if canMove {
				return e.hash(0, -1)
			}
		// west
		case 2:
			for _, dy := range ds {
				_, occupied := pos[e.hash(-1, dy)]
				if occupied {
					canMove = false
					break
				}
			}
			if canMove {
				return e.hash(-1, 0)
			}
		// east
		case 3:
			for _, dy := range ds {
				_, occupied := pos[e.hash(1, dy)]
				if occupied {
					canMove = false
					break
				}
			}
			if canMove {
				return e.hash(1, 0)
			}
		}
	}
	return e.hash(0, 0)
}

func update(pos map[string]elf, dirs []int) (map[string]elf, []int, bool) {
	newPos := make(map[string]elf)
	proposed := make(map[string][]string)
	for _, e := range pos {
		h := proposeMove(e, pos, dirs)
		proposed[h] = append(proposed[h], e.hash(0, 0))
	}
	for h, m := range proposed {
		if len(m) > 1 {
			for _, oldH := range m {
				newPos[oldH] = dehash(oldH)
			}
			continue
		}
		newPos[h] = dehash(h)
	}
	moved := false
	for k := range newPos {
		_, ok := pos[k]
		if !ok {
			moved = true
		}
	}
	dirs = append(dirs[1:], dirs[0])
	if len(pos) != len(newPos) {
		fmt.Println("------ NUMBER OF ELVES ARE CHANGING -------")
		fmt.Printf("%d != %d\n", len(pos), len(newPos))
	}
	return newPos, dirs, moved
}

func findEdges(pos map[string]elf) (int, int, int, int) {
	xMin := math.MaxInt
	yMin := math.MaxInt
	xMax := math.MinInt
	yMax := math.MinInt
	for _, e := range pos {
		if e.x < xMin {
			xMin = e.x
		}
		if e.x > xMax {
			xMax = e.x
		}
		if e.y < yMin {
			yMin = e.y
		}
		if e.y > yMax {
			yMax = e.y
		}
	}
	return xMin, xMax, yMin, yMax
}

func printState(pos map[string]elf) {
	x0, xM, y0, yM := findEdges(pos)
	for l := yM; l > y0-1; l-- {
		for m := x0; m < xM+1; m++ {
			_, ok := pos[fmt.Sprintf("%d,%d", m, l)]
			if ok {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Println()
}

func main() {
	input := readInput("input.txt")
	elves := placeElfs(input)
	dirs := []int{0, 1, 2, 3}
	if DEBUG {
		fmt.Println("INITIAL")
		printState(elves)
	}
	for i := 0; i < 10; i++ {
		elves, dirs, _ = update(elves, dirs)
		if DEBUG {
			fmt.Println("\nROUND", i+1)
			printState(elves)
		}
	}
	x0, xM, y0, yM := findEdges(elves)
	fmt.Println((yM-y0+1)*(xM-x0+1) - len(elves))
	// Assume round of no moves did not occur within 10 rounds
	moved := true
	round := 10
	for moved == true {
		elves, dirs, moved = update(elves, dirs)
		round++
		if DEBUG {
			fmt.Println("\nROUND", round)
			printState(elves)
		}
	}
	fmt.Println(round)
}
