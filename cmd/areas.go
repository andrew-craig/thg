package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrew-craig/thg/db"
)

var areasCmd = &cobra.Command{
	Use:   "areas",
	Short: "List areas",
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return err
		}
		defer database.Close()

		areas, err := db.ListAreas(database)
		if err != nil {
			return err
		}

		if len(areas) == 0 {
			fmt.Println("No areas found.")
			return nil
		}

		headers := []string{"ID", "Title"}
		var rows [][]string
		for _, a := range areas {
			rows = append(rows, []string{a.UUID[:6], a.Title})
		}
		printer.Table(headers, rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(areasCmd)
}
