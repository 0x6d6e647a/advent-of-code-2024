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
	cat
)

func op_add(a int, b int) int {
	return a + b
}

func op_mul(a int, b int) int {
	return a * b
}

func op_cat(a int, b int) int {
	lenB := int(math.Log10(float64(b))) + 1
	return a * powInt(10, lenB) + b
}

func newOperator(num int) operator {
	switch num {
	case 0:
		return add
	case 1:
		return mul
	case 2:
		return cat
	default:
		panic("invalid operator bit")
	}
}

func nextPermuatation(operators []int) []int {
	for index := range operators {
		operators[index] += 1

		if operators[index] > int(cat) {
			operators[index] = 0
		} else {
			break
		}
	}

	return operators
}

func generateOperators(opInts []int) []operator {
	operators := make([]operator, 0, len(opInts))

	for _, opInt := range opInts {
		operators = append(operators, newOperator(opInt))
	}

	return operators
}

func (o operator) toFunc() func(int, int) int {
	switch o {
	case add:
		return op_add
	case mul:
		return op_mul
	case cat:
		return op_cat
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
	permutation := make([]int, numOperations)
	permutationLimit := powInt(3, numOperations) - 1

	for range permutationLimit + 1 {
		operators := generateOperators(permutation)
		if e.compute(operators) == e.goal {
			return true
		}
		permutation = nextPermuatation(permutation)
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
