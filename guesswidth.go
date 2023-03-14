package guesswidth

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

func ToTable(lines []string, header int, trimSpace bool) [][]string {
	pos := widthPositions(lines, header, 4)
	return toTable(lines, pos, trimSpace)
}

func widthPositions(lines []string, header int, limit int) []int {
	var blanks []int
	for n, line := range lines {
		if n == header {
			blanks = lookupBlanks(strings.TrimSuffix(line, " "))
			continue
		}
		blanks = countBlanks(blanks, strings.TrimSuffix(line, " "))
	}
	return positions(blanks, limit)
}

func split(line string, pos []int, trimSpace bool) []string {
	n := 0
	start := 0
	columns := make([]string, len(pos)+1)
	lr := []rune(line)
	for p := 0; p < len(lr); p++ {
		if n > len(pos)-1 {
			columns[n] = strings.TrimSpace(string(lr[start:]))
			break
		}
		if pos[n] == p {
			for ; !unicode.IsSpace(lr[p]); p++ {
			}
			if trimSpace {
				columns[n] = strings.TrimSpace(string(lr[start:p]))
			} else {
				columns[n] = string(lr[start:p])
			}
			n++
			start = p
		}
		if runewidth.RuneWidth(lr[p]) == 2 {
			p++
		}
	}
	return columns
}

func toTable(lines []string, pos []int, trimSpace bool) [][]string {
	tables := make([][]string, 0, len(lines))
	for _, line := range lines {
		columns := split(line, pos, trimSpace)
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

func positions(blanks []int, limit int) []int {
	max := limit
	p := 0
	var pos []int
	for n, v := range blanks {
		if v >= max {
			max = v
			p = n
		}
		if v == 0 {
			max = limit
			if p > 0 {
				pos = append(pos, p)
				p = 0
			}
		}
	}
	return pos
}

func countPrint(line string, blanks []int) {
	fmt.Println(line)
	for _, k := range blanks {
		fmt.Print(k)
	}
	fmt.Println()
}
