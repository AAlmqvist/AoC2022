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

func createList(in []string) [][]int {
	numbers := make([][]int, len(in))
	for i, e := range in {
		val, _ := strconv.Atoi(e)
		numbers[i] = []int{i, val}
	}
	return numbers
}

func decode(nums [][]int, key, rounds int) int {
	p := [][]int{}
	for _, e := range nums {
		e[1] *= key
		p = append(p, e)
	}
	for l := 0; l < rounds; l++ {
		for i := 0; i < len(nums); i++ {
			j := indexOf(i, p)
			toMove := p[j]
			p = append(p[:j], p[j+1:]...)
			n := (j + toMove[1]) % len(p)
			if n < 0 {
				n += len(p)
			}
			p = append(p[:n], append([][]int{toMove}, p[n:]...)...)
		}
	}

	i0 := 0
	for i, e := range p {
		if e[1] == 0 {
			i0 = i
			break
		}
	}
	gc := 0
	for _, e := range []int{1000, 2000, 3000} {
		pos := (i0 + e) % len(p)
		if pos < 0 {
			pos += len(p)
		}
		gc += p[pos][1]
	}
	return gc
}

func indexOf(a int, b [][]int) int {
	for i, j := range b {
		if j[0] == a {
			return i
		}
	}
	return -1
}

func main() {
	input := readInput("input.txt")
	nums := createList(input)
	part1 := decode(nums, 1, 1)
	fmt.Println(part1)
	part2 := decode(nums, 811589153, 10)
	fmt.Println(part2)
}
