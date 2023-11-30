package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var DEBUG = false

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

type Blueprint struct {
	robots    [4]int
	resources [4]int
	costs     [4][4]int
}

func assembleBlueprints(in []string) []Blueprint {
	blueprints := make([]Blueprint, len(in))
	for i, line := range in {
		b := Blueprint{}
		b.robots[0] = 1
		stuff := strings.Split(line, " ")
		x, _ := strconv.Atoi(stuff[6])
		b.costs[0][0] = x
		x, _ = strconv.Atoi(stuff[12])
		b.costs[1][0] = x
		x, _ = strconv.Atoi(stuff[18])
		b.costs[2][0] = x
		x, _ = strconv.Atoi(stuff[21])
		b.costs[2][1] = x
		x, _ = strconv.Atoi(stuff[27])
		b.costs[3][0] = x
		x, _ = strconv.Atoi(stuff[30])
		b.costs[3][2] = x
		blueprints[i] = b
	}
	return blueprints
}

func (bp *Blueprint) timeToBuild(robot_idx, maxTime int) int {
	timeToAfford := 0
	for i, cx := range bp.costs[robot_idx] {
		// This resource is not needed to build this robot
		if cx == 0 {
			continue
		}
		// Resource needed but we have no robot to mine it
		if bp.robots[i] == 0 {
			return maxTime + 1
		}
		minsLeft := 0
		for cx > bp.resources[i]+bp.robots[i]*minsLeft {
			minsLeft++
		}
		if minsLeft > timeToAfford {
			timeToAfford = minsLeft
		}
	}
	return timeToAfford + 1
}

func (bp *Blueprint) maxRobotCost(resource int) int {
	maxCost := 0
	for i := range bp.costs {
		if bp.costs[i][resource] > maxCost {
			maxCost = bp.costs[i][resource]
		}
	}
	return maxCost
}

func (bp *Blueprint) build(i int) {
	for j, c := range bp.costs[i] {
		bp.resources[j] -= c
	}
	bp.robots[i] += 1
}

func (bp *Blueprint) deconstruct(i int) {
	for j, c := range bp.costs[i] {
		bp.resources[j] += c
	}
	bp.robots[i] -= 1
}

func (bp *Blueprint) mine(min int) {
	for i, cx := range bp.robots {
		bp.resources[i] += min * cx
	}
}

func factorial(i int) int {
	switch i {
	case 1:
		return 1
	case 2:
		return 2
	case 3:
		return 6
	case 4:
		return 24
	// If anything scores above this idk
	default:
		return 120
	}
}

func buildAndDig(min, maxTime int, score *int, bp *Blueprint) {
	if min >= maxTime {
		if bp.resources[3] > *score {
			*score = bp.resources[3]
		}
		return
	}
	timeLeft := maxTime - min
	if bp.resources[3]+timeLeft*bp.robots[3]+factorial(timeLeft) < *score {
		return
	}
	for i := 3; i >= 0; i-- {
		if i == 0 && bp.robots[i] >= bp.maxRobotCost(i) {
			continue
		}
		timeToBuild := bp.timeToBuild(i, maxTime)
		if timeToBuild > timeLeft {
			continue
		}
		bp.mine(timeToBuild)
		bp.build(i)
		buildAndDig(min+timeToBuild, maxTime, score, bp)
		bp.deconstruct(i)
		bp.mine(-timeToBuild)
	}
	bp.mine(timeLeft)
	buildAndDig(min+timeLeft, maxTime, score, bp)
	bp.mine(-timeLeft)
	return
}

func tryBlueprints(bps []Blueprint, maxTime int) int {
	p1 := 0
	for id, bp := range bps {
		if DEBUG {
			fmt.Printf("------------ Blueprint %d ------------\n", id+1)
			fmt.Println("Costs:")
			for i, costs := range bp.costs {
				fmt.Println(i, ":", costs)
			}
			fmt.Printf("maxCosts: %d, %d, %d, %d\n", bp.maxRobotCost(0), bp.maxRobotCost(1), bp.maxRobotCost(2), bp.maxRobotCost(3))
		}
		score := 0
		buildAndDig(0, maxTime, &score, &bp)
		if DEBUG {
			fmt.Println(bp.resources, bp.robots)
			fmt.Println(score)
		}
		p1 += (id + 1) * score
	}
	return p1
}

func tryBlueprints2(bps []Blueprint, maxTime int) int {
	p1 := 1
	for id, bp := range bps {
		if DEBUG {
			fmt.Printf("------------ Blueprint %d ------------\n", id+1)
		}
		score := 0
		buildAndDig(0, maxTime, &score, &bp)
		if DEBUG {
			fmt.Println(score)
		}
		p1 *= score
	}
	return p1
}

func main() {
	input := readInput("input.txt")
	blueprints := assembleBlueprints(input)
	part1 := tryBlueprints(blueprints, 24)
	fmt.Println(part1)
	if DEBUG {
		fmt.Println("-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-o-")
	}
	newBlueprints := blueprints
	if len(newBlueprints) > 3 {
		newBlueprints = newBlueprints[:3]
	}
	part2 := tryBlueprints2(newBlueprints, 32)
	fmt.Println(part2)
}
