package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

var regex_instruction *regexp.Regexp = regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)`)

func scan(line string, sum *int) {
	instructions := regex_instruction.FindAllStringSubmatch(line, -1)

	for _, instruction := range instructions {
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
	sum := 0

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		scan(line, &sum)
	}

	fmt.Println(sum)
}
