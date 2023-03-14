package cmd

import (
	"fmt"
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
	fence      string
	header     int
	LimitSplit int
)

func writeTable(args []string) {
	g := guesswidth.New(os.Stdin)
	g.Header = header - 1
	g.LimitSplit = LimitSplit
	g.TrimSpace = false
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&fence, "fence", "|", "fence")
	rootCmd.PersistentFlags().IntVar(&header, "header", 1, "header line number")
	rootCmd.PersistentFlags().IntVar(&LimitSplit, "split", -1, "maximum number of splits")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

}
