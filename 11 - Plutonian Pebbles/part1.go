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

type stoneLine []stone

func newStoneLine(r io.Reader) stoneLine {
	scanner := bufio.NewScanner(r)
	if !scanner.Scan() {
		panic("no input")
	}
	line := scanner.Text()
	strs := strings.Split(line, " ")
	stones := make([]stone, 0, len(strs))

	for _, str := range strs {
		num, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		stones = append(stones, stone(num))
	}

	return stones
}

func (sl *stoneLine) transform() {
	newStones := make([]stone, 0, len(*sl)*2)

	for _, stone := range *sl {
		transformed := stone.transform()
		newStones = append(newStones, transformed...)
	}

	*sl = newStones
}

func (sl *stoneLine) doTransforms(times int) {
	for range times {
		sl.transform()
	}
}

func main() {
	stones := newStoneLine(os.Stdin)
	stones.doTransforms(25)
	fmt.Println(len(stones))
}
