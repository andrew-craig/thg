package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/andrew-craig/thg/db"
)

var showCmd = &cobra.Command{
	Use:   "show <id-or-title>",
	Short: "Show task details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return err
		}
		defer database.Close()

		task, err := db.ResolveTask(database, args[0])
		if err != nil {
			return err
		}

		// Load tags and checklist
		tags, _ := db.TagsForTask(database, task.UUID)
		checklist, _ := db.ChecklistForTask(database, task.UUID)

		typeName := "To-Do"
		if task.Type == 1 {
			typeName = "Project"
		}

		pairs := [][]string{
			{"ID", task.UUID},
			{"Title", task.Title},
			{"Type", typeName},
			{"Status", task.StatusText()},
			{"When", task.WhenText()},
		}

		if task.StartDate > 0 {
			pairs = append(pairs, []string{"Start Date", db.FormatDate(task.StartDate)})
		}
		if task.Deadline > 0 {
			pairs = append(pairs, []string{"Deadline", db.FormatDate(task.Deadline)})
		}
		if task.ProjectTitle != "" {
			pairs = append(pairs, []string{"Project", task.ProjectTitle})
		}
		if task.HeadingTitle != "" {
			pairs = append(pairs, []string{"Heading", task.HeadingTitle})
		}
		if task.AreaTitle != "" {
			pairs = append(pairs, []string{"Area", task.AreaTitle})
		}
		if len(tags) > 0 {
			pairs = append(pairs, []string{"Tags", strings.Join(tags, ", ")})
		}
		if task.Type == 1 {
			pairs = append(pairs, []string{"Tasks", fmt.Sprintf("%d open / %d total", task.OpenTasks, task.TotalTasks)})
		}
		if task.CreationDate > 0 {
			pairs = append(pairs, []string{"Created", db.FormatTimestamp(task.CreationDate)})
		}
		if task.ModDate > 0 {
			pairs = append(pairs, []string{"Modified", db.FormatTimestamp(task.ModDate)})
		}

		printer.Detail(pairs)

		if !jsonOutput {
			if task.Notes != "" {
				printer.Block("Notes", task.Notes, os.Stdout)
			}

			if len(checklist) > 0 {
				fmt.Println("\nChecklist:")
				for _, ci := range checklist {
					check := "[ ]"
					if ci.Status == 3 {
						check = "[x]"
					}
					fmt.Printf("  %s %s\n", check, ci.Title)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(showCmd)
}
