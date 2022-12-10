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

func tick(reg, clock, val, score *int, row string, rows *[]string) string {
	if *clock%40 == 20 {
		*score += *clock * *reg
	}
	pos := (*clock - 1) % 40
	pixel := " "
	if pos >= *reg-1 && pos < *reg+2 {
		pixel = "#"
	}
	row += pixel
	if pos == 39 {
		*rows = append(*rows, row)
		row = ""
	}
	*clock++
	if val != nil {
		*reg += *val
	}
	return row
}

func run(instructions []string) ([]string, int) {
	rows := &[]string{}
	x := 1
	c := 1
	score := 0
	row := ""
	for _, line := range instructions {
		op := strings.Split(line, " ")
		row = tick(&x, &c, nil, &score, row, rows)
		if len(op) > 1 {
			val, _ := strconv.Atoi(op[1])
			row = tick(&x, &c, &val, &score, row, rows)
		}
	}
	return *rows, score
}

func main() {
	input := readInput("input.txt")
	prompt, part1 := run(input)
	fmt.Println(part1)
	for _, line := range prompt {
		fmt.Println(line)
	}
}
