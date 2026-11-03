package guesswidth

import (
	"bufio"
	"reflect"
	"strings"
	"testing"
)

func TestGuessWidth_ReadAll(t *testing.T) {
	type fields struct {
		reader     *bufio.Reader
		pos        []int
		preLines   []string
		preCount   int
		ScanNum    int
		Header     int
		LimitSplit int
		MinLines   int
		TrimSpace  bool
	}
	tests := []struct {
		name   string
		fields fields
		want   [][]string
		want2  []Cols
	}{
		{
			name: "ps",
			fields: fields{
				reader: bufio.NewReader(strings.NewReader(`   PID TTY          TIME CMD
302965 pts/3    00:00:11 zsh
709737 pts/3    00:00:00 ps`)),
				ScanNum:  100,
				MinLines: 2,
			},
			want: [][]string{
				{"   PID", " TTY     ", "     TIME", "CMD"},
				{"302965", " pts/3   ", " 00:00:11", "zsh"},
				{"709737", " pts/3   ", " 00:00:00", "ps"},
			},
			want2: []Cols{
				{6, 0, 1},
				{5, 0, 0},
				{8, 1, 3},
				{3, 0, 0},
			},
		},
		{
			name: "ps overflow",
			fields: fields{
				reader: bufio.NewReader(strings.NewReader(`USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root           1  0.0  0.0 168576 13788 ?        Ss   Mar11   0:49 /sbin/init splash
noborus   703052  2.1  0.7 1184814400 230920 ?   Sl   10:03   0:45 /opt/google/chrome/chrome
noborus   721971  0.0  0.0  13716  3524 pts/3    R+   10:39   0:00 ps aux`)),
				ScanNum:  100,
				MinLines: 2,
			},
			want: [][]string{
				{"USER     ", "    PID", " %CPU", " %MEM", "    VSZ", "   RSS", " TTY     ", " STAT", " START  ", " TIME", "COMMAND"},
				{"root     ", "      1", "  0.0", "  0.0", " 168576", " 13788", " ?       ", " Ss  ", " Mar11  ", " 0:49", "/sbin/init splash"},
				{"noborus  ", " 703052", "  2.1", "  0.7", " 1184814400", " 230920", " ?  ", " Sl  ", " 10:03  ", " 0:45", "/opt/google/chrome/chrome"},
				{"noborus  ", " 721971", "  0.0", "  0.0", "  13716", "  3524", " pts/3   ", " R+  ", " 10:39  ", " 0:00", "ps aux"},
			},
			want2: []Cols{
				{7, 0, 0},
				{6, 1, 4},
				{4, 1, 4},
				{4, 1, 4},
				{10, 1, 4},
				{6, 1, 4},
				{5, 0, 0},
				{4, 0, 1},
				{5, 0, 0},
				{4, 1, 4},
				{25, 0, 0},
			},
		},
		{
			name: "ps limit",
			fields: fields{
				reader: bufio.NewReader(strings.NewReader(`   PID TTY          TIME CMD
302965 pts/3    00:00:11 zsh
709737 pts/3    00:00:00 ps`)),
				ScanNum:    100,
				MinLines:   2,
				LimitSplit: 2,
			},
			want: [][]string{
				{"   PID", " TTY     ", "    TIME CMD"},
				{"302965", " pts/3   ", "00:00:11 zsh"},
				{"709737", " pts/3   ", "00:00:00 ps"},
			},
			want2: []Cols{
				{6, 0, 1},
				{5, 0, 0},
				{12, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuessWidth{
				reader:     tt.fields.reader,
				pos:        tt.fields.pos,
				preLines:   tt.fields.preLines,
				preCount:   tt.fields.preCount,
				ScanNum:    tt.fields.ScanNum,
				Header:     tt.fields.Header,
				LimitSplit: tt.fields.LimitSplit,
				MinLines:   tt.fields.MinLines,
				TrimSpace:  tt.fields.TrimSpace,
			}
			if got := g.ReadAll(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GuessWidth.ReadAll() = \n%#v, want \n%#v", got, tt.want)
			}
			if got2 := g.Widths; !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("GuessWidth.ReadAll() = \n%v, want \n%v", got2, tt.want2)
			}
		})
	}
}

func TestGuessWidth_ReadAllLongLines(t *testing.T) {
	prefix := strings.Repeat("x", 8192)
	input := prefix + "  H\n" + prefix + "  1\n" + prefix + "  2\n"
	g := NewReader(strings.NewReader(input))
	g.ScanNum = 2

	want := [][]string{
		{prefix, "H"},
		{prefix, "1"},
		{prefix, "2"},
	}
	if got := g.ReadAll(); !reflect.DeepEqual(got, want) {
		t.Errorf("GuessWidth.ReadAll() = %#v, want %#v", got, want)
	}
}

func TestToTable_LastHeaderColumn(t *testing.T) {
	lines := []string{
		"OWNER         ID     STATUS DETAILS",
		"noborus      16101   ready /Applications/Ghostty.app/Contents/MacOS/ghostty",
		"root         1       stopped -/bin/zsh",
	}

	got := ToTable(lines, 0, true)
	if len(got) != 3 {
		t.Fatalf("len(ToTable) = %d, want 3", len(got))
	}
	if len(got[1]) != 4 {
		t.Fatalf("len(row1) = %d, want 4: %#v", len(got[1]), got[1])
	}
	if got[1][2] != "ready" || got[1][3] != "/Applications/Ghostty.app/Contents/MacOS/ghostty" {
		t.Fatalf("split row1 = %#v, want STATUS and DETAILS separated", got[1])
	}
	if got[2][2] != "stopped" || got[2][3] != "-/bin/zsh" {
		t.Fatalf("split row2 = %#v, want STATUS and DETAILS separated", got[2])
	}
}

func TestToTable_ShiftedColumnWithDoubleSpace(t *testing.T) {
	lines := []string{
		"USER               PID  %CPU %MEM      VSZ    RSS   TT  STAT STARTED      TIME COMMAND",
		"_mds_stores        845  13.8  2.3 444217440 579536   ??  Ss   24Aug26 124:52.34 mds_stores",
		"_windowserver      595   9.3  0.6 436109072 148256   ??  Ss   24Aug26 1910:22.18 WindowServer",
		"noborusai        42356   2.8  0.1 435375200  25984   ??  S    10:14AM   0:00.11 mdworker_shared",
		"noborusai        16101   2.8  0.8 437316592 204272   ??  S    Wed06AM   1:06.83 ghostty",
		"noborusai         2141   0.8  1.0 441002176 253392   ??  Ss   24Aug26   2:35.70 JapaneseIM",
		"noborusai        40973   0.7  0.0 435309392   6080 s000  S     7:26AM   0:01.06 zsh",
		"noborusai        40462   0.4  9.4 1951534752 2376992   ??  S    11:11PM  14:37.53 Chrome",
	}

	got := ToTable(lines, 0, true)
	if got[7][6] != "??" || got[7][7] != "S" {
		t.Fatalf("split row = %#v at positions %v, want TT and STAT separated", got[7], Positions(lines, 0, 2))
	}
}

func TestToTable_DoubleSpaceInCommand(t *testing.T) {
	command1 := "printf 'a  b'"
	command2 := "echo hello"
	commandWidth := 24
	lines := []string{
		"ID  COMMAND" + strings.Repeat(" ", 19) + "STATUS",
		"1   " + command1 + strings.Repeat(" ", commandWidth-len(command1)) + "  RUNNING",
		"2   " + command2 + strings.Repeat(" ", commandWidth-len(command2)) + "  STOPPED",
	}

	got := ToTable(lines, 0, true)
	want := [][]string{
		{"ID", "COMMAND", "STATUS"},
		{"1", command1, "RUNNING"},
		{"2", command2, "STOPPED"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToTable() = %#v, want %#v", got, want)
	}
}

func TestToTable(t *testing.T) {
	type args struct {
		lines     []string
		header    int
		trimSpace bool
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "PS",
			args: args{
				lines: []string{
					"    PID TTY          TIME CMD",
					" 302965 pts/3    00:00:00 zsh",
					" 323990 pts/3    00:00:00 ps",
				},
				header:    0,
				trimSpace: true,
			},
			want: [][]string{
				{"PID", "TTY", "TIME", "CMD"},
				{"302965", "pts/3", "00:00:00", "zsh"},
				{"323990", "pts/3", "00:00:00", "ps"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTable(tt.args.lines, tt.args.header, tt.args.trimSpace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToTableN(t *testing.T) {
	type args struct {
		lines     []string
		header    int
		numSplit  int
		trimSpace bool
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "no header",
			args: args{
				lines: []string{
					"2022-12-21T09:50:16+0000 WARN A warning that should be ignored is usually at this level and should be actionable.",
					"2022-12-21T09:50:17+0000 INFO This is less important than debug log and is often used to provide context in the current task.",
				},
				numSplit: 2,
			},
			want: [][]string{
				{"2022-12-21T09:50:16+0000", " WARN", "A warning that should be ignored is usually at this level and should be actionable."},
				{"2022-12-21T09:50:17+0000", " INFO", "This is less important than debug log and is often used to provide context in the current task."},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToTableN(tt.args.lines, tt.args.header, tt.args.numSplit, tt.args.trimSpace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToTableN() = \n%#v, want \n%#v", got, tt.want)
			}
		})
	}
}

func TestGuessWidth_UpdateMaxWidth(t *testing.T) {
	g := &GuessWidth{Widths: make([]Cols, 2)}

	got := g.UpdateMaxWidth([]string{" 12", "abc"})
	want := []Cols{{Width: 2, rightCount: 1}, {Width: 3}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("UpdateMaxWidth() = %v, want %v", got, want)
	}
}

func TestGuessWidth_SetJustified(t *testing.T) {
	tests := []struct {
		name      string
		threshold int
		want      []Cols
	}{
		{
			name:      "threshold one",
			threshold: 1,
			want:      []Cols{{Justified: Left}, {Justified: Left, rightCount: 1}, {Justified: Right, rightCount: 2}},
		},
		{
			name:      "threshold two",
			threshold: 2,
			want:      []Cols{{Justified: Left}, {Justified: Left, rightCount: 1}, {Justified: Left, rightCount: 2}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GuessWidth{Widths: []Cols{{}, {rightCount: 1}, {rightCount: 2}}}

			if got := g.SetJustified(tt.threshold); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetJustified(%d) = %v, want %v", tt.threshold, got, tt.want)
			}
		})
	}
}

func TestGuessWidth_Scan(t *testing.T) {
	g := NewReader(strings.NewReader("A  B\n1  2\n"))
	g.Scan(2)

	if !reflect.DeepEqual(g.pos, []int{2}) {
		t.Errorf("Scan() positions = %v, want [2]", g.pos)
	}
	if !reflect.DeepEqual(g.preLines, []string{"A  B", "1  2"}) {
		t.Errorf("Scan() lines = %v, want [A  B 1  2]", g.preLines)
	}
}

func TestPositions(t *testing.T) {
	lines := []string{"ignored", "A  B", "1  2"}
	got := Positions(lines, 1, 2)
	if !reflect.DeepEqual(got, []int{2}) {
		t.Errorf("Positions() = %v, want [2]", got)
	}
}

func TestSplitGraphemes(t *testing.T) {
	tests := []struct {
		name string
		line string
		pos  []int
		trim bool
		want []string
	}{
		{
			name: "emoji grapheme",
			line: "👩‍👩‍👧‍👦  next",
			pos:  []int{2},
			trim: false,
			want: []string{"👩‍👩‍👧‍👦", " next"},
		},
		{
			name: "combining mark",
			line: "e\u0301  next",
			pos:  []int{2},
			trim: true,
			want: []string{"e\u0301", "next"},
		},
		{
			name: "no separator",
			line: "hello",
			trim: true,
			want: []string{"hello"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.line, tt.pos, tt.trim); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("split() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	if got := Version(); got == "" {
		t.Error("Version() returned an empty string")
	}
}
