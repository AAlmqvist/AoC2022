package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Range struct {
	Start int
	Stop  int
}

func (r *Range) Contains(o *Range) bool {
	return r.Start <= o.Start && r.Stop >= o.Stop
}

func (r *Range) Overlaps(o *Range) bool {
	if r.Start <= o.Start && r.Stop >= o.Start {
		return true
	}
	if r.Start <= o.Stop && r.Stop >= o.Stop {
		return true
	}
	return false
}

func Map[T, U any](in []T, fn func(T) U) []U {
	res := make([]U, len(in))
	for i, e := range in {
		res[i] = fn(e)
	}
	return res
}

func Filter[T any](in []T, fn func(T) bool) []T {
	res := make([]T, len(in))
	c := 0
	for _, e := range in {
		if fn(e) {
			res[c] = e
			c += 1
		}
	}
	return res[:c]
}

func readInput(filename string) []string {
	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return Filter(strings.Split(string(input), "\n"), func(in string) bool { return len(in) > 0 })
}

func parseRange(ran string) Range {
	startStop := strings.Split(ran, "-")
	start, err := strconv.Atoi(startStop[0])
	if err != nil {
		fmt.Println("start", err)
	}
	stop, err := strconv.Atoi(startStop[1])
	if err != nil {
		fmt.Println("stop", err)
	}
	return Range{Start: start, Stop: stop}
}

func getRanges(input []string) [][]Range {
	groups := Map(input, func(line string) []Range {
		group := Map(strings.Split(line, ","), parseRange)
		return group
	})
	return groups
}

func countContains(groups [][]Range) int {
	num := 0
	_ = Map(groups, func(rans []Range) int {
		if rans[0].Contains(&rans[1]) {
			num += 1
		}
		if rans[1].Contains(&rans[0]) && !rans[0].Contains(&rans[1]) {
			num += 1
		}
		return 0
	})
	return num
}

func countOverlap(groups [][]Range) int {
	num := 0
	_ = Map(groups, func(rans []Range) int {
		if rans[0].Overlaps(&rans[1]) {
			num += 1
		}
		if rans[1].Overlaps(&rans[0]) && !rans[0].Overlaps(&rans[1]) {
			num += 1
		}
		return 0
	})
	return num
}

func main() {
	input := readInput("input.txt")
	groups := getRanges(input)
	part1 := countContains(groups)
	fmt.Println("part1: ", part1)
	part2 := countOverlap(groups)
	fmt.Println("part2: ", part2)
}
