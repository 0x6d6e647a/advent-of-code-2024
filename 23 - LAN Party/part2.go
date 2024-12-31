package main

import (
	"bufio"
	"fmt"
	"io"
	"maps"
	"math"
	"os"
	"regexp"
	"slices"
	"strings"
)

type set[T comparable] map[T]struct{}

func newSet[T comparable]() set[T] {
	return make(map[T]struct{}, 0)
}

func (s set[T]) clone() set[T] {
	return maps.Clone(s)
}

func (s set[T]) insert(values ...T) {
	for _, value := range values {
		s[value] = struct{}{}
	}
}

func (s set[T]) erase(values ...T) {
	for _, v := range values {
		delete(s, v)
	}
}

func (s set[T]) contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s set[T]) intersection(o set[T]) (i set[T]) {
	i = make(set[T])

	for v := range s {
		if o.contains(v) {
			i.insert(v)
		}
	}

	return i
}

func (s set[T]) toSlice() []T {
	ss := make([]T, 0, len(s))
	for v := range s {
		ss = append(ss, v)
	}
	return ss
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

func bronKerbosch[T comparable](graph map[T]set[T], r, p, x set[T], maxLenClique *set[T], maxLen *int) {
	if len(p) == 0 &&
		len(x) == 0 &&
		*maxLen < len(r) {
		*maxLenClique = r
		*maxLen = len(r)
		return
	}

	for v := range p.clone() {
		rr := r.clone()
		rr.insert(v)
		neighbors := graph[v]
		pp := p.intersection(neighbors)
		xx := x.intersection(neighbors)
		bronKerbosch(graph, rr, pp, xx, maxLenClique, maxLen)
		p.erase(v)
		x.insert(v)
	}
}

func (lm lanMap) findPassword() string {
	maxLenClique := newSet[string]()
	maxLen := math.MinInt
	r := newSet[string]()
	p := lm.computers.clone()
	x := newSet[string]()
	bronKerbosch(lm.connections, r, p, x, &maxLenClique, &maxLen)

	result := maxLenClique.toSlice()
	slices.Sort(result)

	return strings.Join(result, ",")
}

func main() {
	lm := parseLanMap(os.Stdin)
	pw := lm.findPassword()
	fmt.Println(pw)
}
