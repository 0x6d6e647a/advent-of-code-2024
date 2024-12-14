package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
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


func (b bathroom) isChristmasTree() bool {
	const TEST_DEPTH = 4

eachRobot:
	for pos := range b.robots {
		for depth := range TEST_DEPTH {
			row := pos.row + depth
			r := coord{row, pos.col + depth}
			l := coord{row, pos.col - depth}

			if len(b.robots[l]) == 0 || len(b.robots[r]) == 0 {
				continue eachRobot
			}
		}

		return true
	}

	return false
}

const WIDTH = 101
const HEIGHT = 103

func main() {
	b := newBathroom(os.Stdin, WIDTH, HEIGHT)
	for n := range math.MaxInt64 {
		b.tick()
		if b.isChristmasTree() {
			fmt.Println(n + 1)
			return
		}
	}
}
