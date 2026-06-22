package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrew-craig/thg/db"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "List tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return err
		}
		defer database.Close()

		tags, err := db.ListTags(database)
		if err != nil {
			return err
		}

		if len(tags) == 0 {
			fmt.Println("No tags found.")
			return nil
		}

		headers := []string{"Title", "Shortcut"}
		var rows [][]string
		for _, t := range tags {
			rows = append(rows, []string{t.Title, t.Shortcut})
		}
		printer.Table(headers, rows)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tagsCmd)
}
