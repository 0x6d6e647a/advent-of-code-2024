package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
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

type bounds struct {
	minRow int
	maxRow int
	minCol int
	maxCol int
}

func newBounds() bounds {
	return bounds{
		minRow: math.MaxInt,
		maxRow: math.MinInt,
		minCol: math.MaxInt,
		maxCol: math.MinInt,
	}
}

func (b *bounds) expand(pos coord) {
	if pos.row > b.maxRow {
		b.maxRow = pos.row
	}
	if pos.row < b.minRow {
		b.minRow = pos.row
	}
	if pos.col > b.maxCol {
		b.maxCol = pos.col
	}
	if pos.col < b.minCol {
		b.minCol = pos.col
	}
}

type region struct {
	plant  byte
	coords set[coord]
	limits bounds
}

func (r region) countSides() int {
	var pos coord
	sides := 0

	// -- Do vertical.
	for pos.col = r.limits.minCol; pos.col <= r.limits.maxCol; pos.col += 1 {
		onWestFence := false
		onEastFence := false

		for pos.row = r.limits.minRow; pos.row <= r.limits.maxRow; pos.row += 1 {
			for _, direction := range [...]direction{WEST, EAST} {
				var onDirFence *bool
				switch direction {
				case WEST:
					onDirFence = &onWestFence
				case EAST:
					onDirFence = &onEastFence
				default:
					panic("invalid direction")
				}

				if r.coords.contains(pos) && !r.coords.contains(pos.move(direction)) {
					if !*onDirFence {
						sides += 1
					}
					*onDirFence = true
				} else {
					*onDirFence = false
				}
			}
		}
	}

	// -- Do horizontal.
	for pos.row = r.limits.minRow; pos.row <= r.limits.maxRow; pos.row += 1 {
		onNorthFence := false
		onSouthFence := false

		for pos.col = r.limits.minCol; pos.col <= r.limits.maxCol; pos.col += 1 {
			for _, direction := range [...]direction{NORTH, SOUTH} {
				var onDirFence *bool
				switch direction {
				case NORTH:
					onDirFence = &onNorthFence
				case SOUTH:
					onDirFence = &onSouthFence
				default:
					panic("invalid direction")
				}

				if r.coords.contains(pos) && !r.coords.contains(pos.move(direction)) {
					if !*onDirFence {
						sides += 1
					}
					*onDirFence = true
				} else {
					*onDirFence = false
				}
			}
		}
	}

	return sides
}

func (r region) price() int {
	return r.countSides() * len(r.coords)
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
	limits := newBounds()
	limits.expand(pos)

	queue := []coord{pos}
	visited.insert(pos)

	for len(queue) != 0 {
		curr := queue[0]
		queue = queue[1:]

		for _, direction := range directions {
			next := curr.move(direction)

			if !g.inBounds(next) || g.at(next) != plant || visited.contains(next) {
				continue
			}

			coords.insert(next)
			limits.expand(next)
			queue = append(queue, next)
			visited.insert(next)
		}
	}

	return region{plant, coords, limits}

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

	total := 0
	for _, region := range regions {
		price := region.price()
		total += price
	}
	fmt.Println(total)
}
