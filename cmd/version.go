package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本号",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("wolf-k8s-cli version v0.0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
