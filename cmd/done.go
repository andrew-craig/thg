package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrew-craig/thg/config"
	"github.com/andrew-craig/thg/db"
	"github.com/andrew-craig/thg/things"
)

var doneAuthToken string

var doneCmd = &cobra.Command{
	Use:   "done <id>",
	Short: "Complete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token, err := config.LoadAuthToken(doneAuthToken)
		if err != nil {
			return err
		}

		// Resolve task to get full UUID
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
			"completed":  "true",
		}

		if err := things.UpdateTask(params); err != nil {
			return err
		}

		fmt.Printf("Completed: %s\n", task.Title)
		return nil
	},
}

func init() {
	doneCmd.Flags().StringVar(&doneAuthToken, "auth-token", "", "Things URL scheme auth token")
	rootCmd.AddCommand(doneCmd)
}
