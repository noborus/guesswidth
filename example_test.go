package guesswidth_test

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/noborus/guesswidth"
)

func ExampleReader() {
	in := `USER         PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root           1  0.0  0.0 168720 13788 ?        Ss   Mar11   0:50 /sbin/init splash
root           2  0.0  0.0      0     0 ?        S    Mar11   0:00 [kthreadd]
root           3  0.0  0.0      0     0 ?        I<   Mar11   0:00 [rcu_gp]
noborus   660384  0.0  0.3 3505316 115544 ?      Sl   07:43   0:09 gjs /usr/share/gnome-shell/extensions/ding@rastersoft.com/ding.js -E -P /usr/share/gnome-shell/extensions/ding@rastersoft.com -M 1 -D 1920:0:1920:1080:1:0:0:0:0:0 -D 0:0:1920:1200:2:27:0:0:0:1
noborus   735125  0.0  0.0  10268  2420 pts/3    S+   11:24   0:00 grep --color gnome-shell
`
	r := guesswidth.NewReader(strings.NewReader(in))

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(record)
	}
	// Output:
	// [USER PID %CPU %MEM VSZ RSS TTY STAT START TIME COMMAND]
	// [root 1 0.0 0.0 168720 13788 ? Ss Mar11 0:50 /sbin/init splash]
	// [root 2 0.0 0.0 0 0 ? S Mar11 0:00 [kthreadd]]
	// [root 3 0.0 0.0 0 0 ? I< Mar11 0:00 [rcu_gp]]
	// [noborus 660384 0.0 0.3 3505316 115544 ? Sl 07:43 0:09 gjs /usr/share/gnome-shell/extensions/ding@rastersoft.com/ding.js -E -P /usr/share/gnome-shell/extensions/ding@rastersoft.com -M 1 -D 1920:0:1920:1080:1:0:0:0:0:0 -D 0:0:1920:1200:2:27:0:0:0:1]
	// [noborus 735125 0.0 0.0 10268 2420 pts/3 S+ 11:24 0:00 grep --color gnome-shell]
}

func ExampleToTable() {
	lines := []string{
		"    PID TTY          TIME CMD",
		"1595989 pts/6    00:00:01 zsh",
		"1690373 pts/6    00:00:00 ps",
	}
	table := guesswidth.ToTable(lines, 1, true)
	fmt.Println(table)
	// Output:
	// [[PID TTY TIME CMD] [1595989 pts/6 00:00:01 zsh] [1690373 pts/6 00:00:00 ps]]
}

func ExampleToTableN() {
	lines := []string{
		"2022-12-21T09:50:16+0000 WARN A warning that should be ignored is usually at this level and should be actionable.",
		"2022-12-21T09:50:17+0000 INFO This is less important than debug log and is often used to provide context in the current task.",
		"2022-12-10T05:33:53+0000 DEBUG This is a debug log that shows a log that can be ignored.",
		"2022-12-10T05:33:53+0000 INFO This is less important than debug log and is often used to provide context in the current task.",
	}
	table := guesswidth.ToTableN(lines, 1, 2, true)
	for _, columns := range table {
		fmt.Println(strings.Join(columns, ","))
	}
	// Output:
	// 2022-12-21T09:50:16+0000,WARN,A warning that should be ignored is usually at this level and should be actionable.
	// 2022-12-21T09:50:17+0000,INFO,This is less important than debug log and is often used to provide context in the current task.
	// 2022-12-10T05:33:53+0000,DEBUG,This is a debug log that shows a log that can be ignored.
	// 2022-12-10T05:33:53+0000,INFO,This is less important than debug log and is often used to provide context in the current task.
}
