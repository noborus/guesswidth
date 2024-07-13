// Package guesswidth handles the format as formatted by printf.
// Spaces exist as delimiters, but spaces are not always delimiters.
// The width seems to be a fixed length, but it doesn't always fit.
// guesswidth finds the column separation position
// from the reference line(header) and multiple lines(body).
package guesswidth

import (
	"bufio"
	"fmt"
	"io"
	"runtime/debug"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
)

// GuessWidth reads records from printf-like output.
type GuessWidth struct {
	reader *bufio.Reader
	// pos is a list of separator positions.
	pos []int
	// Widths is the width of the column.
	Widths []Cols
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
	// MinLines is the minimum number of lines to recognize as a separator.
	// 1 if only the header, 2 or more if there is a blank in the body.
	MinLines int
	// TrimSpace is whether to trim the space in the value.
	TrimSpace bool
}

type Cols struct {
	Width      int
	Justified  int
	rightCount int
}

const (
	Left = iota
	Right
)

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *GuessWidth {
	reader := bufio.NewReader(r)
	g := &GuessWidth{
		reader:     reader,
		ScanNum:    100,
		preCount:   0,
		Header:     0,
		LimitSplit: 0,
		MinLines:   2,
		TrimSpace:  true,
	}
	return g
}

// ReadAll reads all rows
// and returns a two-dimensional slice of rows and columns.
func (g *GuessWidth) ReadAll() [][]string {
	if len(g.preLines) == 0 {
		g.Scan(g.ScanNum)
	}

	g.Widths = make([]Cols, len(g.pos)+1)
	var rows [][]string
	for {
		columns, err := g.Read()
		if err != nil {
			break
		}
		g.UpdateMaxWidth(columns)
		rows = append(rows, columns)
	}

	g.SetJustified(len(rows) / 2)
	return rows
}

// UpdateMaxWidth updates the maximum width of the column.
func (g *GuessWidth) UpdateMaxWidth(columns []string) []Cols {
	if len(g.Widths) < len(columns) {
		for n := len(g.Widths); n < len(columns); n++ {
			g.Widths = append(g.Widths, Cols{})
		}
	}

	for n, col := range columns {
		width := runewidth.StringWidth(col)
		if width > g.Widths[n].Width {
			g.Widths[n].Width = width
		}
		if isRightAlign(col) {
			g.Widths[n].rightCount++
		}
	}
	return g.Widths
}

// SetJustified sets the justification of the column.
func (g *GuessWidth) SetJustified(threshold int) []Cols {
	for n, col := range g.Widths {
		if col.rightCount < threshold {
			col.Justified = Left
		} else {
			col.Justified = Right
		}
		g.Widths[n] = col
	}
	return g.Widths
}

func isRightAlign(str string) bool {
	if str == "" {
		return false
	}
	for n := 0; n < len(str); n++ {
		if str[n] != ' ' {
			return false
		}
		if str[len(str)-n-1] != ' ' {
			return true
		}
	}
	return false
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

	g.pos = Positions(g.preLines, g.Header, g.MinLines)

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

// ToTable parses a slice of lines and returns a table.
func ToTable(lines []string, header int, trimSpace bool) [][]string {
	pos := Positions(lines, header, 2)
	return toRows(lines, pos, trimSpace)
}

// ToTableN parses a slice of lines and returns a table, but limits the number of splits.
func ToTableN(lines []string, header int, numSplit int, trimSpace bool) [][]string {
	pos := Positions(lines, header, 2)
	if len(pos) > numSplit {
		pos = pos[:numSplit]
	}
	return toRows(lines, pos, trimSpace)
}

// Positions returns separator positions
// from multiple lines and header line number.
// Lines before the header line are ignored.
func Positions(lines []string, header int, minLines int) []int {
	var blanks []int

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
	return positions(blanks, minLines)
}

func separatorPosition(lr []rune, p int, pos []int, n int) int {
	if unicode.IsSpace(lr[p]) {
		return p
	}

	f := p
	fp := 0

	for ; f < len(lr) && !unicode.IsSpace(lr[f]); f++ {
		fp++
	}

	b := p
	bp := 0
	for ; b > 0 && !unicode.IsSpace(lr[b]); b-- {
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
	w := 0
	for p := 0; p < len(lr); p++ {
		if n > len(pos)-1 {
			start = p
			break
		}
		if pos[n] <= w {
			end := separatorPosition(lr, p, pos, n)
			if start > end {
				break
			}
			col := string(lr[start:end])
			if trimSpace {
				columns[n] = strings.TrimSpace(col)
			} else {
				columns[n] = col
			}
			n++
			start = end
		}
		w += runewidth.RuneWidth(lr[p])
	}
	if n < len(columns) {
		col := string(lr[start:])
		if trimSpace {
			columns[n] = strings.TrimSpace(col)
		} else {
			columns[n] = col
		}
	}
	return columns
}

// roRows returns rows separated by columns.
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
	first := true
	for _, v := range line {
		if v == ' ' {
			if first {
				blanks = append(blanks, 0)
				continue
			}
			blanks = append(blanks, 1)
			continue
		}

		first = false
		blanks = append(blanks, 0)
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
		if n >= len(blanks) {
			break
		}
		if r == ' ' && blanks[n] > 0 {
			blanks[n] += 1
		}

		n++
		if runewidth.RuneWidth(r) == 2 {
			n++
		}
	}
	return blanks
}

// Generates a list of separator positions from a blank slice.
func positions(blanks []int, minLines int) []int {
	max := minLines
	p := 0
	var pos []int
	for n, v := range blanks {
		if v >= max {
			max = v
			p = n
		}
		if v == 0 {
			max = minLines
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

var (
	version  string
	revision string
)

func Version() string {
	if version != "" {
		return version + " rev:" + revision
	}
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "(devel)"
	}
	return info.Main.Version
}
