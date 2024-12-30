package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"os"
	"regexp"
	"strconv"
)

type set[T comparable] map[T]struct{}

func newSet[T comparable]() set[T] {
	return make(map[T]struct{}, 0)
}

func (s set[T]) insert(value T) {
	s[value] = struct{}{}
}

func (s set[T]) contains(value T) bool {
	_, ok := s[value]
	return ok
}

type coord struct {
	row int
	col int
}

var directions = map[byte]coord{
	'N': {-1, 0},
	'S': {+1, 0},
	'E': {0, +1},
	'W': {0, -1},
}

type memorySpace struct {
	size  int
	bytes []coord
}

var byteRegex = regexp.MustCompile(`(\d+),(\d+)`)

func newMemorySpace(size int, r io.Reader) (ms memorySpace) {
	ms.size = size

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		matches := byteRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			panic("invalid byte line")
		}

		col, err := strconv.Atoi(matches[1])
		if err != nil {
			panic("invalid byte col")
		}
		row, err := strconv.Atoi(matches[2])
		if err != nil {
			panic("invalid byte row")
		}

		ms.bytes = append(ms.bytes, coord{row, col})
	}

	return ms
}

func (ms memorySpace) inBounds(pos coord) bool {
	return pos.row >= 0 &&
		pos.col >= 0 &&
		pos.row <= ms.size &&
		pos.col <= ms.size
}

type elf struct {
	pos     coord
	time    int
	visited set[coord]
}

func (e elf) isBlocked(bytes []coord) bool {
	for _, b := range bytes[:e.time] {
		if b == e.pos {
			return true
		}
	}

	return false
}

func (e elf) nextElves(ms memorySpace) (elves []elf) {
	e.visited.insert(e.pos)

	for _, delta := range directions {
		next := e
		next.pos.row += delta.row
		next.pos.col += delta.col
		next.visited = maps.Clone(e.visited)
		if !e.visited.contains(next.pos) && ms.inBounds(next.pos) {
			elves = append(elves, next)
		}
	}

	return elves
}

func (ms memorySpace) minStepsToExit(time int) int {
	var open []elf
	visited := newSet[coord]()

	var start elf
	start.time = time
	start.visited = newSet[coord]()
	open = append(open, start)

	for len(open) != 0 {
		curr := open[0]
		open = open[1:]

		// -- Don't revisit positions.
		if visited.contains(curr.pos) {
			continue
		}
		visited.insert(curr.pos)

		// -- Skip blocked positions.
		if curr.isBlocked(ms.bytes) {
			continue
		}

		// -- Found exit.
		if curr.pos.row == ms.size && curr.pos.col == ms.size {
			return len(curr.visited)
		}

		// -- Try next positions.
		open = append(open, curr.nextElves(ms)...)
	}

	return -1
}

func main() {
	ms := newMemorySpace(70, os.Stdin)
	minSteps := ms.minStepsToExit(1024)
	fmt.Println(minSteps)
}
