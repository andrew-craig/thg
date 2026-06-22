package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/andrew-craig/thg/db"
)

var (
	projectsArea string
	projectsAll  bool
)

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List open projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		database, err := db.Open()
		if err != nil {
			return err
		}
		defer database.Close()

		projects, err := db.ListProjects(database, projectsArea, projectsAll)
		if err != nil {
			return err
		}

		if len(projects) == 0 {
			fmt.Println("No projects found.")
			return nil
		}

		headers := []string{"ID", "Title", "Area", "Open", "Total", "When"}
		var rows [][]string
		for _, p := range projects {
			rows = append(rows, []string{
				p.ShortID(),
				p.Title,
				p.AreaTitle,
				fmt.Sprintf("%d", p.OpenTasks),
				fmt.Sprintf("%d", p.TotalTasks),
				p.WhenText(),
			})
		}
		printer.Table(headers, rows)
		return nil
	},
}

func init() {
	projectsCmd.Flags().StringVar(&projectsArea, "area", "", "filter by area")
	projectsCmd.Flags().BoolVar(&projectsAll, "all", false, "include completed/cancelled projects")
	rootCmd.AddCommand(projectsCmd)
}
