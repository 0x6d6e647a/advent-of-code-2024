package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func absInt(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

type Report []int

func isSafe(report Report) bool {
	// -- Change direction type.
	type Change uint8

	const (
		Unset Change = iota
		Increasing
		Decreasing
	)

	// -- Check safety for report.
	change := Unset

	for i := 1; i < len(report); i += 1 {
		left := report[i-1]
		right := report[i]
		diff := left - right

		// -- Check if difference is safe.
		diffAbs := absInt(diff)
		if diffAbs < 1 || diffAbs > 3 {
			return false
		}

		// -- Check if change is consistent.
		var newChange Change
		if diff < 0 {
			newChange = Increasing
		} else if diff > 0 {
			newChange = Decreasing
		} else {
			return false
		}

		if change == Unset {
			change = newChange
		} else if change != newChange {
			return false
		}
	}

	return true
}

func isDampenSafe(report Report) bool {
	if isSafe(report) {
		return true
	}

	for i := range report {
		subReport := slices.Concat(report[:i], report[i+1:])
		if isSafe(subReport) {
			return true
		}
	}

	return false
}

func main() {
	var report Report
	numSafe := 0

	// -- Read input.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// -- Convert line to report.
		line := scanner.Text()
		nums := strings.Fields(line)
		report = make(Report, 0, len(nums))

		for _, num := range nums {
			num, err := strconv.Atoi(num)
			if err != nil {
				panic(err)
			}
			report = append(report, num)
		}

		// -- Count safe reports.
		if isDampenSafe(report) {
			numSafe += 1
		}
	}

	fmt.Println(numSafe)
}
