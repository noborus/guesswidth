package guesswidth

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

type guessWidth struct {
	scanner    *bufio.Scanner
	preLines   []string
	preNum     int
	Header     int
	LimitSplit int
	TrimSpace  bool
}

var preNum = 10

func New(r io.Reader) *guessWidth {
	scanner := bufio.NewScanner(r)
	lines := preRead(scanner, preNum)
	return &guessWidth{
		scanner:    scanner,
		preLines:   lines,
		preNum:     preNum,
		Header:     0,
		LimitSplit: 0,
		TrimSpace:  true,
	}
}

func (g *guessWidth) Rows() [][]string {
	pos := widthPositions(g.preLines, g.Header)
	var rows [][]string
	if g.LimitSplit > 0 {
		if len(pos) > g.LimitSplit {
			pos = pos[:g.LimitSplit]
		}
		rows = toRows(g.preLines, pos, g.TrimSpace)
	} else {
		rows = toRows(g.preLines, pos, g.TrimSpace)
	}

	for g.scanner.Scan() {
		line := g.scanner.Text()
		columns := split(line, pos, g.TrimSpace)
		rows = append(rows, columns)
	}
	return rows
}

func preRead(scanner *bufio.Scanner, preRead int) []string {
	var lines []string
	for i := 0; i < preRead; i++ {
		if !scanner.Scan() {
			return lines
		}
		lines = append(lines, scanner.Text())
	}
	return lines
}

func ToTable(lines []string, header int, trimSpace bool) [][]string {
	pos := widthPositions(lines, header)
	return toRows(lines, pos, trimSpace)
}

func ToTableN(lines []string, header int, numSplit int, trimSpace bool) [][]string {
	pos := widthPositions(lines, header)
	if len(pos) > numSplit {
		pos = pos[:numSplit]
	}
	return toRows(lines, pos, trimSpace)
}

func widthPositions(lines []string, header int) []int {
	var blanks []int
	limit := 2
	if header < 0 {
		header = 0
	}
	for n, line := range lines {
		if n < header {
			continue
		}
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
			start = p
			break
		}
		if pos[n] == p {
			end := p
			for ; !unicode.IsSpace(lr[end]); end++ {
			}
			if trimSpace {
				columns[n] = strings.TrimSpace(string(lr[start:end]))
			} else {
				columns[n] = string(lr[start:end])
			}
			n++
			start = p
		}
		if runewidth.RuneWidth(lr[p]) == 2 {
			p++
		}
	}
	columns[len(columns)-1] = strings.TrimSpace(string(lr[start:]))
	return columns
}

func toRows(lines []string, pos []int, trimSpace bool) [][]string {
	rows := make([][]string, 0, len(lines))
	for _, line := range lines {
		columns := split(line, pos, trimSpace)
		rows = append(rows, columns)
	}
	return rows
}

// Creates a blank(1) and non-blank(0) slice.
// Execute for the base line (header line).
func lookupBlanks(line string) []int {
	blanks := make([]int, 0)
	for _, v := range line {
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

// Count up if the line is blank where the reference line was blank.
func countBlanks(blanks []int, line string) []int {
	n := 0
	for _, r := range line {
		if r == ' ' && blanks[n] > 0 {
			blanks[n] += 1
		}

		n += 1
		if runewidth.RuneWidth(r) == 2 {
			n += 1
		}
		if n >= len(blanks) {
			break
		}
	}
	return blanks
}

// Generates a list of separator positions from a blank slice.
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
