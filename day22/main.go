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
	return Filter(strings.Split(string(input), "\n\n"), func(s string) bool { return len(s) > 0 })
}

type Pos struct {
	x   int
	y   int
	dir int
}

func (p *Pos) Move(dx, dy, w, h int) {
	p.x = (p.x + dx) % w
	if p.x < 0 {
		p.x = p.x + w
	}
	p.y = (p.y + dy) % h
	if p.y < 0 {
		p.y = p.y + h
	}
}

func (p *Pos) MoveTo(o *Pos) {
	p.x = o.x
	p.y = o.y
}

func (p *Pos) Turn(rot int) {
	p.dir = (p.dir + rot) % 4
	if p.dir < 0 {
		p.dir += 4
	}
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
			val *= 10
			s, _ := strconv.Atoi(string(e))
			val += s
		}
	}
	moves = append(moves, []int{val, 0})
	return moves
}

func exploreCove(p *Pos, moves, cove [][]int) (*Pos, [][]int) {
	h, w := len(cove), len(cove[0])
	fmt.Println(p.y, p.x)
	for _, m := range moves {
		v, rot := m[0], m[1]
		var dx, dy int
		switch p.dir {
		case 0:
			dx = 1
		case 1:
			dy = 1
		case 2:
			dx = -1
		case 3:
			dy = -1
		}
		j := 0
		fmt.Println(p.dir)
		o := p.copy()
		for j < v && cove[o.y][o.x] != 1 {
			o.Move(dx, dy, w, h)
			for cove[o.y][o.x] == -1 {
				o.Move(dx, dy, w, h)
			}
			// if cove[o.y][o.x] != 1 && cove[o.y][o.x] != -1 {
			if cove[o.y][o.x] == 0 {
				p.MoveTo(o)
				cove[p.y][p.x] = 2 + o.dir
			}
			j++
		}
		p.Turn(rot)
		cove[o.y][o.x] = 2 + p.dir
		// fmt.Println(p.y, p.x)
	}
	return p, cove
}

func main() {
	input := readInput("test.txt")
	cove, start := readMap(input[0])
	for _, e := range cove {
		for _, b := range e {
			switch b {
			case -1:
				fmt.Print("  ")
			default:
				fmt.Printf("%d ", b)
			}
		}
		fmt.Print("\n")
	}
	fmt.Println("Starting coordinates", start.y, start.x)
	// moves := parseMoves(input[1])
	moves := parseMoves("4R2L3R3R2L1R7R1L6R4R6R2L5R2L2R4")
	final, cove2 := exploreCove(start, moves, cove)
	for _, e := range cove2 {
		for _, b := range e {
			switch b {
			case -1:
				fmt.Print(" ")
			case 0:
				fmt.Print(".")
			case 1:
				fmt.Print("#")
			case 2:
				fmt.Print(">")
			case 3:
				fmt.Print("v")
			case 4:
				fmt.Print("<")
			case 5:
				fmt.Print("^")
			}
		}
		fmt.Print("\n")
	}
	fmt.Println(final.Score())
}
