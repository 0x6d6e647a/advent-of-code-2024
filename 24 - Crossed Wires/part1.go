package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"slices"
)

type operation uint8

const (
	AND operation = iota
	OR
	XOR
)

func newOperation(s string) operation {
	switch s {
	case "AND":
		return AND
	case "OR":
		return OR
	case "XOR":
		return XOR
	default:
		panic("invalid operation string")
	}
}

type register [3]byte

type gate struct {
	a   register
	op  operation
	b   register
	out register
}

var gateRegex = regexp.MustCompile(`(\w{3}) (\w{2,3}) (\w{3}) -> (\w{3})`)

func newGate(line string) (g gate) {
	matches := gateRegex.FindStringSubmatch(line)
	if len(matches) != 5 {
		panic("invalid gate string")
	}

	copy(g.a[:], []byte(matches[1]))
	g.op = newOperation(matches[2])
	copy(g.b[:], []byte(matches[3]))
	copy(g.out[:], []byte(matches[4]))

	return g
}

func (g gate) run(d *device) {
	a, aOk := d.wires[g.a]
	b, bOk := d.wires[g.b]
	if !aOk || !bOk {
		panic("missing device state")
	}

	switch g.op {
	case AND:
		d.wires[g.out] = a && b
	case OR:
		d.wires[g.out] = a || b
	case XOR:
		d.wires[g.out] = a != b
	default:
		panic("invalid operation")
	}
}

type device struct {
	wires map[register]bool
	gates []gate
}

var stateRegex = regexp.MustCompile(`(\w{3}): ([01])`)

func newDevice(r io.Reader) (d device) {
	scanner := bufio.NewScanner(r)

	// -- Parse wires.
	d.wires = make(map[register]bool)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		matches := stateRegex.FindStringSubmatch(line)
		if len(matches) != 3 {
			panic("invalid state line")
		}

		var reg register
		copy(reg[:], []byte(matches[1]))

		switch matches[2][0] {
		case '0':
			d.wires[reg] = false
		case '1':
			d.wires[reg] = true
		default:
			panic("invalid state boolean")
		}
	}

	// -- Parse gates.
	for scanner.Scan() {
		line := scanner.Text()
		gate := newGate(line)
		d.gates = append(d.gates, gate)
	}

	return d
}

func (d device) runGates() device {
	open := make([]gate, len(d.gates))
	copy(open, d.gates)

	for len(open) != 0 {
		for i, g := range open {
			_, aOk := d.wires[g.a]
			_, bOk := d.wires[g.b]
			_, outOk := d.wires[g.out]
			if aOk && bOk && !outOk {
				g.run(&d)
				open = slices.Delete(open, i, i+1)
				break
			}
		}
	}

	return d
}

func (d device) getZNumber() (z uint64) {
	reg := register{'z', 0, 0}

	for i := range 100 {
		reg[1] = byte(i/10) + '0'
		reg[2] = byte(i%10) + '0'
		v, ok := d.wires[reg]
		if !ok {
			break
		}

		if v {
			z |= 1 << i
		}
	}

	return z
}

func main() {
	device := newDevice(os.Stdin)
	done := device.runGates()
	n := done.getZNumber()
	fmt.Println(n)
}
