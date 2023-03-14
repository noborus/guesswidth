package cmd

import (
	"encoding/csv"
	"log"
	"os"

	"github.com/noborus/guesswidth"
	"github.com/spf13/cobra"
)

// csvCmd represents the csv command
var csvCmd = &cobra.Command{
	Use:   "csv",
	Short: "Output in csv format",
	Run: func(cmd *cobra.Command, args []string) {
		delimiter, err := cmd.PersistentFlags().GetString("delimiter")
		if err != nil {
			log.Fatal(err)
		}
		var d rune = ','
		if delimiter != "" {
			d = []rune(delimiter)[0]
		}
		toCSV(d)
	},
}

func toCSV(delimiter rune) {
	g := guesswidth.NewReader(os.Stdin)
	g.Header = header - 1
	g.LimitSplit = LimitSplit
	g.TrimSpace = true

	w := csv.NewWriter(os.Stdout)
	w.Comma = delimiter
	for {
		record, err := g.Read()
		if err != nil {
			break
		}
		if err := w.Write(record); err != nil {
			log.Fatal(err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(csvCmd)
	csvCmd.PersistentFlags().String("delimiter", ",", "delimiter")
}
