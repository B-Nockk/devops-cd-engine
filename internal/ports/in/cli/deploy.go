package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"cd-engine/internal/adapters/healthcheck"
	"cd-engine/internal/adapters/store/sqlite"
	"cd-engine/internal/app"
	"cd-engine/internal/domain"
	"cd-engine/internal/domain/strategy"
)

// Dummy Executor mock matching our out.Executor interface exactly
type MockExecutor struct{}

func (m *MockExecutor) Execute(ctx context.Context, target string, command string) (string, error) {
	fmt.Printf("[MockExecutor] Target: %s | Executing: %s\n", target, command)
	return "mock output", nil
}

var deployCmd = &cobra.Command{
	Use:   "deploy <release_id>",
	Short: "Trigger a deployment for a given release ID",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Assuming domain.ID is just a string type wrapper, cast it.
		// If you built an IDFromString parser, use that.
		releaseID := domain.ID(args[0])

		ctx := context.Background()

		// 1. Initialize adapters
		store, err := sqlite.NewStore("data.db")
		if err != nil {
			return fmt.Errorf("failed to init store: %w", err)
		}

		healthChecker := healthcheck.NewChecker()
		executor := &MockExecutor{}

		// 2. Initialize the Planner (Defaulting to Blue/Green for this test)
		planner := &strategy.BlueGreenPlanner{}

		// 3. Instantiate orchestrator
		// Pass 'store' multiple times because it satisfies Tenant, Env, and Release store interfaces
		orchestrator := app.NewDeploymentOrchestrator(
			store,
			store,
			store,
			planner,
			executor,
			healthChecker,
		)

		// 4. Run deployment
		fmt.Printf("Starting deployment for Release ID: %s...\n", releaseID)
		if err := orchestrator.RunDeploy(ctx, releaseID); err != nil {
			fmt.Printf("[FAIL]: Deployment failed: %v\n", err)
			return err
		}

		fmt.Println("[SUCCESS] Deployment succeeded!")
		return nil
	},
}

func init() {
	// Assuming rootCmd is defined in root.go
	rootCmd.AddCommand(deployCmd)
}
