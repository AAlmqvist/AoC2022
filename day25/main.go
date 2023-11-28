package main

import (
	"fmt"
	"os"
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

func convertToDecimal(in []string) []int {
	dec := make([]int, len(in))
	for i, snafu := range in {
		nbr := 0
		for _, s := range snafu {
			nbr = nbr * 5
			switch s {
			case '2':
				nbr = nbr + 2
			case '1':
				nbr = nbr + 1
			case '-':
				nbr = nbr - 1
			case '=':
				nbr = nbr - 2
			}
			dec[i] = nbr
		}
	}
	return dec
}

func sum(a []int) int {
	b := 0
	for _, n := range a {
		b += n
	}
	return b
}

// Makes more sense than you first thought
func decToSnafu(a int) string {
	fmt.Println(a)
	powers := []int{}
	for a > 0 {
		powers = append(powers, a%5)
		a = a / 5
	}
	snafu := ""
	for i := range powers {
		if powers[i] > 2 {
			powers[i+1] += 1
			powers[i] -= 5
		}
		switch powers[i] {
		case 2:
			snafu = "2" + snafu
		case 1:
			snafu = "1" + snafu
		case 0:
			snafu = "0" + snafu
		case -1:
			snafu = "-" + snafu
		case -2:
			snafu = "=" + snafu
		}
	}
	return snafu
}

func main() {
	input := readInput("input.txt")
	decNbrs := convertToDecimal(input)
	result := sum(decNbrs)
	convRes := decToSnafu(result)
	fmt.Println(convRes)
}
