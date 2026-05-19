package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cd-engine",
	Short: "cd-engine is a Continuous Deployment engine",
	Long:  "cd-engine provides deployment orchestration via CLI commands.",
}

func Execute() error {
	return rootCmd.Execute()
}
