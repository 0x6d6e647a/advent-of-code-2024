package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

func pyModulo(numerator int, denominator int) int {
	return ((numerator % denominator) + denominator) % denominator
}

type coord struct {
	row int
	col int
}

func (c coord) add(delta coord) coord {
	c.row += delta.row
	c.col += delta.col
	return c
}

var regexRobot = regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)

func newRobot(line string) (pos coord, vel coord) {
	matches := regexRobot.FindStringSubmatch(line)
	if len(matches) != 5 {
		panic("invalid robot string")
	}

	var err error

	pos.col, err = strconv.Atoi(matches[1])
	if err != nil {
		panic(err)
	}
	pos.row, err = strconv.Atoi(matches[2])
	if err != nil {
		panic(err)
	}
	vel.col, err = strconv.Atoi(matches[3])
	if err != nil {
		panic(err)
	}
	vel.row, err = strconv.Atoi(matches[4])
	if err != nil {
		panic(err)
	}

	return pos, vel
}

type bathroom struct {
	rowMax int
	colMax int
	rowMid int
	colMid int
	robots map[coord][]coord
}

func newBathroom(r io.Reader, width int, height int) (b bathroom) {
	b.rowMax = height
	b.colMax = width
	b.rowMid = height / 2
	b.colMid = width / 2
	b.robots = make(map[coord][]coord)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		pos, vel := newRobot(line)
		b.robots[pos] = append(b.robots[pos], vel)
	}

	return b
}

func (b *bathroom) moveRobot(pos coord, delta coord) coord {
	newPos := pos.add(delta)
	newPos.row = pyModulo(newPos.row, b.rowMax)
	newPos.col = pyModulo(newPos.col, b.colMax)
	return newPos
}

func (b *bathroom) tick() {
	newRobots := make(map[coord][]coord)

	for pos, vels := range b.robots {
		for _, vel := range vels {
			newPos := b.moveRobot(pos, vel)
			newRobots[newPos] = append(newRobots[newPos], vel)
		}
	}

	b.robots = newRobots
}

type quadrant uint8

const (
	NO_QUADRANT quadrant = iota
	FIRST_QUADRANT
	SECOND_QUADRANT
	THIRD_QUADRANT
	FOURTH_QUADRANT
)

func (b *bathroom) getQuadrant(pos coord) quadrant {
	if pos.row < b.rowMid {
		if pos.col < b.colMid {
			return FIRST_QUADRANT
		} else if pos.col > b.colMid {
			return SECOND_QUADRANT
		}
	} else if pos.row > b.rowMid {
		if pos.col < b.colMid {
			return THIRD_QUADRANT
		} else if pos.col > b.colMid {
			return FOURTH_QUADRANT
		}
	}

	return NO_QUADRANT
}

func (b *bathroom) score() int {
	quads := make(map[quadrant]int)

	for pos, vels := range b.robots {
		q := b.getQuadrant(pos)
		quads[q] += len(vels)
	}

	score := 1

	for q, n := range quads {
		if q == NO_QUADRANT {
			continue
		}

		score *= n
	}

	return score
}

const WIDTH = 101
const HEIGHT = 103

func main() {
	b := newBathroom(os.Stdin, WIDTH, HEIGHT)
	for range 100 {
		b.tick()
	}
	fmt.Println(b.score())
}
