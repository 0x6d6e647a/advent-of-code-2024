package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"slices"
)

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
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

type racetrack struct {
	start coord
	end   coord
	limit coord
	walls []coord
}

func newRacetrack(rdr io.Reader) (rt racetrack) {
	row := 0

	scanner := bufio.NewScanner(rdr)
	for scanner.Scan() {
		line := scanner.Text()
		for col, ch := range line {
			switch ch {
			case '.':
				continue
			case '#':
				rt.walls = append(rt.walls, coord{row, col})
			case 'S':
				rt.start.row = row
				rt.start.col = col
			case 'E':
				rt.end.row = row
				rt.end.col = col
			default:
				panic("invalid race character")
			}
		}
		rt.limit.col = len(line)
		row += 1
	}

	rt.limit.row = row
	return rt
}

func (rt racetrack) inBounds(pos coord) bool {
	return pos.row >= 0 &&
		pos.col >= 0 &&
		pos.row < rt.limit.row &&
		pos.col < rt.limit.col
}

func (rt racetrack) isWalled(pos coord) bool {
	return slices.Contains(rt.walls, pos)
}

func (rt racetrack) getPath() []coord {
	pos := rt.start
	path := []coord{pos}

	for pos != rt.end {
		var next coord

		for _, delta := range directions {
			next = pos
			next.row += delta.row
			next.col += delta.col

			if rt.inBounds(next) &&
				!rt.isWalled(next) &&
				!slices.Contains(path, next) {
				break
			}
		}

		pos = next
		path = append(path, pos)
	}

	return path
}

func (rt racetrack) findCheats(numCheats int, numSaved int) (count int) {
	path := rt.getPath()

	for i, src := range path[:len(path)-numSaved] {
		for j, dst := range path[i+numSaved:] {
			d := dst
			d.row -= src.row
			d.col -= src.col
			dd := absInt(d.row) + absInt(d.col)

			if dd <= numCheats && dd <= j {
				count += 1
			}
		}
	}

	return count
}

func main() {
	rt := newRacetrack(os.Stdin)
	cheats := rt.findCheats(20, 100)
	fmt.Println(cheats)
}
