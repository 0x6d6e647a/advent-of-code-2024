package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
)

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
func (tm topoMap) rate(stack []coord) int {
	rating := 0

	stackLen := len(stack)
	currPos := stack[stackLen-1]
	currTopo := tm.getTopo(currPos)

	if currTopo == 9 {
		rating += 1
	}

	for _, direction := range directions {
		nextPos := currPos.next(direction)
		if slices.Contains(stack, nextPos) || !tm.inBounds(nextPos) {
			continue
		}

		nextTopo := tm.getTopo(nextPos)
		topoDiff := nextTopo - currTopo
		if topoDiff != 1 {
			continue
		}

		newStack := append(stack, nextPos)
		rating += tm.rate(newStack)
	}

	return rating
}

func (tm topoMap) rateTrailhead(trailhead coord) int {
	return tm.rate([]coord{trailhead})
}

func (tm topoMap) totalRating() int {
	sum := 0

	for _, trailhead := range tm.start {
		rating := tm.rateTrailhead(trailhead)
		sum += rating
	}

	return sum
}

func main() {
	tm := newTopoMap(os.Stdin)
	rating := tm.totalRating()
	fmt.Println(rating)
}
