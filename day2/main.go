package main

import (
	"fmt"
	"os"
	"strings"
)

// Values for RPC-enum, should had beed a real enum but that wasn't fast enough
const (
	rock    = 1
	paper   = 2
	scissor = 3
)

// Parse encoded input into rock/paper/scissor manual enum
func GetRPC(char string) int {
	switch char {
	case "A", "X":
		return rock
	case "B", "Y":
		return paper
	case "C", "Z":
		return scissor
	default:
		return 0
	}
}

// Get the correctly chosen move against the opponent
func GetRPC2(char string, op int) int {
	switch char {
	// lose
	case "X":
		if op == 1 {
			return 3
		}
		return op - 1
	// play equal
	case "Y":
		return op
	// win
	case "Z":
		return op%3 + 1
	default:
		return 0
	}
}

// Calculate the score for the outcome of the round
// (without value for chosen rock/paper/scissor)
func RPC(op, me int) int {
	if op == me%3+1 {
		return 0
	}
	if op == me {
		return 3
	}
	return 6
}

func ReadInput(filename string) []string {
	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
	}
	return strings.Split(string(input), "\n")
}

func main() {
	lines := ReadInput("input.txt")
	score1 := 0
	score2 := 0
	for _, line := range lines {
		var op, me, me2 int
		stuff := strings.Split(line, " ")
		if len(stuff) > 1 {
			op = GetRPC(stuff[0])
			me = GetRPC(stuff[1])
			me2 = GetRPC2(stuff[1], op)
			score1 += me + RPC(op, me)
			score2 += me2 + RPC(op, me2)
		}
	}
	fmt.Println(score1)
	fmt.Println(score2)
}
