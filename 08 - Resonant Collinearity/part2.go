package main

import (
	"bufio"
	"fmt"
	"io"
	"math/bits"
	"os"
)

type set[T comparable] map[T]struct{}

func newSet[T comparable]() set[T] {
	return make(map[T]struct{}, 0)
}

func (s set[T]) insert(value T) {
	s[value] = struct{}{}
}

func generateCombinations[T any](arr []T, size int) <-chan []T {
	// -- Normalize requested size.
	length := uint(len(arr))

	if size > len(arr) {
		size = len(arr)
	}

	// -- Output to channel.
	out := make(chan []T)

	go func() {
		defer close(out)

		for comboBits := 1; comboBits < (1 << length); comboBits += 1 {
			if size > 0 && bits.OnesCount(uint(comboBits)) != size {
				continue
			}

			var combo []T

			for index := uint(0); index < length; index += 1 {
				if (comboBits>>index)&1 == 1 {
					combo = append(combo, arr[index])
				}
			}

			out <- combo
		}
	}()

	return out
}

type coord struct {
	row int
	col int
}

type slope struct {
	rise int
	run  int
}

type line struct {
	start coord
	end   coord
}

func (l line) getSlope() slope {
	rise := l.end.row - l.start.row
	run := l.end.col - l.start.col
	return slope{rise, run}
}

func (l line) getAntiNodes(cm cityMap) []coord {
	var antinodes []coord
	slope := l.getSlope()

	// -- Get backwards antinodes.
	rev := l.start
	for {
		rev.row -= slope.rise
		rev.col -= slope.run

		if !cm.inBounds(rev) {
			break
		}

		antinodes = append(antinodes, rev)
	}

	// -- Get forwards antinodes.
	fwd := l.end
	for {
		fwd.row += slope.rise
		fwd.col += slope.run

		if !cm.inBounds(fwd) {
			break
		}

		antinodes = append(antinodes, fwd)
	}

	return antinodes
}

type cityMap struct {
	numRows  int
	numCols  int
	antennas map[byte][]coord
}

func newCityMap(r io.Reader) cityMap {
	row := 0
	col := 0
	antennas := make(map[byte][]coord)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		for index, char := range line {
			char := byte(char)
			if char == '.' {
				continue
			}

			antennas[char] = append(antennas[char], coord{row, index})
		}

		row += 1
		col = len(line)
	}

	return cityMap{row, col, antennas}
}

func (cm cityMap) inBounds(pos coord) bool {
	return pos.row >= 0 && pos.col >= 0 && pos.row < cm.numRows && pos.col < cm.numCols
}

func (cm cityMap) getAntiNodes() set[coord] {
	antinodes := newSet[coord]()

	for _, coords := range cm.antennas {
		for _, coord := range coords {
			antinodes.insert(coord)
		}

		for combo := range generateCombinations(coords, 2) {
			if len(combo) != 2 {
				continue
			}

			line := line{combo[0], combo[1]}

			for _, antinode := range line.getAntiNodes(cm) {
				antinodes.insert(antinode)
			}
		}
	}

	return antinodes
}

func main() {
	cityMap := newCityMap(os.Stdin)
	antinodes := cityMap.getAntiNodes()
	fmt.Println(len(antinodes))
}
