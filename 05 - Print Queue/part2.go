package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"slices"
	"strconv"
	"strings"
)

func getIndices[T comparable](slice []T, value T) []int {
	var indices []int

	for index, curr := range slice {
		if curr == value {
			indices = append(indices, index)
		}
	}

	return indices
}

type pageOrderRule struct {
	before int
	after  int
}

func newPageOrderRule(line string) pageOrderRule {
	nums := strings.Split(line, "|")
	before, err := strconv.Atoi(nums[0])
	if err != nil {
		panic(err)
	}
	after, err := strconv.Atoi(nums[1])
	if err != nil {
		panic(err)
	}
	return pageOrderRule{before, after}
}

func (p pageOrderRule) isSatisfied(pageNums []int) bool {
	beforeIndex := slices.Index(pageNums, p.before)
	if beforeIndex == -1 {
		return true
	}

	afterIndex := slices.Index(pageNums, p.after)
	if afterIndex == -1 {
		return true
	}

	return beforeIndex < afterIndex
}

type pageOrderRules []pageOrderRule

func newPageOrderRules(scanner *bufio.Scanner) pageOrderRules {
	var rules pageOrderRules

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			break
		}

		rule := newPageOrderRule(line)
		rules = append(rules, rule)
	}

	return rules
}

func (p pageOrderRules) isSatisfied(pageNums []int) bool {
	for _, rule := range p {
		if !rule.isSatisfied(pageNums) {
			return false
		}
	}
	return true
}

func (p pageOrderRules) getInvalids(pageNumss [][]int) [][]int {
	var invalids [][]int

	for _, pageNums := range pageNumss {
		if !p.isSatisfied(pageNums) {
			invalids = append(invalids, pageNums)
			continue
		}
	}

	return invalids
}

func (p pageOrderRules) getBrokenRules(pageNums []int) pageOrderRules {
	var broken pageOrderRules

	for _, rule := range p {
		if !rule.isSatisfied(pageNums) {
			broken = append(broken, rule)
		}
	}

	return broken
}

func (p pageOrderRules) fixOrdering(pageNums []int) {
	for !p.isSatisfied(pageNums) {
		broken := p.getBrokenRules(pageNums)

		for _, rule := range broken {
			beforeIndices := getIndices(pageNums, rule.before)
			beforeIndex := beforeIndices[len(beforeIndices)-1]
			afterIndices := getIndices(pageNums, rule.after)
			afterIndex := afterIndices[0]

			pageNums[beforeIndex], pageNums[afterIndex] = pageNums[afterIndex], pageNums[beforeIndex]
		}

		// -- Shuffle rules to prevent deadlock.
		rand.Shuffle(len(p), func(i int, j int) {
			p[i], p[j] = p[j], p[i]
		})
	}
}

func parseUpdates(scanner *bufio.Scanner) [][]int {
	var updates [][]int

	for scanner.Scan() {
		line := scanner.Text()
		strs := strings.Split(line, ",")
		update := make([]int, 0, len(strs))

		for _, str := range strs {
			num, err := strconv.Atoi(str)
			if err != nil {
				panic(err)
			}

			update = append(update, num)
		}

		updates = append(updates, update)
	}

	return updates
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	rules := newPageOrderRules(scanner)
	updates := parseUpdates(scanner)
	invalids := rules.getInvalids(updates)

	// -- Fix invalid update ordering.
	sum := 0

	for _, update := range invalids {
		rules.fixOrdering(update)
		middle := len(update) / 2
		sum += update[middle]
	}

	fmt.Println(sum)
}
