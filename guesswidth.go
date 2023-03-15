package guesswidth

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

type GuessWidth struct {
	reader *bufio.Reader
	// pos is a list of separator positions.
	pos []int
	// preLines stores the lines read for scan.
	preLines []string
	// preCount is the number returned by read.
	preCount int

	// ScanNum is the number to scan to analyze.
	ScanNum int
	// Header is the base line number. It starts from 0.
	Header int
	// limitSplit is the maximum number of columns to split.
	LimitSplit int
	// TrimSpace is whether to trim the space in the value.
	TrimSpace bool
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *GuessWidth {
	reader := bufio.NewReader(r)
	g := &GuessWidth{
		reader:     reader,
		ScanNum:    100,
		preCount:   0,
		Header:     0,
		LimitSplit: 0,
		TrimSpace:  true,
	}
	return g
}

// ReadAll reads all rows
// and returns a two-dimensional slice of rows and columns.
func (g *GuessWidth) ReadAll() [][]string {
	g.Scan(g.ScanNum)

	var rows [][]string
	for {
		columns, err := g.Read()
		if err != nil {
			break
		}
		rows = append(rows, columns)
	}
	return rows
}

// Scan preReads and parses the lines.
func (g *GuessWidth) Scan(num int) {
	for i := 0; i < num; i++ {
		buf, _, err := g.reader.ReadLine()
		if err != nil {
			break
		}
		g.preLines = append(g.preLines, string(buf))
	}

	g.pos = widthPositions(g.preLines, g.Header)

	if g.LimitSplit > 0 {
		if len(g.pos) > g.LimitSplit {
			g.pos = g.pos[:g.LimitSplit]
		}
	}
}

// Read reads one row and returns a slice of columns.
// Scan is executed first if it is not preRead.
func (g *GuessWidth) Read() ([]string, error) {
	if len(g.preLines) == 0 {
		g.Scan(g.ScanNum)
	}

	var line string
	if g.preCount < len(g.preLines) {
		line = g.preLines[g.preCount]
		g.preCount++
	} else {
		buf, _, err := g.reader.ReadLine()
		if err != nil {
			return nil, err
		}
		line = string(buf)
	}

	return split(line, g.pos, g.TrimSpace), nil
}

// ToTable parses a slice of rows and returns a table.
func ToTable(lines []string, header int, trimSpace bool) [][]string {
	pos := widthPositions(lines, header)
	return toRows(lines, pos, trimSpace)
}

// ToTableN parses a slice of rows and returns a table, but limits the number of splits.
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

func separatorPosition(lr []rune, p int, pos []int, n int) int {
	if unicode.IsSpace(lr[p]) {
		return p
	}

	f := p
	fp := 0
	for ; !unicode.IsSpace(lr[f]); f++ {
		fp++
	}

	b := p
	bp := 0
	for ; !unicode.IsSpace(lr[b]); b-- {
		bp++
	}

	if b == pos[n] {
		return f
	}
	if n < len(pos)-1 {
		if f == pos[n+1] {
			return b
		}
		if b == pos[n] {
			return f
		}
		if b > pos[n] && b < pos[n+1] {
			return b
		}
	}
	return f
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
			p = separatorPosition(lr, p, pos, n)
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

// debugCountPrint is for debugging which prints the space count.
func debugCountPrint(line string, blanks []int) {
	fmt.Println(line)
	for _, k := range blanks {
		fmt.Print(k)
	}
	fmt.Println()
}
