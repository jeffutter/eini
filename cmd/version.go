package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "master"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the eini version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func Version() string {
	return version
}
