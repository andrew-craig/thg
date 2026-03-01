package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrewstuart/thg/config"
	"github.com/andrewstuart/thg/db"
	"github.com/andrewstuart/thg/things"
)

var (
	updateAuthToken   string
	updateTitle       string
	updateWhen        string
	updateDeadline    string
	updateTags        string
	updateNotes       string
	updateAppendNotes string
	updateList        string
	updateCanceled    bool
)

var updateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := config.LoadAuthToken(updateAuthToken)
		if err != nil {
			return err
		}

		database, err := db.Open()
		if err != nil {
			return err
		}
		defer database.Close()

		task, err := db.ResolveTask(database, args[0])
		if err != nil {
			return err
		}

		params := map[string]string{
			"auth-token": token,
			"id":         task.UUID,
		}

		if updateTitle != "" {
			params["title"] = updateTitle
		}
		if updateWhen != "" {
			params["when"] = updateWhen
		}
		if cmd.Flags().Changed("deadline") {
			params["deadline"] = updateDeadline
		}
		if updateTags != "" {
			params["add-tags"] = updateTags
		}
		if updateNotes != "" {
			params["notes"] = updateNotes
		}
		if updateAppendNotes != "" {
			params["append-notes"] = updateAppendNotes
		}
		if updateList != "" {
			params["list"] = updateList
		}
		if updateCanceled {
			params["canceled"] = "true"
		}

		if err := things.UpdateTask(params); err != nil {
			return err
		}

		fmt.Printf("Updated: %s\n", task.Title)
		return nil
	},
}

func init() {
	updateCmd.Flags().StringVar(&updateAuthToken, "auth-token", "", "Things URL scheme auth token")
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "new title")
	updateCmd.Flags().StringVar(&updateWhen, "when", "", "when to schedule")
	updateCmd.Flags().StringVar(&updateDeadline, "deadline", "", "deadline (YYYY-MM-DD, or empty to clear)")
	updateCmd.Flags().StringVar(&updateTags, "tags", "", "tags to add (comma-separated)")
	updateCmd.Flags().StringVar(&updateNotes, "notes", "", "replace notes")
	updateCmd.Flags().StringVar(&updateAppendNotes, "append-notes", "", "append to notes")
	updateCmd.Flags().StringVar(&updateList, "list", "", "move to project/area")
	updateCmd.Flags().BoolVar(&updateCanceled, "canceled", false, "cancel the task")
	rootCmd.AddCommand(updateCmd)
}
