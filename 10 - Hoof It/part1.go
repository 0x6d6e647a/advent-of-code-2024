package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

type direction uint8

const (
	NORTH direction = iota
	EAST
	SOUTH
	WEST
)

var directions = [...]direction{NORTH, EAST, SOUTH, WEST}

func (d direction) getDelta() coord {
	switch d {
	case NORTH:
		return coord{-1, 0}
	case EAST:
		return coord{0, 1}
	case SOUTH:
		return coord{1, 0}
	case WEST:
		return coord{0, -1}
	default:
		panic("invalid direction")
	}
}

type coord struct {
	row int
	col int
}

func (c coord) next(dir direction) coord {
	delta := dir.getDelta()
	c.row += delta.row
	c.col += delta.col
	return c
}

type topoMap struct {
	grid    [][]int
	start   []coord
	numRows int
	numCols int
}

func newTopoMap(reader io.Reader) topoMap {
	row := 0
	var grid [][]int
	var start []coord

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		curr := make([]int, 0, len(line))

		for col, ch := range line {
			num, err := strconv.Atoi(string(ch))
			if err != nil {
				panic(err)
			}

			if num == 0 {
				start = append(start, coord{row, col})
			}

			curr = append(curr, num)
		}

		grid = append(grid, curr)
		row += 1
	}

	return topoMap{grid, start, row, len(grid[0])}
}

func (tm topoMap) inBounds(pos coord) bool {
	return pos.row >= 0 && pos.col >= 0 && pos.row < tm.numRows && pos.col < tm.numCols
}

func (tm topoMap) getTopo(pos coord) int {
	return tm.grid[pos.row][pos.col]
}

func (tm topoMap) scoreTrailhead(trailhead coord) int {
	type todo struct {
		pos coord
		topo int
	}

	score := 0
	visited := newSet[coord]()
	visited.insert(trailhead)
	queue := []todo{{trailhead, tm.getTopo(trailhead)}}

	for len(queue) != 0 {
		curr := queue[0]
		queue = queue[1:]

		if curr.topo == 9 {
			score += 1
		}

		for _, direction := range directions {
			nextPos := curr.pos.next(direction)
			if visited.contains(nextPos) || !tm.inBounds(nextPos) {
				continue
			}

			nextTopo := tm.getTopo(nextPos)
			topoDiff := nextTopo - curr.topo
			if topoDiff != 1 {
				continue
			}

			queue = append(queue, todo{nextPos, nextTopo})
			visited.insert(nextPos)
		}
	}

	return score
}

func (tm topoMap) totalScore() int {
	sum := 0

	for _, trailhead := range tm.start {
		score := tm.scoreTrailhead(trailhead)
		sum += score
	}

	return sum
}

func main() {
	tm := newTopoMap(os.Stdin)
	score := tm.totalScore()
	fmt.Println(score)
}
