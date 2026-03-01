package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrewstuart/thg/db"
)

var (
	listArea string
	listTag  string
	listAll  bool
)

var listCmd = &cobra.Command{
	Use:   "list [filter]",
	Short: "List open tasks",
	Long: `List open tasks. Filters: today, inbox, someday, upcoming, or a project name.
Without a filter, shows today's tasks.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filter := "today"
		if len(args) > 0 {
			filter = args[0]
		}
		return listRun(filter)
	},
}

func listRun(filter string) error {
	database, err := db.Open()
	if err != nil {
		return err
	}
	defer database.Close()

	var tasks []db.Task

	// Check for area/tag filters first
	if listArea != "" {
		area, err := db.ResolveArea(database, listArea)
		if err != nil {
			return err
		}
		tasks, err = db.ListByArea(database, area.UUID)
		if err != nil {
			return err
		}
	} else if listTag != "" {
		tag, err := db.ResolveTag(database, listTag)
		if err != nil {
			return fmt.Errorf("tag %q not found", listTag)
		}
		tasks, err = db.ListByTag(database, tag.UUID)
		if err != nil {
			return err
		}
	} else if listAll {
		tasks, err = db.ListAll(database)
		if err != nil {
			return err
		}
	} else {
		switch filter {
		case "today":
			tasks, err = db.ListToday(database)
		case "inbox":
			tasks, err = db.ListInbox(database)
		case "someday":
			tasks, err = db.ListSomeday(database)
		case "upcoming":
			tasks, err = db.ListUpcoming(database)
		default:
			// Treat as project name
			project, err := db.ResolveProject(database, filter)
			if err != nil {
				return err
			}
			tasks, err = db.ListByProject(database, project.UUID)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}

	if len(tasks) == 0 {
		if !jsonOutput {
			fmt.Println("No tasks found.")
		} else {
			printer.Table([]string{"ID", "Title", "Project", "When", "Deadline"}, nil)
		}
		return nil
	}

	headers := []string{"ID", "Title", "Project", "When", "Deadline"}
	var rows [][]string
	for _, t := range tasks {
		deadline := db.FormatDate(t.Deadline)
		rows = append(rows, []string{
			t.ShortID(),
			t.Title,
			t.ProjectTitle,
			t.WhenText(),
			deadline,
		})
	}
	printer.Table(headers, rows)
	return nil
}

func init() {
	listCmd.Flags().StringVar(&listArea, "area", "", "filter by area name")
	listCmd.Flags().StringVar(&listTag, "tag", "", "filter by tag name")
	listCmd.Flags().BoolVar(&listAll, "all", false, "show all open tasks")
	rootCmd.AddCommand(listCmd)
}
