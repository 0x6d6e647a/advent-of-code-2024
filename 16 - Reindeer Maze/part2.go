package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
)

type set[T comparable] map[T]struct{}

func newSet[T comparable](values ...T) set[T] {
	s := make(set[T])
	s.insert(values...)
	return s
}

func (s set[T]) insert(values ...T) {
	for _, v := range values {
		s[v] = struct{}{}
	}
}

func (s set[T]) contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s set[T]) clone() set[T] {
	return maps.Clone(s)
}

type coord struct {
	row int
	col int
}

type direction uint8

const (
	NORTH direction = iota
	SOUTH
	EAST
	WEST
)

var deltas = map[direction]coord{
	NORTH: {-1, 0},
	SOUTH: {+1, 0},
	EAST:  {0, +1},
	WEST:  {0, -1},
}

var turns = map[direction][2]direction{
	NORTH: {WEST, EAST},
	SOUTH: {EAST, WEST},
	EAST:  {NORTH, SOUTH},
	WEST:  {SOUTH, NORTH},
}

func (c *coord) move(dir direction) {
	delta := deltas[dir]
	c.row += delta.row
	c.col += delta.col
}

type maze struct {
	walls set[coord]
	start coord
	end   coord
}

func newMaze(r io.Reader) (m maze) {
	m.walls = newSet[coord]()

	row := 0
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		for col, ch := range line {
			switch ch {
			case '.':
				continue
			case '#':
				m.walls.insert(coord{row, col})
			case 'S':
				m.start = coord{row, col}
			case 'E':
				m.end = coord{row, col}
			}
		}

		row += 1
	}

	return m
}

func (m maze) isBlocked(pos coord) bool {
	return m.walls.contains(pos)
}

type deer struct {
	pos     coord
	dir     direction
	score   int
	visited set[coord]
}

func (d *deer) fwd() {
	d.visited.insert(d.pos)
	d.pos.move(d.dir)
	d.score += 1
}

func (m maze) countSeats() int {
	// -- Find all paths to the goal.
	var open []deer
	var closed []deer
	scores := make(map[coord]int)

	start := deer{m.start, EAST, 0, newSet[coord]()}
	open = append(open, start)

	for len(open) != 0 {
		curr := open[0]
		open = open[1:]

		for !m.isBlocked(curr.pos) {
			// -- Check and update position score.
			posScore, ok := scores[curr.pos]
			if !ok {
				posScore = math.MaxInt
			}
			if posScore < curr.score-1000 {
				break
			}
			if posScore > curr.score {
				scores[curr.pos] = curr.score
			}

			// -- Add left and right turns.
			for _, dir := range turns[curr.dir] {
				next := curr
				next.dir = dir
				next.score += 1000
				next.visited = curr.visited.clone()
				open = append(open, next)
			}

			// -- Check if at goal.
			if curr.pos == m.end {
				closed = append(closed, curr)
				break
			}

			curr.fwd()
		}
	}

	// -- Count seats.
	seats := newSet[coord]()

	for _, d := range closed {
		if d.score == scores[m.end] {
			for p := range d.visited {
				seats.insert(p)
			}
		}
	}

	return len(seats) + 1
}

func main() {
	m := newMaze(os.Stdin)
	fmt.Println(m.countSeats())
}
