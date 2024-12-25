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

func (c *coord) move(dir byte) {
	delta := directions[dir]
	c.row += delta.row
	c.col += delta.col
}

type box [2]coord

func newBox(pos coord) (b box) {
	b[0], b[1] = pos, pos
	b[1].col += 1
	return b
}

func (b box) toward(dir byte) box {
	b[0].move(dir)
	b[1].move(dir)
	return b
}

func (b box) score() int {
	return (b[0].row * 100) + b[0].col
}

func (b box) collidePos(pos coord) bool {
	return b[0] == pos ||
		b[1] == pos
}

func (b box) collideBox(other box) bool {
	return b[0] == other[0] ||
		b[0] == other[1] ||
		b[1] == other[0] ||
		b[1] == other[1]
}

type warehouse struct {
	walls  []coord
	boxes  []box
	robot  coord
	maxRow int
	maxCol int
}

func newWarehouse(scanner *bufio.Scanner) (w warehouse) {
	row := 0

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			break
		}

		w.maxCol = len(line)

		for col, ch := range line {
			realCol := col * 2

			switch ch {
			case '#':
				w.walls = append(w.walls, coord{row, realCol})
				w.walls = append(w.walls, coord{row, realCol + 1})
			case 'O':
				w.boxes = append(w.boxes, newBox(coord{row, realCol}))
			case '@':
				w.robot.row, w.robot.col = row, realCol
			case '.':
				continue
			default:
				panic("invalid warehouse character")
			}
		}

		row += 1
	}

	w.maxRow = row
	w.maxCol *= 2

	return w
}

func (w *warehouse) moveBox(boxIndices []int, dir byte) bool {
	if len(boxIndices) == 0 {
		return true
	}

	altWh := *w

	currIndex := boxIndices[0]
	boxIndices = boxIndices[1:]
	next := altWh.boxes[currIndex].toward(dir)

	// -- Impossible due to wall collision.
	for _, wall := range altWh.walls {
		if next.collidePos(wall) {
			return false
		}
	}

	// -- Find collisions.
	for index, bb := range altWh.boxes {
		if index == currIndex {
			continue
		} else if next.collideBox(bb) {
			boxIndices = append(boxIndices, index)
		}
	}

	// -- Add collisions to queue.
	if !altWh.moveBox(boxIndices, dir) {
		return false
	}

	altWh.boxes[currIndex] = next
	w.boxes = altWh.boxes

	return true
}

func (w *warehouse) moveRobot(dir byte) {
	next := w.robot.toward(dir)

	if slices.Contains(w.walls, next) {
		return
	}

	boxIndex := slices.IndexFunc(w.boxes, func(bb box) bool {
		return bb.collidePos(next)
	})

	if boxIndex != -1 {
		if !w.moveBox([]int{boxIndex}, dir) {
			return
		}
	}

	w.robot = next
}

func (w *warehouse) score() (score int) {
	for _, b := range w.boxes {
		score += b.score()
	}
	return score
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	wh := newWarehouse(s)
	dirs := parseDirections(s)
	for _, dir := range dirs {
		wh.moveRobot(dir)
	}
	fmt.Println(wh.score())
}
