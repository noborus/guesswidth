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
				{6, 1, 1},
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
				{6, 1, 1},
				{5, 0, 0},
				{12, 1, 1},
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
