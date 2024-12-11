package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
)

type diskSpace struct {
	fileId int
	size   int
	prev   *diskSpace
	next   *diskSpace
}

func (ds *diskSpace) insert(value diskSpace) {
	if ds.fileId != -1 {
		panic("attempt to insert into non-empty disk space")
	}

	if ds.size-value.size < 0 {
		panic("attempt to insert with insufficient space")
	}

	value.prev = ds.prev
	value.next = ds
	ds.prev.next = &value
	ds.prev = &value

	ds.size -= value.size

	if ds.size == 0 {
		ds.prev.next = ds.next
		ds.next.prev = ds.prev
	}
}

func (ds *diskSpace) remove() *diskSpace {
	ds.fileId = -1

	if ds.prev.fileId == -1 {
		ds.size += ds.prev.size
		ds.prev.prev.next = ds
		ds.prev = ds.prev.prev
	}

	if ds.next != nil && ds.next.fileId == -1 {
		ds.size += ds.next.size
		if ds.next.next != nil {
			ds.next.next.prev = ds
		}
		ds.next = ds.next.next
	}

	return ds
}

type filesystem struct {
	sentinel diskSpace
	numFiles int
	last     *diskSpace
}

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
	sentinel := diskSpace{math.MinInt, math.MinInt, nil, nil}
	numFiles := 0
	curr := &sentinel

	for index, size := range diskMap {
		// -- Determine file id to use.
		isEmpty := index%2 != 0
		fileId := -1
		if !isEmpty {
			fileId = numFiles
			numFiles += 1
		}

		// -- Skip fully empty spots.
		if size == 0 {
			continue
		}

		// -- Create disk space.
		ds := &diskSpace{fileId, size, curr, nil}
		curr.next = ds
		curr = ds
	}

	return filesystem{sentinel, numFiles, curr}
}

func (fs *filesystem) compress() {
	src := fs.last

eachFileId:
	for srcIndex := fs.numFiles - 1; srcIndex >= 0; srcIndex -= 1 {
		// -- Find source.
		for src.fileId != srcIndex {
			src = src.prev
		}

		// -- Find destination.
		dst := fs.sentinel.next
		for dst != nil {
			if dst == src {
				continue eachFileId
			}

			if dst.fileId == -1 && dst.size >= src.size {
				break
			}

			dst = dst.next
		}

		// -- Insert source at destination.
		dst.insert(*src)
		src = src.remove()
	}
}

func (fs filesystem) checksum() int {
	checksum := 0
	index := 0
	curr := fs.sentinel.next
	for curr != nil {
		for range curr.size {
			if curr.fileId != -1 {
				checksum += index * curr.fileId
			}
			index += 1
		}

		curr = curr.next
	}

	return checksum
}

func main() {
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
