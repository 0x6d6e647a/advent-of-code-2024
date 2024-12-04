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

var directions = map[string]coord{
	"N":  {-1, 0},
	"NE": {-1, 1},
	"E":  {0, 1},
	"SE": {1, 1},
	"S":  {1, 0},
	"SW": {1, -1},
	"W":  {0, -1},
	"NW": {-1, -1},
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

func (c crossword) countDirectionMatches(pos coord, word []byte) int {
	matches := 0

	for _, delta := range directions {
		if c.isMatchDirection(pos, word, delta) {
			matches += 1
		}
	}

	return matches
}

func (c crossword) countMatches(word []byte) int {
	matches := 0

	for row := 0; row < c.numRows; row += 1 {
		for col := 0; col < c.numCols; col += 1 {
			if c.grid[row][col] == word[0] {
				pos := coord{row, col}
				matches += c.countDirectionMatches(pos, word)
			}
		}
	}

	return matches
}

func main() {
	puzzle := newCrossword(os.Stdin)
	word := []byte("XMAS")
	matches := puzzle.countMatches(word)
	fmt.Println(matches)
}
