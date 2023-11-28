package main

import (
	"fmt"
	"os"
)

func readInput(filename string) []int {
	input, _ := os.ReadFile(filename)
	winds := []int{}
	for _, run := range string(input) {
		switch run {
		case '<':
			winds = append(winds, -1)
		case '>':
			winds = append(winds, 1)
		}
	}
	return winds
}

type rock struct {
	anchor   []int
	relative [][]int
}

func (r *rock) print() {
	fmt.Println(r.anchor[0], r.anchor[1])
	for _, e := range r.relative {
		fmt.Println(r.anchor[0]+e[0], r.anchor[1]+e[1])
	}
}

func (r *rock) move(x, y int, in map[string][]int) bool {
	for _, b := range r.trans() {
		if b[0]+x < 0 || b[0]+x > 6 {
			return false
		}
		bs := toHash([]int{b[0] + x, b[1] + y})
		_, hit := in[bs]
		if hit {
			return false
		}
	}
	r.anchor[0] += x
	r.anchor[1] += y
	return true
}

func (r *rock) trans() [][]int {
	out := [][]int{r.anchor}
	for _, e := range r.relative {
		out = append(out, []int{e[0] + r.anchor[0], e[1] + r.anchor[1]})
	}
	return out
}

func createRel(in []int) [][]int {
	i := 0
	rels := [][]int{}
	for i+1 < len(in) {
		rels = append(rels, []int{in[i], in[i+1]})
		i += 2
	}
	return rels
}

func getRock(x, y, i int) *rock {
	switch i {
	// a###
	case 0:
		return &rock{
			anchor:   []int{x, y},
			relative: [][]int{{1, 0}, {2, 0}, {3, 0}},
		}
	//  #
	// a##
	//  #
	case 1:
		return &rock{
			anchor:   []int{x, y + 1},
			relative: [][]int{{1, 0}, {2, 0}, {1, 1}, {1, -1}},
		}
	//   #
	//   #
	// a##
	case 2:
		return &rock{
			anchor:   []int{x, y},
			relative: [][]int{{1, 0}, {2, 0}, {2, 1}, {2, 2}},
		}
	// #
	// #
	// #
	// a
	case 3:
		return &rock{
			anchor:   []int{x, y},
			relative: [][]int{{0, 1}, {0, 2}, {0, 3}},
		}
	// ##
	// a#
	case 4:
		return &rock{
			anchor:   []int{x, y},
			relative: [][]int{{1, 0}, {0, 1}, {1, 1}},
		}
	}
	return nil
}

func toHash(in []int) string {
	return fmt.Sprintf("%d,%d", in[0], in[1])
}

func printState(s int, state map[string][]int) {
	for s > 0 {
		fmt.Print("|")
		for l := 0; l < 7; l++ {
			ls := toHash([]int{l, s})
			_, in := state[ls]
			if in {
				fmt.Print("#")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("|\n")
		s--
	}
	fmt.Println("+-------+")
	fmt.Println(" ")
}

func fallingRocks(winds []int, settled map[string][]int) (map[string][]int, int, []int) {
	maxY := 0
	for _, e := range settled {
		if e[1] > maxY {
			maxY = e[1]
		}
	}
	t := 0 // type of rock next
	j := 0 // wind index
	i := 0
	heightIncreases := []int{}
	// Assume pattern is reached within 2022 iterations
	for i < 2022 {
		next := getRock(2, maxY+4, t)
		stuck := false
		for !stuck {
			dir := winds[j]
			next.move(dir, 0, settled)
			stuck = !next.move(0, -1, settled)
			j = (j + 1) % len(winds)
		}
		oldMaxY := maxY
		for _, e := range next.trans() {
			if e[1] > maxY {
				maxY = e[1]
			}
			es := toHash(e)
			settled[es] = e
		}
		t = (t + 1) % 5
		i++
		heightIncreases = append(heightIncreases, maxY-oldMaxY)
	}
	return settled, maxY, heightIncreases
}

func findPattern(data []int) (int, []int) {
	n := len(data)
	delta := 20
	// Assume we reached a reoccuring pattern (and that it is longer that delta)
	toMatch := data[n-delta:]
	k := n - delta - 1
	patternLength := -1
	for k > 0 {
		if compareLists(data[k:k+delta], toMatch) {
			patternLength = n - delta - k
			break
		}
		k--
	}
	pattern := data[n-patternLength : n]
	patternSum := 0
	for _, inc := range pattern {
		patternSum += inc
	}
	return patternSum, pattern
}

func compareLists(a, b []int) bool {
	if len(a) != len(b) {
		fmt.Println("Lists are not the same lengths")
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func computeFinalResult(part1, patternSum int, pattern []int) int {
	part2 := part1
	roundsLeft := 1000000000000 - 2022
	nbrPattern := roundsLeft / len(pattern)
	lastRound := roundsLeft % len(pattern)
	part2 += patternSum * nbrPattern
	for i, extra := range pattern {
		if i == lastRound {
			break
		}
		part2 += extra
	}
	return part2
}

func main() {
	winds := readInput("input.txt")
	floor := createRel([]int{0, 0, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0, 6, 0})
	settled := make(map[string][]int)
	for _, f := range floor {
		fs := toHash(f)
		settled[fs] = f
	}
	settled, part1, heightIncreases := fallingRocks(winds, settled)
	fmt.Println(part1)
	patternSum, pattern := findPattern(heightIncreases)
	part2 := computeFinalResult(part1, patternSum, pattern)
	fmt.Println(part2)
}
