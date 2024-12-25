package main

import (
	"bufio"
	"fmt"
	"io"
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

func (m maze) dijkstra() int {
	type node struct {
		pos coord
		dir direction
	}

	start := node{m.start, EAST}
	open := newSet[node](start)
	costs := make(map[node]int)
	costs[start] = 0

	// -- Find lowest cost open node.
	next := func() (n node) {
		lowest := math.MaxInt

		for o := range open {
			c, ok := costs[o]
			if !ok {
				continue
			}

			if c < lowest {
				n = o
				lowest = c
			}
		}

		return n
	}

	// -- Update costs.
	update := func(curr node, next node, cost int) {
		if m.isBlocked(next.pos) {
			return
		}

		currCost, ok := costs[curr]
		if !ok {
			panic("no cost for curr")
		}
		currCost += cost

		nextCost, ok := costs[next]
		if !ok {
			nextCost = math.MaxInt
		}

		if currCost < nextCost {
			costs[next] = currCost
			if !open.contains(next) {
				open.insert(next)
			}
		}
	}

	// -- Dijkstra's Algorithm.
	for len(open) != 0 {
		curr := next()
		delete(open, curr)

		// -- Return cost at goal.
		if curr.pos == m.end {
			for n, c := range costs {
				if n.pos == m.end {
					return c
				}
			}
			panic("no goal")
		}

		// -- Try moving forward and turning.
		fwd := curr
		fwd.pos.move(fwd.dir)
		update(curr, fwd, 1)

		for _, turn := range turns[curr.dir] {
			next := curr
			next.dir = turn
			update(curr, next, 1000)
		}
	}

	panic("no escape")
}

func main() {
	m := newMaze(os.Stdin)
	fmt.Println(m.dijkstra())
}
