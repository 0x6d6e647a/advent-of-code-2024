package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

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

func (p pageOrderRules) getValids(pageNumss [][]int) [][]int {
	var valids [][]int

	for _, pageNums := range pageNumss {
		if !p.isSatisfied(pageNums) {
			continue
		}

		valids = append(valids, pageNums)
	}

	return valids
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
	valids := rules.getValids(updates)

	// -- Check validity of updates.
	sum := 0

	for _, valid := range valids {
		middle := len(valid) / 2
		sum += valid[middle]
	}

	fmt.Println(sum)
}
