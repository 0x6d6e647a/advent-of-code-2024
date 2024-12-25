package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
)

type coord struct {
	row int
	col int
}

var directions = map[byte]coord{
	'^': {-1, 0},
	'v': {+1, 0},
	'>': {0, +1},
	'<': {0, -1},
}

func parseDirections(scanner *bufio.Scanner) (dirs []byte) {
	for scanner.Scan() {
		line := scanner.Text()
		for i := range line {
			dirs = append(dirs, line[i])
		}
	}

	return dirs
}

func (c coord) toward(dir byte) coord {
	delta := directions[dir]
	c.row += delta.row
	c.col += delta.col
	return c
}

type warehouse struct {
	walls []coord
	boxes []coord
	robot coord
}

func newWarehouse(scanner *bufio.Scanner) (w warehouse) {
	row := 0

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			break
		}

		for col, ch := range line {
			switch ch {
			case '#':
				w.walls = append(w.walls, coord{row, col})
			case 'O':
				w.boxes = append(w.boxes, coord{row, col})
			case '@':
				w.robot.row, w.robot.col = row, col
			case '.':
				continue
			default:
				panic("invalid warehouse character")
			}
		}

		row += 1
	}

	return w
}

func (w *warehouse) moveBox(pos coord, dir byte) bool {
	if !slices.Contains(w.boxes, pos) {
		panic("attempt to move non-box")
	}

	next := pos.toward(dir)

	if slices.Contains(w.walls, next) {
		return false
	}

	if slices.Contains(w.boxes, next) {
		if !w.moveBox(next, dir) {
			return false
		}
	}

	i := slices.Index(w.boxes, pos)
	w.boxes = slices.Delete(w.boxes, i, i+1)
	w.boxes = append(w.boxes, next)

	return true
}

func (w *warehouse) moveRobot(dir byte) {
	next := w.robot.toward(dir)

	if slices.Contains(w.walls, next) {
		return
	}

	if slices.Contains(w.boxes, next) {
		if !w.moveBox(next, dir) {
			return
		}
	}

	w.robot = next
}

func (w *warehouse) score() (score int) {
	for _, box := range w.boxes {
		score += (box.row * 100) + box.col
	}
	return score
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	wh := newWarehouse(s)
	for _, dir := range parseDirections(s) {
		wh.moveRobot(dir)
	}
	fmt.Println(wh.score())
}
