package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
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

func parseInput(in []string) map[string]string {
	equations := make(map[string]string)
	for _, e := range in {
		keyval := strings.Split(e, ": ")
		equations[keyval[0]] = keyval[1]
	}
	return equations
}

func solve(x string, equations map[string]string) int {
	eq, ok := equations[x]
	if !ok {
		panic(fmt.Sprintf("Can't find input %s in map", eq))
	}
	val, err := strconv.Atoi(eq)
	if err != nil {
		s := strings.Split(eq, " ")
		a := solve(s[0], equations)
		b := solve(s[2], equations)
		if x == "root" {
			fmt.Println(a, b)
		}
		switch s[1] {
		case "+":
			val = a + b
		case "-":
			val = a - b
		case "*":
			val = a * b
		case "/":
			val = a / b
		}
	}
	return val
}

type method string

const (
	add  method = "+"
	sub  method = "-"
	mult method = "*"
	div  method = "/"
)

type part interface {
	operate() int
	reduce(int, method) int
}

type value struct {
	val int
}

func (v *value) operate() int {
	return v.val
}

func (v *value) reduce(x int, m method) int {
	switch m {
	case add:
		return x - v.val
	case sub:
		return x + v.val
	case mult:
		return x / v.val
	case div:
		return x * v.val
	}
	return v.val
}

type operator struct {
	lhs  part
	rhs  part
	meth method
}

func (o *operator) operate() int {
	switch o.meth {
	case add:
		return o.lhs.operate() + o.rhs.operate()
	case sub:
		return o.lhs.operate() - o.rhs.operate()
	case mult:
		return o.lhs.operate() * o.rhs.operate()
	case div:
		return o.lhs.operate() / o.lhs.operate()
	}
	return -1
}

func (o *operator) reduce(in int) int {
	return 0
}

func mapVals(x string, equations map[string]string) (string, map[string]string) {
	if x == "humn" {
		return "humn", equations
	}
	eq, _ := equations[x]
	val, err := strconv.Atoi(eq)
	if err != nil {
		s := strings.Split(eq, " ")
		var a, b string
		a, equations = mapVals(s[0], equations)
		b, equations = mapVals(s[2], equations)
		if strings.Contains(a, "humn") || strings.Contains(b, "humn") {
			if x == "root" {
				text := fmt.Sprintf("%s = %s", a, b)
				equations[x] = text
				return text, equations
			}
			text := fmt.Sprintf("(%s %s %s)", a, s[1], b)
			equations[x] = text
			return text, equations
		}
		ai, _ := strconv.Atoi(a)
		bi, _ := strconv.Atoi(b)
		switch s[1] {
		case "+":
			val = ai + bi
		case "-":
			val = ai - bi
		case "*":
			val = ai * bi
		case "/":
			val = ai / bi
		}
		fmt.Println(ai, s[1], bi, "=", val)
	}
	return fmt.Sprintf("%d", val), equations
}

func solveEquation(x string, rhs int) int {
	// fmt.Println(x, "=", rhs)
	var lhs, c int
	var newX, op string
	if !strings.Contains(x, "(") {
		stuff := strings.Split(x, " ")
		if stuff[0] == "humn" {
			c, _ = strconv.Atoi(stuff[2])
		}
		if stuff[2] == "humn" {
			c, _ = strconv.Atoi(stuff[0])
		}
		switch stuff[1] {
		case "+":
			lhs = rhs - c
		case "-":
			lhs = rhs + c
		case "*":
			lhs = rhs / c
		case "/":
			lhs = rhs * c
		}
		// fmt.Println("humn =", lhs)
		return lhs
	}
	if x[0] == '(' {
		end := len(x) - 1
		for x[end] != ')' {
			end--
		}
		newX = x[1:end]
		s := strings.Split(x[end+2:], " ")
		op = s[0]
		c, _ = strconv.Atoi(s[1])
		switch op {
		case "+":
			lhs = solveEquation(newX, rhs-c)
		case "-":
			lhs = solveEquation(newX, rhs+c)
		case "*":
			lhs = solveEquation(newX, rhs/c)
		case "/":
			lhs = solveEquation(newX, rhs*c)
		}
	} else {
		ind := strings.Index(x, "(")
		newX = x[ind+1 : len(x)-1]
		s := strings.Split(x[:ind-1], " ")
		op = s[1]
		c, _ = strconv.Atoi(s[0])
		switch op {
		case "+":
			lhs = solveEquation(newX, rhs-c)
		case "-":
			lhs = solveEquation(newX, c-rhs)
		case "*":
			lhs = solveEquation(newX, rhs/c)
		case "/":
			lhs = solveEquation(newX, c/rhs)
		}
	}
	return lhs
}

func main() {
	input := readInput("input.txt")
	equations := parseInput(input)
	start := time.Now()
	part1 := solve("root", equations)
	fmt.Printf("Solved part1 in %fs\n", time.Since(start).Seconds())
	fmt.Println(part1)
	eq, equations := mapVals("root", equations)
	fmt.Println(eq)
	parts := strings.Split(eq, " = ")
	rhs, _ := strconv.Atoi(parts[1])
	lhs := solveEquation(parts[0][1:len(parts[0])-1], rhs)
	fmt.Println(lhs)
}
