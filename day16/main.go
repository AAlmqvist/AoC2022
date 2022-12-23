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

func mapValves(in []string) ([][]int, []int, map[string]int) {
	inds := make(map[string]int)
	flows := make([]int, len(in))
	conns := [][]int{}
	tmpConns := [][]string{}
	for i, li := range in {
		stf := strings.Split(li, ";")
		vInf := strings.Split(stf[0], " ")
		leadsTo := strings.Split(stf[1], ", ")
		leadsTo[0] = leadsTo[0][len(leadsTo[0])-2:]
		name := vInf[1]
		row := make([]int, len(in))
		conns = append(conns, row)
		rate, _ := strconv.Atoi(strings.Split(vInf[4], "=")[1])
		flows[i] = rate
		inds[name] = i
		tmpConns = append(tmpConns, leadsTo)
	}
	for i, neighbs := range tmpConns {
		for _, n := range neighbs {
			j := inds[n]
			conns[i][j] = 1
			conns[j][i] = 1
		}
	}
	return conns, flows, inds
}

func bfs(start int, in [][]int, done []bool) [][]int {
	n := len(in)
	visited := make([]bool, n)
	visited[start] = true
	next := []int{}
	for i, e := range in[start] {
		if e == 1 && !visited[i] {
			next = append(next, i)
			visited[i] = true
		}
	}
	steps := 1
	for len(next) > 0 {
		steps++
		after := []int{}
		for _, l := range next {
			neighs := []int{}
			for i, val := range in[l] {
				if i == start {
					continue
				}
				if val == 1 {
					if !visited[i] {
						neighs = append(neighs, i)
						in[start][i] = steps
						in[i][start] = steps
						visited[i] = true
					}
				}
			}
			after = append(after, neighs...)
		}
		next = after
	}
	return in
}

type Run struct {
	i       int
	visited []bool
	mins    int
	flow    int
	total   int
}

type dRun struct {
	i        []int
	visited  []bool
	nextStop []int
	hist     [2][]int
	mins     int
	flow     int
	total    int
}

type valve struct {
	i      int
	flow   int
	neighs []int
}

func copy[T any](a []T) []T {
	b := make([]T, len(a))
	for i, e := range a {
		b[i] = e
	}
	return b
}

func exploreOptions(startInd int, edges [][]int, flows []valve) int {
	m := 0
	start := Run{
		i:       startInd,
		visited: make([]bool, len(edges)),
	}
	start.visited[startInd] = true
	runs := []Run{start}
	for len(runs) > 0 {
		c := runs[0]
		at := c.i
		runs = runs[1:]
		for _, f := range flows {
			// Don't add valves visited
			if !c.visited[f.i] {
				steps := edges[at][f.i]
				// Don't add to path if there is no time left to reach
				if c.mins+steps >= 30 {
					continue
				}
				newR := Run{
					i: f.i,
					// Have to copy the values because the the slices points to
					// the same underlying array
					visited: copy(c.visited),
					mins:    c.mins,
					flow:    c.flow + f.flow,
					total:   c.total + (steps+1)*c.flow,
				}
				// add minutes for walking there and opening the valve
				newR.mins = newR.mins + edges[at][f.i] + 1
				// newR.order = append(newR.order, f.i)
				newR.visited[f.i] = true
				runs = append(runs, newR)
			}
		}
		tmp := c.total + (30-c.mins)*c.flow
		if tmp > m {
			m = tmp
		}
	}
	return m
}

func hash(a dRun) string {
	val := 0
	for i, v := range a.visited {
		if v {
			val += 1 << i
		}
	}
	return fmt.Sprintf("%d,%d,%d,%d,%d", val, a.i[0]+a.i[1], a.nextStop[0]+a.nextStop[1], a.mins, a.flow)
}

// len(a) <= len(b)
func isPrefix(a, b []int) bool {
	for i, ax := range a {
		if ax != b[i] {
			return false
		}
	}
	return true
}

func withHelp(startInd int, edges [][]int, flows []valve, inds map[int]int) int {
	m := 0
	// p1 := []int{startInd, 45, 26, 56, 8, 17, 54, 2}
	// p2 := []int{startInd, 50, 47, 14, 46, 7, 27}
	start := dRun{
		i:        []int{startInd, startInd},
		nextStop: []int{0, 0},
		visited:  make([]bool, len(edges)),
		hist:     [2][]int{[]int{startInd}, []int{startInd}},
	}
	seen := make(map[string]bool)
	start.visited[startInd] = true
	runs := []dRun{start}
	for len(runs) > 0 {
		c := runs[0]
		runs = runs[1:]
		p := 0
		o := 1
		if c.nextStop[1] < c.nextStop[0] {
			p = 1
			o = 0
		}
		// if isPrefix(c.hist[0], p1) && isPrefix(c.hist[1], p2) {
		// 	fmt.Println(" ---- from queue ----")
		// 	fmt.Println("(", c.i[p], ") ", c.i, c.nextStop)
		// 	fmt.Println("mins", c.mins)
		// 	fmt.Println("flow", c.flow)
		// 	fmt.Println("tot", c.total)
		// }
		passed := c.nextStop[p] - c.mins
		var adding int
		if c.i[0] != startInd && c.i[1] != startInd {
			c.mins = c.nextStop[p]
			c.total = c.total + c.flow*passed
			adding = flows[inds[c.i[p]]].flow
			c.flow = c.flow + adding
		}
		// if isPrefix(c.hist[0], p1) && isPrefix(c.hist[1], p2) {
		// 	fmt.Println(" --- after ---")
		// 	fmt.Println("mins", c.mins)
		// 	fmt.Println("flow", c.flow, "(added:", adding, ")")
		// 	fmt.Println("tot", c.total)
		// 	fmt.Println("Next", c.nextStop)
		// }
		for _, f := range flows {
			// Don't add valves visited
			if !c.visited[f.i] {
				steps := edges[c.i[p]][f.i]
				// Don't add to path if there is no time left to reach
				if c.mins+steps+1 >= 26 {
					continue
				}
				toAdd := c.mins
				newR := dRun{
					// Have to copy the values because the the slices points to
					// the same underlying array
					i:        copy(c.i),
					nextStop: copy(c.nextStop),
					visited:  copy(c.visited),
					hist:     [2][]int{copy(c.hist[0]), copy(c.hist[1])},
					mins:     toAdd,
					flow:     c.flow,
					total:    c.total,
				}
				newR.i[p] = f.i
				newR.nextStop[p] = toAdd + steps + 1
				newR.visited[f.i] = true
				newR.hist[p] = append(newR.hist[p], f.i)
				// if isPrefix(newR.hist[0], p1) && isPrefix(newR.hist[1], p2) {
				// 	fmt.Println(" --- New ---")
				// 	fmt.Println("next", newR.i)
				// 	fmt.Println("in", newR.nextStop)
				// 	fmt.Println(" - stats -")
				// 	fmt.Println("mins", newR.mins)
				// 	fmt.Println("flow", newR.flow)
				// 	fmt.Println("tot", newR.total)
				// 	fmt.Println("hist", newR.hist)
				// }
				h := hash(newR)
				if seen[h] {
					continue
				}
				seen[h] = true
				runs = append(runs, newR)
			}
		}
		if c.nextStop[o]-c.mins > 0 {
			passed = c.nextStop[o] - c.mins
			c.mins = c.nextStop[o]
			c.total = c.total + c.flow*passed
			c.flow = c.flow + flows[inds[c.i[o]]].flow
		}

		// if isPrefix(c.hist[0], p1) && isPrefix(c.hist[1], p2) {
		// 	fmt.Println(" --- final open ---")
		// 	fmt.Println("mins", c.mins)
		// 	fmt.Println("flow", c.flow)
		// 	fmt.Println("Added", flows[inds[c.i[o]]].flow)
		// 	fmt.Println("tot", c.total)
		// 	fmt.Println("Next", c.nextStop)
		// }
		tmp := c.total + (26-c.mins)*c.flow
		// if isPrefix(c.hist[0], p1) && isPrefix(c.hist[1], p2) {
		// 	fmt.Println(" --- final score ---")
		// 	fmt.Println("mins left", 26-c.mins)
		// 	fmt.Println("flow", c.flow)
		// 	fmt.Println("tot", tmp)
		// }
		if tmp > m {
			// if isPrefix(c.hist[0], p1) && isPrefix(c.hist[1], p2) {
			// 	fsum := 0
			// 	for _, e := range c.hist[0][1:] {
			// 		ff := flows[inds[e]].flow
			// 		fsum += ff
			// 	}
			// 	for _, e := range c.hist[1][1:] {
			// 		ff := flows[inds[e]].flow
			// 		fsum += ff
			// 	}
			// 	fmt.Println("Total flow: ", fsum)
			// }
			m = tmp
		}
	}
	// for _, ff := range flows {
	// 	fmt.Println(ff.i, ", ", ff.flow)
	// }
	// fmt.Println(inds)
	return m
}

func main() {
	input := readInput("input.txt")
	edges, flows, strMap := mapValves(input)
	start := time.Now()
	done := make([]bool, len(edges))
	for i := 0; i < len(edges); i++ {
		edges = bfs(i, edges, done)
		done[i] = true
	}
	fmt.Printf("Done mapping (%fs)\n", time.Since(start).Seconds())
	posFlows := []valve{}
	indMap := make(map[int]int)
	flo := []int{}
	for i, f := range flows {
		if f > 0 {
			posFlows = append(posFlows, valve{i: i, flow: f, neighs: edges[i]})
			indMap[i] = len(posFlows) - 1
			flo = append(flo, f)
		}
	}
	start = time.Now()
	part1 := exploreOptions(strMap["AA"], edges, posFlows)
	fmt.Printf("Done running part1 (%fs)\n", time.Since(start).Seconds())
	fmt.Println(part1)
	start = time.Now()
	part2 := withHelp(strMap["AA"], edges, posFlows, indMap)
	fmt.Printf("Done running part2 (%fs)\n", time.Since(start).Seconds())
	fmt.Println(part2)
}
