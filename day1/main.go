package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Println("")
		os.Exit(1)
	}
	elves := strings.Split(string(input), "\n")
	nsum := 0
	best := 0
	second := 0
	third := 0
	for _, elf := range elves {
		next, err := strconv.Atoi(elf)
		if err != nil {
			if nsum > third {
				third = nsum
			}
			if nsum > second {
				third = second
				second = nsum
			}
			if nsum > best {
				third = second
				second = best
				best = nsum
			}
			nsum = 0
		}
		nsum += next
	}
	fmt.Println(best)
	fmt.Println(best + second + third)
}
