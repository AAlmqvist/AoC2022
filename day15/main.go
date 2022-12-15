package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	realfile = "input.txt"
	testfile = "test.txt"
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

func Map[T, U any](in []T, fn func(T) U) []U {
	res := make([]U, len(in))
	for i, e := range in {
		res[i] = fn(e)
	}
	return res
}

func Any[T any](in []T, fn func(T) bool) bool {
	for _, e := range in {
		if fn(e) {
			return true
		}
	}
	return false
}

func readInput(filename string) []string {
	input, _ := os.ReadFile(filename)
	return Filter(strings.Split(string(input), "\n"), func(s string) bool { return len(s) > 0 })
}

type beacon struct {
	x       int
	y       int
	closest int
}

func abs(in int) int {
	if in < 0 {
		return -in
	}
	return in
}

func (b *beacon) dist(o *beacon) int {
	return abs(b.x-o.x) + abs(b.y-o.y)
}

func (b *beacon) findEndpoints(y int) (int, int, bool) {
	r := abs(b.y-y) - b.closest
	if r <= 0 {
		return b.x + r, b.x - r, true
	}
	return 0, 0, false
}

func (b *beacon) closer(o *beacon) bool {
	return b.dist(o) <= b.closest
}

func parseBeacons(in []string) ([]*beacon, []*beacon) {
	sensors := []*beacon{}
	beacons := []*beacon{}
	for _, line := range in {
		parts := strings.Split(line, ":")
		s1 := strings.Split(parts[0], " ")
		x, _ := strconv.Atoi(strings.Split(strings.TrimSuffix(s1[2], ","), "=")[1])
		y, _ := strconv.Atoi(strings.Split(s1[3], "=")[1])
		s2 := strings.Split(parts[1], " ")
		cx, _ := strconv.Atoi(strings.Split(strings.TrimSuffix(s2[5], ","), "=")[1])
		cy, _ := strconv.Atoi(strings.Split(s2[6], "=")[1])
		s := &beacon{x: x, y: y, closest: abs(x-cx) + abs(y-cy)}
		sensors = append(sensors, s)
		b := &beacon{x: cx, y: cy, closest: s.closest}
		if !Any(beacons, func(be *beacon) bool { return be.x == b.x && be.y == b.y }) {
			beacons = append(beacons, b)
		}
	}
	return sensors, beacons
}

func merge(a, b []int) []int {
	if a[0] > b[0] {
		a[0] = b[0]
	}
	if a[1] < b[1] {
		a[1] = b[1]
	}
	return a
}

func countRow(sensors, beacons []*beacon, yk int) int {
	ctr := 0
	ep := [][]int{}
	for _, s := range sensors {
		x1, x2, reached := s.findEndpoints(yk)
		if reached {
			ep = append(ep, []int{x1, x2})
		}
	}
	doMerge := true
	for doMerge {
		ep, doMerge = mergeEndpoints(ep)
	}
	for _, es := range ep {
		ctr += es[1] - es[0]
	}
	return ctr
}

func mergeEndpoints(ep [][]int) ([][]int, bool) {
	didMerge := false
	i := 0
	for i < len(ep)-1 {
		j := i + 1
		for j < len(ep) {
			if ep[i][0] <= ep[j][0] {
				if ep[i][1] >= ep[j][0] {
					ep[i] = merge(ep[i], ep[j])
					ep = append(ep[:j], ep[j+1:]...)
					didMerge = true
					continue
				}
			} else {
				if ep[i][0] <= ep[j][1] {
					ep[i] = merge(ep[i], ep[j])
					ep = append(ep[:j], ep[j+1:]...)
					didMerge = true
					continue
				}
			}
			j++
		}
		i++
	}
	return ep, didMerge
}

// L1-Ball of a sensor
type rombus struct {
	up    []int
	right []int
	down  []int
	left  []int
}

func newRombus(b *beacon) *rombus {
	return &rombus{
		up:    []int{b.x, b.y + b.closest},
		right: []int{b.x + b.closest, b.y},
		down:  []int{b.x, b.y - b.closest},
		left:  []int{b.x - b.closest, b.y},
	}
}

func findLines(rombs []*rombus, up bool) [][]int {
	lines := [][]int{}
	for i, r1 := range rombs {
		for j, r2 := range rombs {
			if i == j {
				continue
			}
			if up {
				if r2.up[0]+r2.up[1]-r1.down[0]-r1.down[1] < 0 {
					continue
				}
				if r1.right[0]+r1.right[1]-r2.left[0]-r2.left[1] < 0 {
					continue
				}
				a := 1
				b := -1
				c := -a*r1.down[0] - b*r1.down[1]
				if r2.left[0]-r2.left[1]+c == 2 {
					found := false
					for _, l := range lines {
						if l[2] == c-1 {
							found = true
						}
					}
					if !found {
						lines = append(lines, []int{a, b, c - 1})
					}
				}
			} else {
				if r1.right[0]-r1.right[1]-r2.left[0]+r2.left[1] < 0 {
					continue
				}
				if r2.down[0]-r2.down[1]-r1.up[0]+r1.up[1] < 0 {
					continue
				}
				a := 1
				b := 1
				c := -a*r1.up[0] - b*r1.up[1]
				if r2.left[0]+r2.left[1]+c == 2 {
					found := false
					for _, l := range lines {
						if l[2] == c-1 {
							found = true
						}
					}
					if !found {
						lines = append(lines, []int{a, b, c - 1})
					}
				}
			}
		}
	}
	return lines
}

func findTuningFrequency(sensors []*beacon, rombs []*rombus, maxVal int) int {
	lines1 := findLines(rombs, true)
	lines2 := findLines(rombs, false)
	for _, l1 := range lines1 {
		for _, l2 := range lines2 {
			// find intersection and check if its within the bounds
			xt := int(float64(l1[1]*l2[2]-l2[1]*l1[2]) / float64(l1[0]*l2[1]-l2[0]*l1[1]))
			yt := int(float64(l1[2]*l2[0]-l2[2]*l1[0]) / float64(l1[0]*l2[1]-l2[0]*l1[1]))
			if xt < 0 || xt > maxVal {
				continue
			}
			if yt < 0 || yt > maxVal {
				continue
			}
			bt := &beacon{xt, yt, 1}
			ok := true
			for _, s := range sensors {
				if s.closer(bt) {
					ok = false
				}
			}
			if ok {
				return xt*4000000 + yt
			}
		}
	}
	return -1
}

func main() {
	filename := realfile
	input := readInput(filename)
	sensors, beacons := parseBeacons(input)
	yt := 2000000
	p2 := 4000000
	if filename == testfile {
		yt = 10
		p2 = 20
	}
	part1 := countRow(sensors, beacons, yt)
	fmt.Println(part1)
	rombs := Map(sensors, newRombus)
	part2 := findTuningFrequency(sensors, rombs, p2)
	fmt.Println(part2)
}
