package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

func pyModulo(numerator uint64, denominator uint64) uint64 {
	return ((numerator % denominator) + denominator) % denominator
}

type set[T comparable] map[T]struct{}

func newSet[T comparable]() set[T] {
	return make(map[T]struct{}, 0)
}

func (s set[T]) insert(value T) {
	s[value] = struct{}{}
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

const SEQ_LENGTH = 4

type priceMapper map[int8]set[[SEQ_LENGTH]int8]

func newPriceMapper() priceMapper {
	pm := make(priceMapper)
	for i := range 10 {
		pm[int8(i)] = newSet[[SEQ_LENGTH]int8]()
	}
	return pm
}

type monkey struct {
	initSecretNumber uint64
	secretNumbers    []uint64
	prices           []int8
	deltas           []int8
	priceMapper      priceMapper
}

const NUM_CHANGES = 2000

func newMonkey(sn uint64, metaPriceMapper *priceMapper) (m monkey) {
	m.initSecretNumber = sn
	m.priceMapper = newPriceMapper()

	for range NUM_CHANGES {
		prevPrice := int8(sn % 10)
		nsn := nextSecretNumber(sn)
		price := int8(nsn % 10)
		delta := price - prevPrice

		m.secretNumbers = append(m.secretNumbers, nsn)
		m.prices = append(m.prices, price)
		m.deltas = append(m.deltas, delta)
		sn = nsn

		if len(m.prices) < SEQ_LENGTH &&
			len(m.deltas) < SEQ_LENGTH {
			continue
		}

		var currDeltas [SEQ_LENGTH]int8
		copy(currDeltas[:], m.deltas[len(m.deltas)-SEQ_LENGTH:])

		m.priceMapper[price].insert(currDeltas)
		(*metaPriceMapper)[price].insert(currDeltas)
	}

	return m
}

func (m monkey) getPrice(target [SEQ_LENGTH]int8) int8 {
	for i := SEQ_LENGTH - 1; i < len(m.deltas); i += 1 {
		var currDeltas [SEQ_LENGTH]int8
		copy(currDeltas[:], m.deltas[i-SEQ_LENGTH+1:i+1])

		if currDeltas == target {
			return m.prices[i]
		}
	}

	return 0
}

func mostBananas(monkeys []monkey, metaPriceMapper priceMapper) (highest int64) {
	highest = math.MinInt64

	for price := int8(9); price > 0; price -= 1 {
		fmt.Printf("Checking %d's...\n", price)
		deltasSet := metaPriceMapper[price]

		for delta := range deltasSet {
			var sum int64 = 0

			for _, m := range monkeys {
				sum += int64(m.getPrice(delta))
			}

			if sum > highest {
				highest = sum
				fmt.Println(highest)
			}
		}
	}

	return highest
}

func main() {
	var monkeys []monkey
	metaPriceMapper := newPriceMapper()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		sn, err := strconv.ParseUint(line, 10, 64)
		if err != nil {
			panic(err)
		}
		monkey := newMonkey(sn, &metaPriceMapper)
		monkeys = append(monkeys, monkey)
	}

	highest := mostBananas(monkeys, metaPriceMapper)
	fmt.Println(highest)
}
