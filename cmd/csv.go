package cmd

import (
	"encoding/csv"
	"log"
	"os"

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
	lines := readAll(os.Stdin)
	table := toTable(lines, true)
	w := csv.NewWriter(os.Stdout)
	w.Comma = delimiter
	for _, record := range table {
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
