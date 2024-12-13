package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

type coord struct {
	x int
	y int
}

type clawMachine struct {
	a coord
	b coord
	p coord
}

var buttonRegex = regexp.MustCompile(`Button ([AB]): X\+(\d+), Y\+(\d+)`)
var prizeRegex = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

func newClawMachine(scanner *bufio.Scanner) *clawMachine {
	// -- Read lines.
	var lines [4]string

	for index := range lines {
		if !scanner.Scan() {
			if index == 3 {
				break
			}
			return nil
		}
		lines[index] = scanner.Text()
	}

	matchesA := buttonRegex.FindStringSubmatch(lines[0])
	matchesB := buttonRegex.FindStringSubmatch(lines[1])
	matchesP := prizeRegex.FindStringSubmatch(lines[2])

	if matchesA[1] != "A" || matchesB[1] != "B" {
		panic("button name mismatch")
	}

	if len(lines[3]) != 0 {
		panic("non-blank separator")
	}

	// -- Parse integers.
	aX, err := strconv.Atoi(matchesA[2])
	if err != nil {
		panic(err)
	}
	aY, err := strconv.Atoi(matchesA[3])
	if err != nil {
		panic(err)
	}
	bX, err := strconv.Atoi(matchesB[2])
	if err != nil {
		panic(err)
	}
	bY, err := strconv.Atoi(matchesB[3])
	if err != nil {
		panic(err)
	}
	pX, err := strconv.Atoi(matchesP[1])
	if err != nil {
		panic(err)
	}
	pY, err := strconv.Atoi(matchesP[2])
	if err != nil {
		panic(err)
	}

	return &clawMachine{
		a: coord{aX, aY},
		b: coord{bX, bY},
		p: coord{pX, pY},
	}
}

type clawMachineSolution struct {
	a int
	b int
}

func (cms clawMachineSolution) cost() int {
	return (cms.a * 3) + cms.b
}

func (cm clawMachine) solve() *clawMachineSolution {
	a := &cm.a
	b := &cm.b
	p := &cm.p

	det := a.x*b.y - b.x*a.y
	i := (b.y*p.x - b.x*p.y) / det
	j := (a.x*p.y - a.y*p.x) / det
	pos := coord{a.x*i + b.x*j, a.y*i + b.y*j}

	if pos != *p {
		return nil
	}

	return &clawMachineSolution{i, j}
}

func parseClawMachines(r io.Reader) []clawMachine {
	var clawMachines []clawMachine

	scanner := bufio.NewScanner(r)

	for {
		clawMachine := newClawMachine(scanner)
		if clawMachine == nil {
			break
		}
		clawMachines = append(clawMachines, *clawMachine)
	}

	return clawMachines
}

func main() {
	clawMachines := parseClawMachines(os.Stdin)
	total := 0

	for _, cm := range clawMachines {
		solution := cm.solve()

		if solution == nil {
			continue
		}

		total += solution.cost()
	}

	fmt.Println(total)
}
