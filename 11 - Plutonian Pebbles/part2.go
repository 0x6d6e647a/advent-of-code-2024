package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type stone int

func powStone(base stone, exponent stone) stone {
	return stone(math.Pow(float64(base), float64(exponent)))
}

func (s stone) transform() []stone {
	// -- Rule 1.
	if s == 0 {
		return []stone{1}
	}

	// -- Rule 2.
	numDigits := stone(math.Log10(float64(s))) + 1
	if numDigits%2 == 0 {
		halfLen := numDigits / 2
		base10 := powStone(10, halfLen)
		left := s / base10
		right := s - (left * base10)
		return []stone{left, right}
	}

	// Rule 3.
	return []stone{s * 2024}
}

type stoneLine map[stone]int

func newStoneLine(r io.Reader) stoneLine {
	stones := make(map[stone]int)

	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		panic("no input")
	}
	line := scanner.Text()
	strs := strings.Split(line, " ")

	for _, str := range strs {
		num, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		stones[stone(num)] += 1
	}

	return stones
}

func (sl *stoneLine) transform() {
	transformed := make(map[stone]int)

	for stone, count := range *sl {
		newStones := stone.transform()
		for _, stone := range newStones {
			transformed[stone] += count
		}
	}

	*sl = transformed
}

func (sl *stoneLine) doTransforms(times int) {
	for range times {
		sl.transform()
	}
}

func (sl stoneLine) len() int {
	total := 0

	for _, count := range sl {
		total += count
	}

	return total
}

func main() {
	stones := newStoneLine(os.Stdin)
	stones.doTransforms(75)
	fmt.Println(stones.len())
}
