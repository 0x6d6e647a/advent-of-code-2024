package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

type coord struct {
	row int
	col int
}

func (c *coord) move(delta coord) {
	c.row += delta.row
	c.col += delta.col
}

func (c coord) add(other coord) coord {
	return coord{c.row + other.row, c.col + other.col}
}

var directions = map[string]coord{
	"NE": {-1, 1},
	"SE": {1, 1},
	"SW": {1, -1},
	"NW": {-1, -1},
}

var direction_opposite = map[string]string{
	"NE": "SW",
	"SE": "NW",
	"SW": "NE",
	"NW": "SE",
}

type crossword struct {
	grid    [][]byte
	numRows int
	numCols int
}

func newCrossword(r io.Reader) crossword {
	var grid [][]byte
	numRows := 0

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		line := []byte(text)
		grid = append(grid, line)
		numRows += 1
	}

	numCols := len(grid[0])
	return crossword{grid, numRows, numCols}
}

func (c crossword) inBounds(pos coord) bool {
	return pos.row >= 0 && pos.col >= 0 && pos.row < c.numRows && pos.col < c.numCols
}

func (c crossword) get(pos coord) (byte, error) {
	if !c.inBounds(pos) {
		return '\x00', errors.New("out of bounds")
	}
	return c.grid[pos.row][pos.col], nil
}

func (c crossword) isMatchDirection(pos coord, word []byte, delta coord) bool {
	for _, word_char := range word {
		cw_char, err := c.get(pos)
		if err != nil {
			return false
		}

		if cw_char != word_char {
			return false
		}

		pos.move(delta)
	}

	return true
}

func (c crossword) isMatchX(pos coord, word []byte) bool {
	matches := 0

	for direction, delta := range directions {
		opposite_direction := direction_opposite[direction]
		start := pos.add(directions[opposite_direction])

		if c.isMatchDirection(start, word, delta) {
			matches += 1

			if matches == 2 {
				return true
			}
		}
	}

	return false
}

func (c crossword) countMatches(word []byte, mid byte) int {
	matches := 0

	for row := 0; row < c.numRows; row += 1 {
		for col := 0; col < c.numCols; col += 1 {
			char := c.grid[row][col]
			if char == mid {
				pos := coord{row, col}
				if c.isMatchX(pos, word) {
					matches += 1
				}
			}
		}
	}

	return matches
}

func main() {
	puzzle := newCrossword(os.Stdin)
	word := []byte("MAS")
	mid := byte('A')
	matches := puzzle.countMatches(word, mid)
	fmt.Println(matches)
}
