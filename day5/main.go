package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Move struct {
	from int
	to   int
	nbr  int
}

type Stack struct {
	top *Crate
	nbr int
}

func (s *Stack) Place(crate *Crate) {
	if s.top == nil {
		s.top = crate
		s.nbr = 1
		return
	}
	crate.under = s.top
	s.top.over = crate
	s.top = crate
	s.nbr += 1
	return
}

func (s *Stack) Pick() *Crate {
	if s.top == nil {
		return nil
	}
	if s.nbr == 1 {
		toReturn := s.top
		s.top = nil
		s.nbr = 0
		return toReturn
	}
	toReturn := s.top
	s.top = s.top.under
	s.top.over = nil
	toReturn.under = nil
	s.nbr -= 1
	return toReturn
}

func (s *Stack) Print() {
	next := s.top
	if next == nil {
		fmt.Print("\n")
		return
	}
	for next != nil {
		fmt.Printf("%s ", next.value)
		next = next.under
	}
	fmt.Print("\n")
}

type Crate struct {
	value string
	over  *Crate
	under *Crate
}

func readInput(filename string) []string {
	input, _ := os.ReadFile(filename)
	return strings.Split(string(input), "\n\n")
}

func fillStacks(rows []string, nbrStacks int) []Stack {
	stacks := make([]Stack, nbrStacks)
	for row := len(rows) - 1; row > -1; row -= 1 {
		crates := rows[row]
		stack := 0
		for i := 1; i < len(crates); i += 4 {
			val := string(crates[i])
			if val != " " {
				crate := &Crate{value: val}
				stacks[stack].Place(crate)
			}
			stack += 1
		}
	}
	return stacks
}

func stackCrates(description string) []Stack {
	lines := strings.Split(description, "\n")
	// Apparently integer division by default. Nice
	nbrStacks := len(lines[0])/4 + 1
	stacks := fillStacks(lines[:len(lines)-1], nbrStacks)
	return stacks
}

func parseMoves(moves []string) []Move {
	parsedMoves := []Move{}
	for _, row := range moves {
		stuff := strings.Split(row, " ")
		if len(stuff) < 6 {
			continue
		}
		nbr, _ := strconv.Atoi(stuff[1])
		from, _ := strconv.Atoi(stuff[3])
		to, _ := strconv.Atoi(stuff[5])
		parsedMoves = append(parsedMoves, Move{from, to, nbr})
	}
	return parsedMoves
}

func MoveCrates9000(moves []Move, stacks []Stack) []Stack {
	for _, move := range moves {
		for i := 0; i < move.nbr; i++ {
			crate := stacks[move.from-1].Pick()
			stacks[move.to-1].Place(crate)
		}
	}
	return stacks
}

func MoveCrates9001(moves []Move, stacks []Stack) []Stack {
	for _, move := range moves {
		tempStack := Stack{}
		for i := 0; i < move.nbr; i++ {
			crate := stacks[move.from-1].Pick()
			tempStack.Place(crate)
		}
		for i := 0; i < move.nbr; i++ {
			crate := tempStack.Pick()
			stacks[move.to-1].Place(crate)
		}
	}
	return stacks
}

func printOutput(out string, stacks []Stack) {
	for _, stack := range stacks {
		fmt.Printf("%s", stack.top.value)
	}
	fmt.Print("\n")
}

func main() {
	input := readInput("input.txt")
	stacks := stackCrates(input[0])
	moves := parseMoves(strings.Split(input[1], "\n"))
	part1 := MoveCrates9000(moves, stacks)
	stacks = stackCrates(input[0])
	part2 := MoveCrates9001(moves, stacks)
	printOutput("part1: ", part1)
	printOutput("part2: ", part2)
}
