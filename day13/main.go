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

func Map[T, U any](in []T, fn func(T) U) []U {
	res := make([]U, len(in))
	for i, e := range in {
		res[i] = fn(e)
	}
	return res
}

func readInput(filename string) [][]string {
	input, _ := os.ReadFile(filename)
	return Map(Filter(strings.Split(string(input), "\n\n"), func(s string) bool { return len(s) > 0 }),
		func(in string) []string { return strings.Split(in, "\n") })
}

type packet interface {
	vals() []packet
}

type packInt struct {
	val int
}

func (pi *packInt) vals() []packet {
	return nil
}

type packList struct {
	values []packet
}

func (pl *packList) vals() []packet {
	return pl.values
}

func parsePacket(in string) packet {
	packet, _ := parseInner(in, 1, len(in)-1)
	return packet
}

func parseInner(in string, start, end int) (packet, int) {
	i := start
	outer := packList{}
	for i < end {
		switch in[i] {
		case '[':
			stack := 1
			i2 := i + 1
			for stack != 0 {
				if in[i2] == ']' {
					stack--
				}
				if in[i2] == '[' {
					stack++
				}
				i2++
			}
			packin, newI := parseInner(in, i+1, i2-1)
			outer.values = append(outer.values, packin)
			i = newI
		case ']', ',':
			i++
		default:
			i2 := i + 1
			for i2 < end && in[i2] != ',' {
				i2++
			}
			newInt, err := strconv.Atoi(in[i:i2])
			if err != nil {
				fmt.Println("Error parsing ", in[i:i2], " to int")
			}
			outer.values = append(outer.values, &packInt{newInt})
			i = i2 + 1
		}
	}
	return &outer, i
}

func parsePairs(input [][]string) [][]packet {
	packetPairs := [][]packet{}
	for _, line := range input {
		pair := []packet{}
		pair = append(pair, parsePacket(line[0]))
		pair = append(pair, parsePacket(line[1]))
		packetPairs = append(packetPairs, pair)
	}
	return packetPairs
}

func parsePackets(input [][]string) []packet {
	packetPairs := []packet{}
	for _, line := range input {
		packetPairs = append(packetPairs, parsePacket(line[0]))
		packetPairs = append(packetPairs, parsePacket(line[1]))
	}
	return packetPairs
}

func comparePairs(pairs [][]packet) int {
	correct := 0
	for ind, pair := range pairs {
		right, final := compare(pair[0], pair[1])
		if !final {
			fmt.Println("Inconclusive")
			continue
		}
		if right {
			correct += ind + 1
		}
	}
	return correct
}

func compare(p1, p2 packet) (bool, bool) {
	switch c1 := p1.(type) {
	case *packInt:
		c2, ok := p2.(*packInt)
		if !ok {
			wrapped := packList{values: []packet{p1}}
			smaller, final := compare(&wrapped, p2)
			if final {
				return smaller, final
			}
		} else {
			if c1.val < c2.val {
				return true, true
			}
			if c2.val < c1.val {
				return false, true
			}
		}
	case *packList:
		c2, ok := p2.(*packList)
		if !ok {
			wrapped := packList{values: []packet{p2}}
			smaller, final := compare(p1, &wrapped)
			if final {
				return smaller, final
			}
		} else {
			i := 0
			packs1 := c1.vals()
			packs2 := c2.vals()
			for i < len(packs1) && i < len(packs2) {
				smaller, final := compare(packs1[i], packs2[i])
				if final {
					return smaller, final
				}
				i++
			}
			if i < len(packs2) {
				return true, true
			}
			if i < len(packs1) {
				return false, true
			}

		}
	}
	return false, false
}

func printP(p packet, layer int) {
	_, ok := p.(*packList)
	defer func() {
		if ok {
			fmt.Print("]")
		}
		if layer == 0 {
			fmt.Print("\n")
		}
	}()
	if ok {
		fmt.Print("[")
	}
	inner := p.vals()
	if inner == nil {
		p1, ok2 := p.(*packInt)
		if ok2 {
			fmt.Print(p1.val)
		}
		return
	}
	for i, pr := range p.vals() {
		printP(pr, layer+1)
		if i < (len(p.vals()) - 1) {
			fmt.Print(",")
		}
	}
}

func cmp(p1, p2 packet) bool {
	s, _ := compare(p1, p2)
	return s
}

func sort(in []packet) []packet {
	if len(in) < 12 {
		// smaller than 12 = do bubble sort
		swapped := true
		for swapped {
			swapped = false
			i := 0
			for i < len(in)-1 {
				if !cmp(in[i], in[i+1]) {
					in[i], in[i+1] = in[i+1], in[i]
					swapped = true
				}
				i++
			}
		}
		return in
	}
	out := []packet{}
	f1 := sort(in[:len(in)/2])
	f2 := sort(in[len(in)/2:])
	var i1, i2 int
	for i1 < len(f1) && i2 < len(f2) {
		if cmp(f1[i1], f2[i2]) {
			out = append(out, f1[i1])
			i1++
		} else {
			out = append(out, f2[i2])
			i2++
		}
	}
	if i1 < len(f1) {
		out = append(out, f1[i1:]...)
	}
	if i2 < len(f2) {
		out = append(out, f2[i2:]...)
	}
	return out
}

func findDecoderKey(ps []packet) int {
	i1 := 0
	findFirst := true
	first := parsePacket("[2]")
	second := parsePacket("[6]")
	for ind, pack := range ps {
		if findFirst {
			s1, _ := compare(first, pack)
			if s1 {
				i1 = ind + 1
				findFirst = false
			}
			continue
		}
		s1, _ := compare(second, pack)
		if s1 {
			return i1 * (ind + 2)
		}
	}
	return -1
}

func main() {
	input := readInput("input.txt")
	packetPairs := parsePairs(input)
	part1 := comparePairs(packetPairs)
	fmt.Println(part1)
	allPackets := parsePackets(input)
	sortedPackets := sort(allPackets)
	part2 := findDecoderKey(sortedPackets)
	fmt.Println(part2)
}
