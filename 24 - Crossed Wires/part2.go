package main

import (
	"bufio"
	"fmt"
	"io"
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

func (s set[T]) insert(values ...T) {
	for _, value := range values {
		s[value] = struct{}{}
	}
}

func (s set[T]) toSlice() []T {
	ss := make([]T, 0, len(s))
	for v := range s {
		ss = append(ss, v)
	}
	return ss
}

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

func (r register) String() string {
	var s [3]rune
	s[0] = rune(r[0])
	s[1] = rune(r[1])
	s[2] = rune(r[2])
	return string(s[:])
}

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
		fmt.Println(line)
		panic("invalid gate string")
	}

	copy(g.a[:], []byte(matches[1]))
	g.op = newOperation(matches[2])
	copy(g.b[:], []byte(matches[3]))
	copy(g.out[:], []byte(matches[4]))

	return g
}

type device struct {
	highestZ register
	wires    map[register]bool
	gates    map[register]gate
}

var stateRegex = regexp.MustCompile(`(\w+): ([01])`)

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
	highestZ := math.MinInt
	d.gates = make(map[register]gate)

	for scanner.Scan() {
		line := scanner.Text()
		gate := newGate(line)
		d.gates[gate.out] = gate

		if gate.out[0] == 'z' {
			tens := int(gate.out[1] - '0')
			ones := int(gate.out[2] - '0')

			if tens < 0 ||
				tens > 9 ||
				ones < 0 ||
				ones > 9 {
				continue
			}

			z := tens*10 + ones

			if z > highestZ {
				highestZ = z
				d.highestZ = gate.out
			}
		}
	}

	return d
}

func (d device) findSwappedWires() []register {
	wrong := newSet[register]()

	prefixes := [...]byte{'x', 'y', 'z'}
	x0 := register{'x', '0', '0'}

	for _, gate := range d.gates {
		out0 := gate.out[0]

		if out0 == 'z' &&
			gate.op != XOR &&
			gate.out != d.highestZ {
			wrong.insert(gate.out)
		}

		if gate.op == XOR &&
			!slices.Contains(prefixes[:], out0) &&
			!slices.Contains(prefixes[:], gate.a[0]) &&
			!slices.Contains(prefixes[:], gate.b[0]) {
			wrong.insert(gate.out)
		}

		if gate.op == AND &&
			gate.a != x0 &&
			gate.b != x0 {
			for _, other := range d.gates {
				if (gate.out == other.a || gate.out == other.b) &&
					other.op != OR {
					wrong.insert(gate.out)
				}
			}
		}

		if gate.op == XOR {
			for _, other := range d.gates {
				if (gate.out == other.a || gate.out == other.b) &&
					other.op == OR {
					wrong.insert(gate.out)
				}
			}
		}
	}

	return wrong.toSlice()
}

func main() {
	device := newDevice(os.Stdin)
	wires := device.findSwappedWires()

	wireNames := make([]string, len(wires))
	for i, w := range wires {
		wireNames[i] = w.String()
	}
	slices.Sort(wireNames)
	fmt.Println(strings.Join(wireNames, ","))
}
