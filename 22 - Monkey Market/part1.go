package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func pyModulo(numerator uint64, denominator uint64) uint64 {
	return ((numerator % denominator) + denominator) % denominator
}

func nextSecretNumber(curr uint64) uint64 {
	mix := func(sn uint64, v uint64) uint64 {
		return v ^ sn
	}
	prune := func(sn uint64) uint64 {
		return pyModulo(sn, 16777216)
	}
	step0 := func(sn uint64) uint64 {
		v0 := sn * 64
		v1 := mix(sn, v0)
		v2 := prune(v1)
		return v2
	}
	step1 := func(sn uint64) uint64 {
		v0 := sn / 32
		v1 := mix(sn, v0)
		v2 := prune(v1)
		return v2
	}
	step2 := func(sn uint64) uint64 {
		v0 := sn * 2048
		v1 := mix(sn, v0)
		v2 := prune(v1)
		return v2
	}
	s0 := step0(curr)
	s1 := step1(s0)
	s2 := step2(s1)
	return s2
}

func main() {
	var buyers []uint64
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		sn, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			panic(err)
		}
		buyers = append(buyers, sn)
	}

	for range 2000 {
		for i, sn := range buyers {
			buyers[i] = nextSecretNumber(sn)
		}
	}

	var sum uint64 = 0
	for _, sn := range buyers {
		sum += sn
	}

	fmt.Println(sum)
}
