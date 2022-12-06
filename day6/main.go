package main

import (
	"fmt"
	"os"
)

func readInput(filename string) string {
	input, _ := os.ReadFile(filename)
	return string(input)
}

func Unique(sequence string) bool {
	unique := uint32(0)
	for _, char := range sequence {
		val := uint32(1 << (char - 'a'))
		if unique&val > 0 {
			return false
		}
		unique += val
	}
	return true
}

func findFirstAllUniquePacket(buffer string, packetLen int) int {
	ind := packetLen
	sequence := buffer[ind-packetLen : ind]
	for !Unique(sequence) && ind < len(buffer) {
		ind++
		sequence = buffer[ind-packetLen : ind]
	}
	if ind >= len(buffer) {
		return -1
	}
	return ind
}

func main() {
	input := readInput("input.txt")
	part1 := findFirstAllUniquePacket(input, 4)
	fmt.Println("Part2: ", part1)
	part2 := findFirstAllUniquePacket(input, 14)
	fmt.Println("Part2: ", part2)
}
