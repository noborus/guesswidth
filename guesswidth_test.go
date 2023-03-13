package guesswidth

import (
	"reflect"
	"testing"
)

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
			name: "test1",
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

func Test_toTable(t *testing.T) {
	type args struct {
		lines     []string
		pos       []int
		trimSpace bool
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toTable(tt.args.lines, tt.args.pos, tt.args.trimSpace); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
