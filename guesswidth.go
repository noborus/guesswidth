package guesswidth

import (
	"fmt"
	"strings"

	"github.com/mattn/go-runewidth"
)

func ToTable(lines []string, header int, trimSpace bool) [][]string {
	pos := widthPosition(lines, header)
	return toTable(lines, pos, trimSpace)
}

func widthPosition(lines []string, header int) []int {
	var blanks []int
	for n, line := range lines {
		if n == header {
			blanks = lookupBlanks(strings.TrimSuffix(line, " "))
			continue
		}
		blanks = countBlanks(blanks, strings.TrimSuffix(line, " "))
	}
	return positions(blanks)
}

func toTable(lines []string, pos []int, trimSpace bool) [][]string {
	tables := make([][]string, 0, len(lines))
	for _, line := range lines {
		p := 0
		n := 0
		start := 0
		columns := make([]string, len(pos)+1)
		for _, c := range line {
			if n > len(pos)-1 {
				columns[n] = strings.TrimSpace(line[start:])
				break
			}
			if pos[n] == p {
				for ; line[p] != ' '; p++ {
				}
				if trimSpace {
					columns[n] = strings.TrimSpace(line[start:p])
				} else {
					columns[n] = line[start:p]
				}
				n++
				start = p + 1
			}
			p++
			if runewidth.RuneWidth(c) == 2 {
				p++
			}
		}
		tables = append(tables, columns)
	}
	return tables
}

func lookupBlanks(str string) []int {
	blanks := make([]int, 0)
	for _, v := range str {
		if v == ' ' {
			blanks = append(blanks, 1)
		} else {
			blanks = append(blanks, 0)
		}
		if runewidth.RuneWidth(v) == 2 {
			blanks = append(blanks, 0)
		}
	}
	return blanks
}

func countBlanks(blanks []int, str string) []int {
	width := 0
	for _, v := range str {
		if v == ' ' && blanks[width] > 0 {
			blanks[width] += 1
		}

		width += 1
		if runewidth.RuneWidth(v) == 2 {
			width += 1
		}
		if width >= len(blanks) {
			break
		}
	}
	return blanks
}

func positions(blanks []int) []int {
	max := 2
	p := 0
	var pos []int
	for n, v := range blanks {
		if v >= max {
			max = v
			p = n
		}
		if v == 0 {
			max = 2
			if p > 0 {
				pos = append(pos, p)
				p = 0
			}
		}
	}
	return pos
}

func countPrint(header string, blanks []int) {
	fmt.Print(string(header))
	for _, k := range blanks {
		fmt.Print(k)
	}
	fmt.Println()
}
