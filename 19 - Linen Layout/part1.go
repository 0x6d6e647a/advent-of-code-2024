package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type color byte

type request struct {
	towels  [][]color
	designs [][]color
	memory  map[string]bool
}

func colorSliceString(colors []color) string {
	var sb strings.Builder

	for _, c := range colors {
		sb.WriteByte(byte(c))
	}

	return sb.String()
}

func newRequest(r io.Reader) (req request) {
	foundBlank := false

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			if foundBlank {
				panic("invalid request input")
			}
			foundBlank = true
			continue
		}

		// -- Parse available towels.
		if !foundBlank {
			for _, str := range strings.Split(line, ", ") {
				var towel []color

				for _, ch := range str {
					c := color(ch)
					towel = append(towel, c)
				}

				req.towels = append(req.towels, towel)
			}

			continue
		}

		// -- Parse requested designs.
		var design []color

		for _, ch := range line {
			c := color(ch)
			design = append(design, c)
		}

		req.designs = append(req.designs, design)
	}

	req.memory = make(map[string]bool)
	return req
}

func (req request) isPossible(design []color) bool {
	if len(design) == 0 {
		return true
	}

	result, ok := req.memory[colorSliceString(design)]
	if ok {
		return result
	}

	for _, towel := range req.towels {
		if len(towel) > len(design) {
			continue
		}

		// -- See if towel fits design.
		numMatch := 0
		for i, c := range towel {
			if c == design[i] {
				numMatch += 1
			} else {
				break
			}
		}

		if numMatch != len(towel) {
			continue
		}

		// -- Check remaining design.
		if req.isPossible(design[len(towel):]) {
			req.memory[colorSliceString(design)] = true
			return true
		}
	}

	req.memory[colorSliceString(design)] = false
	return false
}

func (req request) score() int {
	score := 0

	for _, d := range req.designs {
		if req.isPossible(d) {
			score += 1
		}
	}

	return score
}

func main() {
	req := newRequest(os.Stdin)
	score := req.score()
	fmt.Println(score)
}
