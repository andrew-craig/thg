package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appVersion = "dev"

func SetVersion(v string) {
	appVersion = v
	rootCmd.Version = v
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("thg", appVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
