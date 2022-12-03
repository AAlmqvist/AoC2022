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

// Turn string into uint64 where each bit in place 1-52 represents if badge is present
func sortBadges(badges string) uint64 {
	var sorted uint64
	for _, badge := range badges {
		val := uint64(1 << getValue(badge))
		if sorted&val == 0 {
			sorted += val
		}
	}
	return sorted
}

func splitRucksack(line string) (uint64, uint64) {
	return sortBadges(line[:len(line)/2]), sortBadges(line[len(line)/2:])
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
		// Sort the badges by flipping bit to 1 in places where badge is present
		first, second := splitRucksack(line)
		match := first & second
		i := 1
		for i < 64 {
			if match&(1<<i) > 0 {
				part1 += i
				break
			}
			i += 1
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
		match := sortBadges(group[0]) & sortBadges(group[1]) & sortBadges(group[2])
		i := 1
		for i < 64 {
			if match&(1<<i) > 0 {
				part2 += i
				break
			}
			i += 1
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
