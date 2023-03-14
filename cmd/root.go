package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/noborus/guesswidth"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "guesswidth",
	Short: "Guess the width of the column and split it",
	Long: `Guess the width of the columns from the header and body,
split them, insert fences and output.`,
	Run: func(cmd *cobra.Command, args []string) {
		writeTable(args)
	},
}

var (
	fence    string
	header   int
	numSplit int
)

func writeTable(args []string) {
	lines := readAll(os.Stdin)
	table := toTable(lines, false)
	write(table)
}

func readAll(r io.Reader) []string {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func toTable(lines []string, trimSpace bool) [][]string {
	var table [][]string
	if numSplit > 0 {
		table = guesswidth.ToTableN(lines, header-1, numSplit, trimSpace)
	} else {
		table = guesswidth.ToTable(lines, header-1, trimSpace)
	}
	return table
}

func write(table [][]string) {
	for _, row := range table {
		for n, col := range row {
			if n > 0 {
				fmt.Print(fence)
			}
			fmt.Printf("%s", col)
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
	rootCmd.PersistentFlags().IntVar(&numSplit, "split", -1, "number to split")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}
