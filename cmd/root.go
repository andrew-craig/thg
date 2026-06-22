package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/andrew-craig/thg/format"
)

var (
	jsonOutput bool
	printer    *format.Printer
)

var rootCmd = &cobra.Command{
	Use:   "thg",
	Short: "Things 3 CLI",
	Long:  "Command-line interface for reading and writing tasks in Things 3.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		printer = &format.Printer{JSON: jsonOutput}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Default: show today list
		return listRun("today")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&jsonOutput, "json", false, "output as JSON")
}
