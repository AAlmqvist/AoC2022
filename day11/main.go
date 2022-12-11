package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Map[T, U any](in []T, fn func(T) U) []U {
	res := make([]U, len(in))
	for i, e := range in {
		res[i] = fn(e)
	}
	return res
}

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
	return Filter(strings.Split(string(input), "\n\n"), func(s string) bool { return len(s) > 0 })
}

func getVal(in int, desc string) int {
	switch desc {
	case "old":
		return in
	default:
		nbr, _ := strconv.Atoi(desc)
		return nbr
	}
}

func makeOpFunc(indesc string) func(int) int {
	stuff := strings.Split(strings.Split(indesc, "= ")[1], " ")
	return func(in int) int {
		switch stuff[1] {
		case "+":
			return getVal(in, stuff[0]) + getVal(in, stuff[2])
		case "*":
			return getVal(in, stuff[0]) * getVal(in, stuff[2])
		}
		return 0
	}
}

func getLastInt(indesc string) int {
	inst := strings.Split(indesc, " ")
	lastInt, _ := strconv.Atoi(inst[len(inst)-1])
	return lastInt
}

type Prime struct {
	val   int
	count int
}

type Monkey struct {
	items     []int
	operation func(int) int
	divBy     int
	monkTrue  int
	monkFalse int
	inspected int
}

func NewMonkey(description string) *Monkey {
	lines := strings.Split(description, "\n")
	items := Map(strings.Split(strings.Split(lines[1], ": ")[1], ", "), func(in string) int {
		nbr, _ := strconv.Atoi(in)
		return nbr
	})
	opFunc := makeOpFunc(lines[2])
	divBy := getLastInt(lines[3])
	m1 := getLastInt(lines[4])
	m2 := getLastInt(lines[5])
	return &Monkey{
		items:     items,
		operation: opFunc,
		divBy:     divBy,
		monkTrue:  m1,
		monkFalse: m2,
	}
}

func (m *Monkey) takeTurn(other []*Monkey, extraWorried bool, mgn int) {
	for len(m.items) > 0 {
		item := m.items[0]
		if len(m.items) == 1 {
			m.items = []int{}
		} else {
			m.items = m.items[1:]
		}
		item = m.operation(item)
		if extraWorried {
			item = item % mgn
		} else {
			item = item / 3
		}
		if item < 0 {
			panic("Overflow!")
		}
		if item%m.divBy == 0 {
			other[m.monkTrue].recieve(item)
		} else {
			other[m.monkFalse].recieve(item)
		}
		m.inspected++
	}
}

func (m *Monkey) recieve(in int) {
	m.items = append(m.items, in)
}

func run(input []string, nbrRounds int, extraWorried bool) []*Monkey {
	monkeys := []*Monkey{}
	mgn := 1
	for _, monkDesc := range input {
		monk := NewMonkey(monkDesc)
		mgn *= monk.divBy
		monkeys = append(monkeys, monk)
	}
	for i := 0; i < nbrRounds; i++ {
		for _, monk := range monkeys {
			monk.takeTurn(monkeys, extraWorried, mgn)
		}
	}
	return monkeys
}

func monkeyBusiness(monkeys []*Monkey) int {
	var t1, t2 int
	for _, monk := range monkeys {
		if monk.inspected > t1 {
			t2 = t1
			t1 = monk.inspected
			continue
		}
		if monk.inspected > t2 {
			t2 = monk.inspected
		}
	}
	return t1 * t2
}

func main() {
	input := readInput("input.txt")
	monkeys := run(input, 20, false)
	fmt.Println(monkeyBusiness(monkeys))
	monkeys = run(input, 10000, true)
	fmt.Println(monkeyBusiness(monkeys))
}
