package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
)

type set[T comparable] map[T]struct{}

func newSet[T comparable]() set[T] {
	return make(map[T]struct{}, 0)
}

func (s set[T]) insert(value T) {
	s[value] = struct{}{}
}

type direction uint8

const (
	NORTH direction = iota
	EAST
	SOUTH
	WEST
)

func (d direction) turnRight() direction {
	switch d {
	case NORTH:
		return EAST
	case EAST:
		return SOUTH
	case SOUTH:
		return WEST
	case WEST:
		return NORTH
	default:
		panic("invalid direction")
	}
}

type coord struct {
	row int
	col int
}

func (c coord) move(dir direction) coord {
	switch dir {
	case NORTH:
		c.row -= 1
	case EAST:
		c.col += 1
	case SOUTH:
		c.row += 1
	case WEST:
		c.col -= 1
	}

	return c
}

type guard struct {
	pos coord
	dir direction
}

type labMap struct {
	numRows   int
	numCols   int
	obstacles []coord
	guard     guard
}

func newLabMap(r io.Reader) labMap {
	row := 0
	col := -1
	var obstacles []coord
	guard := guard{coord{-1, -1}, SOUTH}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		for col, char := range line {
			switch char {
			case '#':
				obstacles = append(obstacles, coord{row, col})
			case '^':
				guard.pos.row = row
				guard.pos.col = col
				guard.dir = NORTH
			}
		}

		row += 1
		col = len(line)
	}

	return labMap{row, col, obstacles, guard}
}

func (l labMap) inBounds(pos coord) bool {
	return pos.row >= 0 && pos.col >= 0 && pos.row < l.numRows && pos.col < l.numCols
}

func (l labMap) isBlocked(pos coord) bool {
	return slices.Contains(l.obstacles, pos)
}

func (l labMap) walkGuard() []guard {
	var guardWalk []guard

	for {
		guardWalk = append(guardWalk, l.guard)

		// -- Move forward avoiding obstacles.
		var newPos coord

		for {
			newPos = l.guard.pos.move(l.guard.dir)
			if !l.isBlocked(newPos) {
				break
			}
			l.guard.dir = l.guard.dir.turnRight()
		}

		// -- Done if out of bounds.
		if !l.inBounds(newPos) {
			return guardWalk
		}

		l.guard.pos = newPos
	}
}

func main() {
	labMap := newLabMap(os.Stdin)
	guardWalk := labMap.walkGuard()

	guardPosSet := newSet[coord]()
	for _, guard := range guardWalk {
		guardPosSet.insert(guard.pos)
	}

	fmt.Println(len(guardPosSet))
}
