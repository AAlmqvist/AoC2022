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

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sign(n int) int {
	if n == 0 {
		return 0
	}
	return n / abs(n)
}

func readInput(filename string) []string {
	input, _ := os.ReadFile(filename)
	return Filter(strings.Split(string(input), "\n"), func(s string) bool { return len(s) > 0 })
}

type End struct {
	x int
	y int
}

func (e *End) Add(dx, dy int) {
	e.x += dx
	e.y += dy
}

func (e *End) Follow(head *End) bool {
	dx := head.x - e.x
	dy := head.y - e.y
	if abs(dx) > 1 || abs(dy) > 1 {
		e.x += sign(dx)
		e.y += sign(dy)
		return true
	}
	return false
}

func (e *End) String() string {
	return fmt.Sprintf("%d,%d", e.x, e.y)
}

type Rope struct {
	Parts []*End
}

func NewRope(nbr int) *Rope {
	r := Rope{}
	for i := 0; i < nbr; i++ {
		r.Parts = append(r.Parts, &End{})
	}
	return &r
}

func (r *Rope) Move(dir string) bool {

	switch dir {
	case "U":
		r.Parts[0].Add(0, 1)
	case "R":
		r.Parts[0].Add(1, 0)
	case "D":
		r.Parts[0].Add(0, -1)
	case "L":
		r.Parts[0].Add(-1, 0)
	}
	i := 1
	for i < len(r.Parts) {
		if !r.Parts[i].Follow(r.Parts[i-1]) {
			return false
		}
		i++
	}
	return true
}

func Move(r *Rope, moves []string) int {
	visited := []string{}
	for _, move := range moves {
		stuff := strings.Split(move, " ")
		nbr, _ := strconv.Atoi(stuff[1])
		for i := 0; i < nbr; i++ {
			r.Move(stuff[0])
			tailPos := r.Parts[len(r.Parts)-1].String()
			toAdd := true
			for _, pos := range visited {
				if pos == tailPos {
					toAdd = false
					break
				}
			}
			if toAdd {
				visited = append(visited, tailPos)
			}
		}
	}
	return len(visited)
}

func main() {
	input := readInput("input.txt")
	fmt.Println(Move(NewRope(2), input))
	fmt.Println(Move(NewRope(10), input))
}
