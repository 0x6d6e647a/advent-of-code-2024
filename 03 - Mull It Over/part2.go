package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var regex_instruction *regexp.Regexp = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do(n't)?\(\)`)

func scan(line string, sum *int, enabled *bool) {
	instructions := regex_instruction.FindAllStringSubmatch(line, -1)

	for _, instruction := range instructions {
		// -- Check for control flow.
		if instruction[0] == "don't()" {
			*enabled = false
			continue
		}

		if instruction[0] == "do()" {
			*enabled = true
			continue
		}

		// -- Skip if not enabled.
		if !*enabled {
			continue
		}

		// -- Parse and add to sum.
		left, err := strconv.Atoi(instruction[1])
		if err != nil {
			panic(err)
		}
		right, err := strconv.Atoi(instruction[2])
		if err != nil {
			panic(err)
		}

		*sum += left * right
	}
}

func main() {
	enabled := true
	sum := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		scan(line, &sum, &enabled)
	}

	fmt.Println(sum)
}
