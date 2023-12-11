package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const DEBUG = true

const (
	right int = iota
	down
	left
	up
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
	return Filter(strings.Split(string(input), "\n\n"), func(s string) bool { return len(s) > 0 })
}

func wrapAround(x, n int) int {
	xn := x % n
	if xn < 0 {
		xn += n
	}
	return xn
}

type Pos struct {
	x   int
	y   int
	dir int
}

func (p *Pos) Move(dx, dy, w, h int) {
	p.x = wrapAround(p.x+dx, w)
	p.y = wrapAround(p.y+dy, h)
}

func (p *Pos) MoveTo(o *Pos) {
	p.x = o.x
	p.y = o.y
}

func (p *Pos) Turn(rot int) {
	p.dir = wrapAround(p.dir+rot, 4)
}

func (p *Pos) Score() int {
	return (p.y+1)*1000 + (p.x+1)*4 + p.dir
}

func (p *Pos) copy() *Pos {
	nx := p.x
	ny := p.y
	ndir := p.dir
	return &Pos{nx, ny, ndir}
}

func findWidth(rows []string) int {
	w := 0
	for _, row := range rows {
		if len(row) > w {
			w = len(row)
		}
	}
	return w
}

func readMap(in string) ([][]int, *Pos) {
	var start *Pos
	coveMap := strings.Split(in, "\n")
	w := findWidth(coveMap)
	grid := [][]int{}
	for i, row := range coveMap {
		gRow := make([]int, w)
		j := 0
		for j < w {
			if j >= len(row) {
				gRow[j] = -1
				j++
				continue
			}
			switch row[j] {
			case ' ':
				gRow[j] = -1
			case '#':
				gRow[j] = 1
			case '.':
				if start == nil {
					start = &Pos{y: i, x: j}
				}
			}
			j++
		}
		grid = append(grid, gRow)
	}
	return grid, start
}

func parseMoves(in string) [][]int {
	moves := [][]int{}
	val := 0
	for _, e := range in {
		switch e {
		case 'R':
			b := []int{val, 1}
			moves = append(moves, b)
			val = 0
		case 'L':
			b := []int{val, -1}
			moves = append(moves, b)
			val = 0
		default:
			s, err := strconv.Atoi(string(e))
			if err != nil {
				continue
			}
			val *= 10
			val += s
		}
	}
	moves = append(moves, []int{val, 0})
	return moves
}

func exploreCove(start *Pos, moves, cove [][]int) *Pos {
	p := start.copy()
	h, w := len(cove), len(cove[0])
	for _, m := range moves {
		v, rot := m[0], m[1]
		var dx, dy int
		switch p.dir {
		case right:
			dx = 1
		case down:
			dy = 1
		case left:
			dx = -1
		case up:
			dy = -1
		}
		j := 0
		o := p.copy()
		for j < v && cove[o.y][o.x] != 1 {
			o.Move(dx, dy, w, h)
			for cove[o.y][o.x] == -1 {
				o.Move(dx, dy, w, h)
			}
			if cove[o.y][o.x] == 0 {
				p.MoveTo(o)
			}
			j++
		}
		p.Turn(rot)
	}
	return p
}

func main() {
	input := readInput("input.txt")
	cove, start := readMap(input[0])
	moves := parseMoves(input[1])
	final := exploreCove(start, moves, cove)
	fmt.Println(final.Score())
	cube := cubeFromMap(cove, start)
	print(cube.p)
}
