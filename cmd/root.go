package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/noborus/guesswidth"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "guesswidth",
	Short: "Guess the width of the column and split it",
	Long: `Guess the width of the columns from the header and body,
split them, insert fences and output.`,
	Version: guesswidth.Version(),
	Run: func(cmd *cobra.Command, args []string) {
		writeTable()
	},
}

var (
	fence      string
	header     int
	limitSplit int
	scanNum    int
	align      bool
)

func writeTable() {
	g := guesswidth.NewReader(os.Stdin)
	g.Header = header - 1
	g.LimitSplit = limitSplit
	g.TrimSpace = false
	if scanNum > 0 {
		g.ScanNum = scanNum
	}
	if align {
		writeAlign(g)
		return
	}
	write(g)
}

func write(g *guesswidth.GuessWidth) {
	for {
		row, err := g.Read()
		if err != nil {
			break
		}
		for n, col := range row {
			if n > 0 {
				fmt.Print(fence)
			}
			fmt.Printf("%s", col)
		}
		fmt.Println()
	}
}

func writeAlign(g *guesswidth.GuessWidth) {
	for _, row := range g.ReadAll() {
		for n, col := range row {
			if n > 0 {
				fmt.Print(fence)
			}
			col = strings.TrimSpace(col)
			if g.Widths[n].Justified == guesswidth.Right {
				fmt.Printf("%*s", g.Widths[n].Width, col)
			} else {
				if len(g.Widths)-1 == n {
					fmt.Printf("%s", col)
				} else {
					fmt.Printf("%-*s", g.Widths[n].Width, col)
				}
			}
		}
		fmt.Println()
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&fence, "fence", "|", "fence")
	rootCmd.PersistentFlags().IntVar(&header, "header", 1, "header line number")
	rootCmd.PersistentFlags().IntVar(&limitSplit, "split", -1, "maximum number of splits")
	rootCmd.PersistentFlags().IntVar(&scanNum, "scannum", 100, "number of line to scan")
	rootCmd.PersistentFlags().BoolVarP(&align, "align", "a", false, "align the output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
