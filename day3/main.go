package main

import (
	"fmt"
	"os"
	"strings"
)

func readInput(filename string) []string {
	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		return []string{}
	}
	return strings.Split(string(input), "\n")
}

func splitRucksack(line string) (string, string) {
	first := line[:len(line)/2]
	second := line[len(line)/2:]
	return first, second
}

func getValue(in rune) int {
	if in < 'a' {
		return int(in-'A') + 27
	}
	return int(in-'a') + 1
}

func part1(input []string) int {
	part1 := 0
	for _, line := range input {
		first, second := splitRucksack(line)
	inner:
		for _, char1 := range first {
			for _, char2 := range second {
				if char1 == char2 {
					part1 += getValue(char1)
					break inner
				}
			}
		}
	}
	return part1
}

func part2(input []string) int {
	part2 := 0
	groups := [][]string{}
	for i := 0; i < len(input)-2; i += 3 {
		groups = append(groups, input[i:i+3])
	}
	for _, group := range groups {
	groupLoop:
		for _, char1 := range group[0] {
			for _, char2 := range group[1] {
				for _, char3 := range group[2] {
					if char1 == char2 && char1 == char3 {
						part2 += getValue(char1)
						break groupLoop
					}
				}
			}
		}
	}
	return part2
}

func main() {
	input := readInput("input.txt")
	part1 := part1(input)
	part2 := part2(input)
	fmt.Printf("Part1: %d\nPart2: %d", part1, part2)
}
