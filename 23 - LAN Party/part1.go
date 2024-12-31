package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
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

func (s set[T]) contains(value T) bool {
	_, ok := s[value]
	return ok
}

type lanMap struct {
	computers   set[string]
	connections map[string]set[string]
}

var lanMapLineRegex = regexp.MustCompile(`([A-Za-z]+)-([A-Za-z]+)`)

func parseLanMap(r io.Reader) (lm lanMap) {
	lm.computers = newSet[string]()
	lm.connections = make(map[string]set[string])

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		matches := lanMapLineRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			panic("invalid lan map string")
		}
		a := matches[1]
		b := matches[2]
		lm.computers.insert(a, b)

		_, aOk := lm.connections[a]
		if !aOk {
			lm.connections[a] = newSet[string]()
		}
		_, bOk := lm.connections[b]
		if !bOk {
			lm.connections[b] = newSet[string]()
		}
		lm.connections[a].insert(b)
		lm.connections[b].insert(a)
	}

	return lm
}

const NUM_EDGES = 3

func (lm lanMap) countInterconnected() int {
	cycles := make(map[string][NUM_EDGES]string)

	for a := range lm.connections {
		if a[0] != 't' {
			continue
		}

		for b := range lm.connections[a] {
			for c := range lm.connections[b] {
				if lm.connections[c].contains(a) {
					cycle := [NUM_EDGES]string{a, b, c}
					slices.Sort(cycle[:])
					key := fmt.Sprint(cycle)
					cycles[key] = cycle
				}
			}
		}
	}

	return len(cycles)
}

func main() {
	lm := parseLanMap(os.Stdin)
	n := lm.countInterconnected()
	fmt.Println(n)
}
