package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/andrewstuart/thg/things"
)

var (
	addWhen      string
	addDeadline  string
	addList      string
	addTags      string
	addNotes     string
	addChecklist string
	addHeading   string
	addReveal    bool
)

var addCmd = &cobra.Command{
	Use:   "add <title>",
	Short: "Add a new task",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := strings.Join(args, " ")

		params := map[string]string{
			"title": title,
		}

		if addWhen != "" {
			params["when"] = addWhen
		}
		if addDeadline != "" {
			params["deadline"] = addDeadline
		}
		if addList != "" {
			params["list"] = addList
		}
		if addTags != "" {
			params["tags"] = addTags
		}
		if addNotes != "" {
			params["notes"] = addNotes
		}
		if addChecklist != "" {
			// Convert comma-separated to newline-separated
			items := strings.ReplaceAll(addChecklist, ",", "\n")
			params["checklist-items"] = items
		}
		if addHeading != "" {
			params["heading"] = addHeading
		}
		if addReveal {
			params["reveal"] = "true"
		}

		if err := things.AddTodo(params); err != nil {
			return err
		}

		fmt.Printf("Added: %s\n", title)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVar(&addWhen, "when", "", "when to schedule (today/tomorrow/evening/anytime/someday/YYYY-MM-DD)")
	addCmd.Flags().StringVar(&addDeadline, "deadline", "", "deadline (YYYY-MM-DD)")
	addCmd.Flags().StringVar(&addList, "list", "", "project or area name")
	addCmd.Flags().StringVar(&addTags, "tags", "", "comma-separated tag names")
	addCmd.Flags().StringVar(&addNotes, "notes", "", "task notes")
	addCmd.Flags().StringVar(&addChecklist, "checklist", "", "comma-separated checklist items")
	addCmd.Flags().StringVar(&addHeading, "heading", "", "heading within a project")
	addCmd.Flags().BoolVar(&addReveal, "reveal", false, "show task in Things after adding")
	rootCmd.AddCommand(addCmd)
}
