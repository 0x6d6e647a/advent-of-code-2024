package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	// -- Read input.
	var left []int
	right := make(map[int]int)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		nums := strings.Fields(line)

		numLeft, err := strconv.Atoi(nums[0])
		if err != nil {
			panic(err)
		}
		left = append(left, numLeft)

		numRight, err := strconv.Atoi(nums[1])
		if err != nil {
			panic(err)
		}
		if _, ok := right[numRight]; !ok {
			right[numRight] = 0
		}
		right[numRight] += 1
	}

	// -- Determine similarity score.
	similarity := 0

	for _, numLeft := range left {
		similarity += numLeft * right[numLeft]
	}

	fmt.Println(similarity)
}
