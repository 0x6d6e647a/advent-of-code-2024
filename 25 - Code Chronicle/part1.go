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

func (s set[T]) insert(values ...T) {
	for _, value := range values {
		s[value] = struct{}{}
	}
}

const HEIGHT = 7
const WIDTH = 5

type keyLock [HEIGHT][WIDTH]bool

func parseInput(r io.Reader) (kls []keyLock) {
	var curr keyLock
	row := 0

	putLock := func() {
		kls = append(kls, curr)
		row = 0
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			putLock()
			continue
		}

		for col, ch := range line {
			curr[row][col] = ch == '#'
		}
		row += 1
	}

	putLock()

	return kls
}

func (kl keyLock) fits(other keyLock) bool {
	for row := 0; row < HEIGHT; row += 1 {
		for col := 0; col < WIDTH; col += 1 {
			if kl[row][col] && other[row][col] {
				return false
			}
		}
	}

	return true
}

func main() {
	kls := parseInput(os.Stdin)

	tracker := newSet[[2]int]()
	for i := 0; i < len(kls); i += 1 {
		for j := 0; j < len(kls); j += 1 {
			if kls[i].fits(kls[j]) {
				tracker.insert([2]int{i, j})
			}
		}
	}

	fmt.Println(len(tracker) / 2)
}
