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

type Blueprint struct {
	robots    [4]int
	resources [4]int
	costs     [4][4]int
	visited   map[int]bool
}

func assembleBlueprints(in []string) []Blueprint {
	blueprints := make([]Blueprint, len(in))
	for i, line := range in {
		b := Blueprint{visited: make(map[int]bool)}
		b.robots[0] = 1
		stuff := strings.Split(line, " ")
		x, _ := strconv.Atoi(stuff[6])
		b.costs[0][0] = x
		x, _ = strconv.Atoi(stuff[12])
		b.costs[1][0] = x
		x, _ = strconv.Atoi(stuff[18])
		b.costs[2][0] = x
		x, _ = strconv.Atoi(stuff[21])
		b.costs[2][1] = x
		x, _ = strconv.Atoi(stuff[27])
		b.costs[3][0] = x
		x, _ = strconv.Atoi(stuff[30])
		b.costs[3][2] = x
		blueprints[i] = b
	}
	return blueprints
}

func (bp *Blueprint) afford(c int) bool {
	for i, cx := range bp.costs[c] {
		if bp.resources[i] < cx {
			return false
		}
	}
	return true
}

func (bp *Blueprint) build(i int) {
	for j, c := range bp.costs[i] {
		bp.resources[j] -= c
	}
	bp.robots[i] += 1
}

func (bp *Blueprint) deconstruct(i int) {
	for j, c := range bp.costs[i] {
		bp.resources[j] += c
	}
	bp.robots[i] -= 1
}

func (bp *Blueprint) unmine() {
	for i, cx := range bp.robots {
		bp.resources[i] -= cx
	}
}

func (bp *Blueprint) mine() {
	for i, cx := range bp.robots {
		bp.resources[i] += cx
	}
}

func (bp *Blueprint) seen(m int) bool {
	_, ok := bp.visited[m]
	if !ok {
		bp.visited[m] = true
		return false
	}
	return true
}

func (bp *Blueprint) max(i int) int {
	m := 0
	for _, v := range bp.costs {
		if v[i] > m {
			m = v[i]
		}
	}
	return m
}

func factorial(i int) int {
	switch i {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 6
	case 4:
		return 24
	default:
		return 30
	}
}

func buildAndDig(min, maxTime int, earliest, score *int, bp *Blueprint, best map[int]int) map[int]int {
	if min >= maxTime {
		if bp.resources[3] > *score {
			*score = bp.resources[3]
		}
		return best
	}
	if bp.resources[3] < best[min] {
		return best
	}
	if bp.robots[0] == 0 && min > *earliest {
		return best
	}
	left := maxTime - min
	if bp.resources[3]+bp.robots[3]*left+factorial(left-1) < *score {
		return best
	}
	if bp.robots[3] == 0 && bp.resources[2]+bp.robots[2]*(left-2)+factorial(left-3) < bp.costs[3][2] {
		return best
	}
	for i := 3; i >= 0; i-- {
		if bp.afford(i) {
			if i == 3 && min < *earliest {
				*earliest = min
			}
			if i == 3 && bp.resources[3] < best[min] {
				best[min] = bp.resources[3]
			}
			bp.mine()
			bp.build(i)
			best = buildAndDig(min+1, maxTime, earliest, score, bp, best)
			bp.deconstruct(i)
			bp.unmine()
			if i == 3 {
				return best
			}
		}
	}
	bp.mine()
	best = buildAndDig(min+1, maxTime, earliest, score, bp, best)
	bp.unmine()
	return best
}

func tryBlueprints(bps []Blueprint, maxTime int) int {
	p1 := 0
	for id, bp := range bps {
		fmt.Printf("------------ Blueprint %d ------------\n", id+1)
		score := 0
		earliest := 0
		bestAtTime := make(map[int]int)
		buildAndDig(0, maxTime, &earliest, &score, &bp, bestAtTime)
		fmt.Println(score)
		p1 += (id + 1) * score
	}
	return p1
}

func tryBlueprints2(bps []Blueprint, maxTime int) int {
	p1 := 1
	for id, bp := range bps {
		fmt.Printf("------------ Blueprint %d ------------\n", id+1)
		score := 0
		earliest := 0
		bestAtTime := make(map[int]int)
		buildAndDig(0, maxTime, &earliest, &score, &bp, bestAtTime)
		fmt.Println(score)
		p1 *= score
	}
	return p1
}

func main() {
	input := readInput("input.txt")
	blueprints := assembleBlueprints(input)
	part1 := tryBlueprints(blueprints, 24)
	fmt.Println(part1)
	fmt.Println("-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-")
	part2 := tryBlueprints2(blueprints[:3], 32)
	fmt.Println(part2)
}
