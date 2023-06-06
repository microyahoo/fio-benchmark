package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of fio benchmark",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("fio-benchmark version is v0.1")
	},
	TraverseChildren: true,
}

func init() {}
