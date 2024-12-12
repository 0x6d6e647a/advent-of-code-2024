package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

type region struct {
	plant  byte
	coords set[coord]
}

func (r region) perimeter() int {
	perimeter := 0

	for curr := range r.coords {
		for _, direction := range directions {
			next := curr.move(direction)

			if !r.coords.contains(next) {
				perimeter += 1
			}
		}
	}

	return perimeter
}

func (r region) price() int {
	return r.perimeter() * len(r.coords)
}

type garden struct {
	grid    [][]byte
	numRows int
	numCols int
}

func newGarden(r io.Reader) garden {
	var grid [][]byte
	numRows := 0
	numCols := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]byte, 0, len(line))

		for _, ch := range line {
			row = append(row, byte(ch))
		}

		numRows += 1
		numCols = len(row)

		grid = append(grid, row)
	}

	return garden{grid, numRows, numCols}
}

func (g garden) inBounds(pos coord) bool {
	return pos.row >= 0 && pos.col >= 0 && pos.row < g.numRows && pos.col < g.numCols
}

func (g garden) at(pos coord) byte {
	return g.grid[pos.row][pos.col]
}

func (g garden) walkRegion(pos coord, visited *set[coord]) region {
	plant := g.at(pos)
	coords := newSet[coord]()
	coords.insert(pos)
	visited.insert(pos)

	queue := []coord{pos}

	for len(queue) != 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, direction := range directions {
			next := curr.move(direction)

			if !g.inBounds(next) || g.at(next) != plant || visited.contains(next) {
				continue
			}

			queue = append(queue, next)
			coords.insert(next)
			visited.insert(next)
		}
	}

	return region{plant, coords}
}

func (g garden) findRegions() []region {
	var regions []region
	visited := newSet[coord]()
	pos := coord{0, 0}

	for ; pos.row < g.numRows; pos.row += 1 {
		for ; pos.col < g.numCols; pos.col += 1 {
			if visited.contains(pos) {
				continue
			}

			region := g.walkRegion(pos, &visited)
			regions = append(regions, region)
		}

		pos.col = 0
	}

	return regions
}

func main() {
	garden := newGarden(os.Stdin)
	regions := garden.findRegions()

	price := 0
	for _, region := range regions {
		price += region.price()
	}

	fmt.Println(price)
}
