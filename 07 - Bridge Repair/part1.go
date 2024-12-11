package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

func powInt(base int, exponent int) int {
	return int(math.Pow(float64(base), float64(exponent)))
}

type operator uint8

const (
	add operator = iota
	mul
)

func op_add(a int, b int) int {
	return a + b
}

func op_mul(a int, b int) int {
	return a * b
}


func newOperator(bit uint) operator {
	switch bit {
	case 0:
		return add
	case 1:
		return mul
	default:
		panic("invalid operator bit")
	}
}

func generateOperators(byte uint, num int) []operator {
	operators := make([]operator, 0, num)

	for range num {
		bit := byte & 1
		byte >>= 1
		operators = append(operators, newOperator(bit))
	}

	return operators
}

func (o operator) toFunc() func(int, int) int {
	switch o {
	case add:
		return op_add
	case mul:
		return op_mul
	default:
		panic("invalid operator")
	}
}

type equation struct {
	goal       int
	components []int
}

func newEquation(line string) equation {
	colSplit := strings.Split(line, ":")

	// -- Parse goal number.
	goal, err := strconv.Atoi(colSplit[0])
	if err != nil {
		panic(err)
	}

	// Parse each component.
	var components []int
	for _, str := range strings.Split(colSplit[1], " ") {
		if len(str) == 0 {
			continue
		}

		component, err := strconv.Atoi(str)
		if err != nil {
			panic(err)
		}
		components = append(components, component)
	}

	return equation{goal, components}
}

func (e equation) compute(operators []operator) int {
	result := e.components[0]

	for index, component := range e.components[1:] {
		opFunc := operators[index].toFunc()
		result = opFunc(result, component)
	}

	return result
}

func (e equation) isPossible() bool {
	numOperations := len(e.components) - 1
	permutationLimit := powInt(2, numOperations) - 1

	for permutation := 0; permutation <= permutationLimit; permutation += 1 {
		operators := generateOperators(uint(permutation), numOperations)
		if e.compute(operators) == e.goal {
			return true
		}
	}
	return false
}

func main() {
	// -- Parse equations.
	var equations []equation

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		equation := newEquation(line)
		equations = append(equations, equation)

	}

	// -- Sum possible equations.
	sum := 0

	for _, equation := range equations {
		if equation.isPossible() {
			sum += equation.goal
		}
	}

	fmt.Println(sum)
}
