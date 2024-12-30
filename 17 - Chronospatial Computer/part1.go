package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type instruction uint8

const (
	adv instruction = 0
	bxl instruction = 1
	bst instruction = 2
	jnz instruction = 3
	bxc instruction = 4
	out instruction = 5
	bdv instruction = 6
	cdv instruction = 7
)

type operand uint8

const (
	zero    operand = 0
	one     operand = 1
	two     operand = 2
	three   operand = 3
	regA    operand = 4
	regB    operand = 5
	regC    operand = 6
	unknown operand = 7
)

type computer struct {
	a int64
	b int64
	c int64
}

func (comp computer) getOperand(op operand) int64 {
	switch op {
	case zero:
		return 0
	case one:
		return 1
	case two:
		return 2
	case three:
		return 3
	case regA:
		return comp.a
	case regB:
		return comp.b
	case regC:
		return comp.c
	default:
		panic("invalid operand")
	}
}

func (comp computer) run(prog []uint8) (output []uint8) {
	var step int

	for ip := 0; ip < len(prog); ip += step {
		step = 2
		inst := instruction(prog[ip])
		op := prog[ip+1]
		comp.do(inst, op, &ip, &step, &output)
	}

	return output
}

func (comp *computer) do(inst instruction, op uint8, ip *int, step *int, output *[]uint8) {
	div := func() int64 {
		numerator := comp.getOperand(regA)
		combo := comp.getOperand(operand(op))
		var denominator int64 = 1 << combo
		quotient := numerator / denominator
		return quotient
	}

	switch inst {
	case adv:
		quotient := div()
		comp.a = quotient
	case bxl:
		b := comp.getOperand(regB)
		xor := b ^ int64(op)
		comp.b = xor
	case bst:
		combo := comp.getOperand(operand(op))
		modulo := combo & 0x07
		comp.b = modulo
	case jnz:
		a := comp.getOperand(regA)
		if a != 0 {
			*ip = int(op)
			*step = 0
		}
	case bxc:
		b := comp.getOperand(regB)
		c := comp.getOperand(regC)
		xor := b ^ c
		comp.b = xor
	case out:
		combo := comp.getOperand(operand(op))
		modulo := uint8(combo & 0x07)
		*output = append(*output, modulo)
	case bdv:
		quotient := div()
		comp.b = quotient
	case cdv:
		quotient := div()
		comp.c = quotient
	default:
		panic("instruction not implemented")
	}
}

var registerRegex = regexp.MustCompile(`Register ([A-Z]): (\d+)`)
var programRegex = regexp.MustCompile(`Program: (.*)`)

func parseInput(r io.Reader) (comp computer, prog []uint8) {
	// -- Parse registers.
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}
		matches := registerRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			panic("invalid register string")
		}
		letter := matches[1][0]
		val, err := strconv.Atoi(matches[2])
		if err != nil {
			panic("invalid register value")
		}

		switch letter {
		case 'A':
			comp.a = int64(val)
		case 'B':
			comp.b = int64(val)
		case 'C':
			comp.c = int64(val)
		default:
			panic("invalid register letter")
		}

	}

	// -- Parse program.
	if !scanner.Scan() {
		panic("missing program line")
	}
	line := scanner.Text()
	matches := programRegex.FindStringSubmatch(line)
	if len(matches) != 2 {
		panic("invalid program string")
	}

	numStrs := strings.Split(matches[1], ",")
	prog = make([]uint8, len(numStrs))
	for i, numStr := range numStrs {
		num, err := strconv.ParseUint(numStr, 10, 8)
		if err != nil {
			panic("invalid program number")
		}
		prog[i] = uint8(num)
	}

	return comp, prog
}

func main() {
	comp, prog := parseInput(os.Stdin)
	output := comp.run(prog)
	for i, n := range output {
		fmt.Print(n)
		if i+1 != len(output) {
			fmt.Print(",")
		}
	}
	fmt.Println()
}
