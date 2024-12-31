package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type set[T comparable] map[T]struct{}

func newSet[T comparable]() set[T] {
	return make(map[T]struct{}, 0)
}

func (s set[T]) insert(value T) {
	s[value] = struct{}{}
}

func permutations[T any](arr []T) (res [][]T) {
	var perm func([]T, int)
	perm = func(arr []T, n int) {
		if n == 1 {
			tmp := make([]T, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i += 1 {
				perm(arr, n-1)

				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	perm(arr, len(arr))
	return res
}

func permutationsString(s string) (res []string) {
	perms := permutations([]byte(s))

	for _, perm := range perms {
		res = append(res, string(perm))
	}

	return res
}

type coord struct {
	row int
	col int
}

var dirs = map[byte]coord{
	'^': {-1, 0},
	'v': {+1, 0},
	'>': {0, +1},
	'<': {0, -1},
}

var numpad = map[byte]coord{
	'7': {0, 0}, '8': {0, 1}, '9': {0, 2},
	'4': {1, 0}, '5': {1, 1}, '6': {1, 2},
	'1': {2, 0}, '2': {2, 1}, '3': {2, 2},
	' ': {3, 0}, '0': {3, 1}, 'A': {3, 2},
}

func validNumpadPos(pos coord) bool {
	lim := numpad['A']
	return pos.row >= 0 &&
		pos.col >= 0 &&
		pos.row <= lim.row &&
		pos.col <= lim.col &&
		pos != numpad[' ']
}

var dirpad = map[byte]coord{
	' ': {0, 0}, '^': {0, 1}, 'A': {0, 2},
	'<': {1, 0}, 'v': {1, 1}, '>': {1, 2},
}

func validDirpadPos(pos coord) bool {
	lim := dirpad['>']
	return pos.row >= 0 &&
		pos.col >= 0 &&
		pos.row <= lim.row &&
		pos.col <= lim.col &&
		pos != dirpad[' ']
}

func validKeypadPos(pos coord, isDirpad bool) bool {
	if isDirpad {
		return validDirpadPos(pos)
	}
	return validNumpadPos(pos)
}

func getKeyPresses(start coord, end coord) string {
	var sb strings.Builder

	diffRow := end.row - start.row
	diffCol := end.col - start.col

	for range diffRow {
		sb.WriteByte('v')
	}

	for range -diffRow {
		sb.WriteByte('^')
	}

	for range diffCol {
		sb.WriteByte('>')
	}

	for range -diffCol {
		sb.WriteByte('<')
	}

	return sb.String()
}

type node struct {
	sequence string
	depth    int
	isDirpad bool
	pos      coord
}

var cache = make(map[node]int)

func countPresses(sequence string, depth int, isDirpad bool, pos coord) int {
	// -- Return empty if done.
	if len(sequence) == 0 {
		return 0
	}

	// -- Check if cached result exists.
	node := node{sequence, depth, isDirpad, pos}
	cached, ok := cache[node]
	if ok {
		return cached
	}

	// -- Get current keypad.
	keypad := numpad
	if isDirpad {
		keypad = dirpad
	}

	press := keypad[sequence[0]]
	buttons := getKeyPresses(pos, press)

	// -- Get current minimum length.
	var min_len int

	if depth != 0 {
		var perm_lengths []int

		perms := permutationsString(buttons)
		perm_set := newSet[string]()
		for _, perm := range perms {
			perm_set.insert(perm)
		}

		for perm := range perm_set {
			curr := pos
			invalid := false

			for _, button := range perm {
				delta := dirs[byte(button)]
				curr.row += delta.row
				curr.col += delta.col

				if !validKeypadPos(curr, isDirpad) {
					invalid = true
					break
				}
			}

			if !invalid {
				perm += "A"
				perm_lengths = append(perm_lengths,
					countPresses(perm, depth-1, true, dirpad['A']))

			}
		}

		if len(perm_lengths) == 0 {
			perm_lengths = append(perm_lengths, countPresses("A", depth-1, true, dirpad['A']))
		}

		min_len = slices.Min(perm_lengths)
	} else {
		min_len = len(buttons) + 1
	}

	result := min_len + countPresses(sequence[1:], depth, isDirpad, press)
	cache[node] = result
	return result
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		codeNumStr := strings.TrimSuffix(line, "A")
		codeNum, err := strconv.Atoi(codeNumStr)
		if err != nil {
			panic(err)
		}
		presses := countPresses(line, 25, false, numpad['A'])
		total += codeNum * presses
	}

	fmt.Println(total)
}
