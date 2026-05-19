package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"cd-engine/internal/adapters/store/sqlite"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Initialize the database and run migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running database migrations...")
		_, err := sqlite.NewStore("data.db")
		if err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		fmt.Println("[SUCCESS] Database initialized successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
