package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func absDiffInt(x int, y int) int {
	if x < y {
		return y - x
	}
	return x - y
}

func main() {
	// -- Read input.
	var left []int
	var right []int

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
		right = append(right, numRight)
	}

	// -- Sort lists.
	sort.Ints(left)
	sort.Ints(right)

	// -- Determine total distance.
	distance := 0

	for i := 0; i < len(left); i += 1 {
		distance += absDiffInt(left[i], right[i])
	}

	fmt.Println(distance)
}
