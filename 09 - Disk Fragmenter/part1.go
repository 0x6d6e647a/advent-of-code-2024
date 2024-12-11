package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type filesystem []int

func newFilesystem(line string) filesystem {
	// -- Convert string to disk map.
	diskMap := make([]int, 0, len(line))
	fsSize := 0

	for _, ch := range line {
		// -- Convert character to integer.
		blockSize, err := strconv.Atoi(string(ch))
		if err != nil {
			panic(err)
		}

		diskMap = append(diskMap, blockSize)
		fsSize += blockSize
	}

	// -- Convert disk map to filesystem.
	filesystem := make([]int, 0, fsSize)
	fsIndex := 0

	for chIndex, blockSize := range diskMap {
		// -- Use filesystem index or empty indicator.
		isEmpty := chIndex%2 != 0
		currIndex := -1
		if !isEmpty {
			currIndex = fsIndex
		}

		// -- Add blocks to file system.
		for range blockSize {
			filesystem = append(filesystem, currIndex)
		}

		// -- Increment the filesystem index.
		if isEmpty {
			fsIndex += 1
		}
	}

	return filesystem
}

func (fs filesystem) checksum() int {
	checksum := 0

	for index, blockIndex := range fs {
		if blockIndex == -1 {
			continue
		}
		checksum += index * blockIndex
	}

	return checksum
}

func (fs *filesystem) compress() {
	dst := 0 

	for src := len(*fs)-1; src >= 0; src -= 1 {
		if (*fs)[src] == -1 {
			continue
		}

		for ; dst < len(*fs); dst += 1 {
			if (*fs)[dst] == -1 {
				break
			}
		}

		if dst >= len(*fs) || dst >= src {
			break
		}

		(*fs)[dst] = (*fs)[src]
		(*fs)[src] = -1
	}
}

func main() {
	// -- Get input line.
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		panic("no input")
	}
	line := scanner.Text()
	fs := newFilesystem(line)
	fs.compress()
	checksum := fs.checksum()
	fmt.Println(checksum)
}
